package ai

import (
	"context"
	"log"

	"courses-service/src/model"

	"google.golang.org/genai"
)

type AiClient struct {
	GeminiApiKey string
	Client       *genai.Client
}

// const aiModel = "gemini-2.0-flash"
func NewAiClient(geminiApiKey string) *AiClient {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: geminiApiKey,
	})
	if err != nil {
		log.Fatal("Failed to create Gemini client", err)
		return nil
	}
	log.Println("Gemini client created")
	return &AiClient{
		GeminiApiKey: geminiApiKey,
		Client:       client,
	}
}

func (c *AiClient) SummarizeCourseFeedbacks(feedbacks []model.CourseFeedback) string {
	return ""
}
