package service

import (
	"context"
	"fmt"
	"minos/internal/dto"
	"minos/internal/llm"
	"minos/internal/llm/gemini"
	"minos/internal/model"
	"minos/internal/repository"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
)

type ChatService interface {
	SendMessage(interviewID uuid.UUID, req *dto.SendMessageRequest) (*dto.SendMessageResponse, error)
	GetHistory(interviewID uuid.UUID) ([]model.Message, error)
}

type chatService struct {
	msgRepo       repository.MessageRepository
	interviewRepo repository.InterviewRepository
	geminiClient  *gemini.Client
}

func NewChatService(msgRepo repository.MessageRepository, interviewRepo repository.InterviewRepository, geminiClient *gemini.Client) ChatService {
	return &chatService{
		msgRepo:       msgRepo,
		interviewRepo: interviewRepo,
		geminiClient:  geminiClient,
	}
}

func (s *chatService) SendMessage(interviewID uuid.UUID, req *dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	// 1. Validate Interview
	interview, err := s.interviewRepo.FindInterviewByID(interviewID)
	if err != nil {
		return nil, err
	}
	if interview.Status != model.InterviewStatusActive {
		return nil, fmt.Errorf("interview is not active")
	}

	// 2. Prepare User Content
	userContent := req.Content
	if req.Code != "" {
		lang := req.Language
		if lang == "" {
			lang = "unknown"
		}
		userContent += fmt.Sprintf("\n\n[USER ATTACHED CODE (%s)]:\n```%s\n%s\n```\n(Please review this code as part of the interview context)", lang, lang, req.Code)
	}

	// 3. Save User Message
	userMsg := &model.Message{
		InterviewID: interviewID,
		Role:        model.MessageRoleUser,
		Content:     userContent, // Save full content including attached code
	}
	if err := s.msgRepo.CreateMessage(userMsg); err != nil {
		return nil, err
	}

	// 4. Build History for Gemini
	history, err := s.msgRepo.FindMessagesByInterviewID(interviewID)
	if err != nil {
		return nil, err
	}

	var geminiHistory []*genai.Content

	// Prepend System Prompt logic
	systemInstruction := fmt.Sprintf(llm.SystemPromptInterviewer, string(interview.ProblemSnapshot))

	// We construct the chat history for context
	for _, msg := range history {
		// Skip the current user message we just saved, as we will send it in SendMessage
		if msg.ID == userMsg.ID {
			continue
		}
		role := "user"
		if msg.Role == model.MessageRoleAssistant {
			role = "model"
		}
		geminiHistory = append(geminiHistory, &genai.Content{
			Role:  role,
			Parts: []genai.Part{genai.Text(msg.Content)},
		})
	}

	// Prepend system prompt as fake turn
	fullHistory := []*genai.Content{
		{
			Role:  "user",
			Parts: []genai.Part{genai.Text(systemInstruction)},
		},
		{
			Role:  "model",
			Parts: []genai.Part{genai.Text("Understood. I am ready to conduct the interview.")},
		},
	}
	fullHistory = append(fullHistory, geminiHistory...)

	cs := s.geminiClient.StartChat(fullHistory)
	resp, err := cs.SendMessage(context.Background(), genai.Text(userContent))
	if err != nil {
		return nil, err
	}

	// 5. Extract Response
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from AI")
	}
	aiText := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			aiText += string(txt)
		}
	}

	// 6. Save AI Message
	aiMsg := &model.Message{
		InterviewID: interviewID,
		Role:        model.MessageRoleAssistant,
		Content:     aiText,
	}
	if err := s.msgRepo.CreateMessage(aiMsg); err != nil {
		return nil, err
	}

	return &dto.SendMessageResponse{
		MessageID:  aiMsg.ID,
		AIResponse: aiText,
	}, nil
}

func (s *chatService) GetHistory(interviewID uuid.UUID) ([]model.Message, error) {
	return s.msgRepo.FindMessagesByInterviewID(interviewID)
}
