package service

import (
	"context"
	"encoding/json"
	"fmt"
	"minos/internal/dto"
	"minos/internal/llm"
	"minos/internal/llm/gemini"
	"minos/internal/model"
	"minos/internal/repository"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type InterviewService interface {
	StartInterview(req *dto.StartInterviewRequest) (*dto.StartInterviewResponse, error)
	GetInterview(id uuid.UUID) (*model.Interview, error)
	EndInterview(id uuid.UUID) (*dto.EndInterviewResponse, error)
}

type interviewService struct {
	repo           repository.InterviewRepository
	evalRepo       repository.EvaluationRepository
	msgRepo        repository.MessageRepository
	submissionRepo repository.SubmissionRepository
	geminiClient   *gemini.Client
}

func NewInterviewService(
	repo repository.InterviewRepository,
	evalRepo repository.EvaluationRepository,
	msgRepo repository.MessageRepository,
	submissionRepo repository.SubmissionRepository,
	geminiClient *gemini.Client,
) InterviewService {
	return &interviewService{
		repo:           repo,
		evalRepo:       evalRepo,
		msgRepo:        msgRepo,
		submissionRepo: submissionRepo,
		geminiClient:   geminiClient,
	}
}

func (s *interviewService) StartInterview(req *dto.StartInterviewRequest) (*dto.StartInterviewResponse, error) {
	// 1. Create Interview Record
	interview := &model.Interview{
		UserID:          req.UserID,
		ProblemID:       req.ProblemID,
		ProblemSnapshot: req.ProblemSnapshot,
		Status:          model.InterviewStatusActive,
	}
	if err := s.repo.CreateInterview(interview); err != nil {
		return nil, err
	}

	// 2. Generate Greeting using Gemini
	prompt := fmt.Sprintf(llm.SystemPromptInterviewer, string(req.ProblemSnapshot)) + "\n\nPlease start the interview by greeting the candidate and asking them to explain their initial thought process."
	resp, err := s.geminiClient.GenerateContent(context.Background(), prompt)
	greeting := "Hello! I'm ready to help you with this problem. How would you like to start?" // Default fallback
	if err == nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		// Extract text
		// Note: Simplification here
		greeting = fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	}

	// 3. Save Greeting as first message
	msg := &model.Message{
		InterviewID: interview.ID,
		Role:        model.MessageRoleAssistant,
		Content:     greeting,
	}
	s.msgRepo.CreateMessage(msg)

	return &dto.StartInterviewResponse{
		InterviewID: interview.ID,
		Greeting:    greeting,
	}, nil
}

func (s *interviewService) GetInterview(id uuid.UUID) (*model.Interview, error) {
	return s.repo.FindInterviewByID(id)
}

func (s *interviewService) EndInterview(id uuid.UUID) (*dto.EndInterviewResponse, error) {
	interview, err := s.repo.FindInterviewByID(id)
	if err != nil {
		return nil, err
	}

	if interview.Status == model.InterviewStatusCompleted {
		// Already completed, return existing evaluation
		eval, err := s.evalRepo.FindEvaluationByInterviewID(id)
		if err != nil {
			return nil, err
		}
		return &dto.EndInterviewResponse{
			EvaluationID: eval.ID,
			OverallScore: eval.OverallScore,
			Feedback:     eval.DetailedFeedback,
		}, nil
	}

	// 1. Update Status
	now := time.Now()
	interview.Status = model.InterviewStatusCompleted
	interview.EndedAt = &now
	s.repo.UpdateInterview(interview)

	// 2. Gather Context
	msgs, _ := s.msgRepo.FindMessagesByInterviewID(id)
	submissions, _ := s.submissionRepo.FindSubmissionsByInterviewID(id)

	transcript := ""
	for _, m := range msgs {
		transcript += fmt.Sprintf("[%s]: %s\n", m.Role, m.Content)
	}

	subsText := ""
	for _, sub := range submissions {
		subsText += fmt.Sprintf("Code (%s): %s\nResult: %s\n\n", sub.Language, sub.Code, sub.AIFeedback)
	}

	// 3. Call Gemini for Evaluation
	prompt := fmt.Sprintf(llm.SystemPromptEvaluator, string(interview.ProblemSnapshot), transcript, subsText)
	// Force JSON structure?
	prompt += "\nPlease output the result as a valid JSON object with keys: problem_solving_score, code_quality_score, communication_score, technical_score, overall_score, strengths (array), improvements (array), detailed_feedback."

	resp, err := s.geminiClient.GenerateContent(context.Background(), prompt)
	if err != nil {
		return nil, err
	}

	// 4. Parse JSON Response
	// This is tricky without strict mode. We assume Gemini follows instructions.
	// For MVP, we might need to clean the response string (remove markdown code blocks).
	responseText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])
	responseText = strings.TrimPrefix(responseText, "```json")
	responseText = strings.TrimPrefix(responseText, "```")
	responseText = strings.TrimSuffix(responseText, "```")

	type EvalResult struct {
		ProblemSolvingScore int      `json:"problem_solving_score"`
		CodeQualityScore    int      `json:"code_quality_score"`
		CommunicationScore  int      `json:"communication_score"`
		TechnicalScore      int      `json:"technical_score"`
		OverallScore        int      `json:"overall_score"`
		Strengths           []string `json:"strengths"`
		Improvements        []string `json:"improvements"`
		DetailedFeedback    string   `json:"detailed_feedback"`
	}

	var res EvalResult
	json.Unmarshal([]byte(responseText), &res)

	// 5. Save Evaluation
	strengthsJSON, _ := json.Marshal(res.Strengths)
	improvementsJSON, _ := json.Marshal(res.Improvements)

	evaluation := &model.Evaluation{
		InterviewID:         id,
		ProblemSolvingScore: res.ProblemSolvingScore,
		CodeQualityScore:    res.CodeQualityScore,
		CommunicationScore:  res.CommunicationScore,
		TechnicalScore:      res.TechnicalScore,
		OverallScore:        res.OverallScore,
		Strengths:           datatypes.JSON(strengthsJSON),
		Improvements:        datatypes.JSON(improvementsJSON),
		DetailedFeedback:    res.DetailedFeedback,
	}

	s.evalRepo.CreateEvaluation(evaluation)

	return &dto.EndInterviewResponse{
		EvaluationID: evaluation.ID,
		OverallScore: evaluation.OverallScore,
		Feedback:     evaluation.DetailedFeedback,
	}, nil
}

