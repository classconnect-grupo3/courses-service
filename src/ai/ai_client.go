package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"courses-service/src/config"
	"courses-service/src/model"

	"google.golang.org/genai"
)

type AiClient struct {
	context      context.Context
	GeminiApiKey string
	Client       *genai.Client
}

const aiModel = "gemini-2.0-flash"

func generateCourseFeedbacksPrompt(feedbacks []*model.CourseFeedback) string {
	prompt := SummarizeCourseFeedbacksPrompt
	for _, feedback := range feedbacks {
		prompt += fmt.Sprintf("Puntuacion: %d\n", feedback.Score)
		prompt += fmt.Sprintf("Tipo: %s\n", feedback.FeedbackType)
		prompt += fmt.Sprintf("Feedback: %s\n", feedback.Feedback)
	}
	return prompt
}

func NewAiClient(config *config.Config) *AiClient {

	if config.Environment == "test" {
		return &AiClient{
			context:      context.Background(),
			GeminiApiKey: "",
			Client:       nil,
		}
	}
	geminiApiKey := config.GeminiApiKey
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: geminiApiKey,
	})
	if err != nil {
		log.Fatal("Failed to create Gemini client", err)
		return nil
	}
	log.Println("Gemini client created")
	return &AiClient{
		context:      ctx,
		GeminiApiKey: geminiApiKey,
		Client:       client,
	}
}

func (c *AiClient) SummarizeCourseFeedbacks(feedbacks []*model.CourseFeedback) string {
	prompt := generateCourseFeedbacksPrompt(feedbacks)
	response, err := c.Client.Models.GenerateContent(c.context, aiModel, genai.Text(prompt), nil)
	if err != nil {
		log.Fatal("Failed to generate content", err)
	}
	return debugString(response)
}

func debugString[T any](r *T) string {

	response, err := json.MarshalIndent(*r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(response)
}
