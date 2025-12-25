package gemini

import (
	"context"
	"fmt"
	"minos/config"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

type Client struct {
	client *genai.Client
	model  *genai.GenerativeModel
	conf   *config.Config
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.Gemini.ApiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	modelName := cfg.Gemini.Model
	if modelName == "" {
		modelName = "gemini-1.5-pro"
	}
	model := client.GenerativeModel(modelName)

	log.Info().Str("model", modelName).Msg("Gemini client initialized")

	return &Client{
		client: client,
		model:  model,
		conf:   cfg,
	}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) StartChat(history []*genai.Content) *genai.ChatSession {
	cs := c.model.StartChat()
	if len(history) > 0 {
		cs.History = history
	}
	return cs
}

func (c *Client) GenerateContent(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	return c.model.GenerateContent(ctx, genai.Text(prompt))
}

