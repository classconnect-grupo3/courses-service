package ai

import (
	"context"
	"errors"
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

func generateStudentFeedbacksPrompt(feedbacks []*model.StudentFeedback) string {
	prompt := SummarizeStudentFeedbacksPrompt
	for _, feedback := range feedbacks {
		prompt += fmt.Sprintf("Puntuacion: %d\n", feedback.Score)
		prompt += fmt.Sprintf("Tipo: %s\n", feedback.FeedbackType)
		prompt += fmt.Sprintf("Feedback: %s\n", feedback.Feedback)
	}
	return prompt
}

func generateCourseFeedbacksPrompt(feedbacks []*model.CourseFeedback) string {
	prompt := SummarizeCourseFeedbacksPrompt
	for _, feedback := range feedbacks {
		prompt += fmt.Sprintf("Puntuacion: %d\n", feedback.Score)
		prompt += fmt.Sprintf("Tipo: %s\n", feedback.FeedbackType)
		prompt += fmt.Sprintf("Feedback: %s\n", feedback.Feedback)
	}
	return prompt
}

func generateSubmissionFeedbackPrompt(score *float64, feedback string) string {
	prompt := SummarizeSubmissionFeedbackPrompt
	if score != nil {
		prompt += fmt.Sprintf("Puntuacion: %.2f\n", *score)
	} else {
		prompt += "Puntuacion: No asignada\n"
	}
	prompt += fmt.Sprintf("Feedback: %s\n", feedback)
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

func (c *AiClient) SummarizeCourseFeedbacks(feedbacks []*model.CourseFeedback) (string, error) {
	prompt := generateCourseFeedbacksPrompt(feedbacks)
	response, err := c.Client.Models.GenerateContent(c.context, aiModel, genai.Text(prompt), nil)
	if err != nil {
		log.Fatal("Failed to generate content", err)
		return "", err
	}
	return obtainAnswerFromModel(response)
}

func (c *AiClient) SummarizeStudentFeedbacks(feedbacks []*model.StudentFeedback) (string, error) {
	prompt := generateStudentFeedbacksPrompt(feedbacks)
	response, err := c.Client.Models.GenerateContent(c.context, aiModel, genai.Text(prompt), nil)
	if err != nil {
		log.Fatal("Failed to generate content", err)
		return "", err
	}
	return obtainAnswerFromModel(response)
}

func (c *AiClient) SummarizeSubmissionFeedback(score *float64, feedback string) (string, error) {
	if c.Client == nil {
		// Return a mock response for test environment
		return "Resumen de retroalimentación generado por IA (entorno de test)", nil
	}

	prompt := generateSubmissionFeedbackPrompt(score, feedback)
	response, err := c.Client.Models.GenerateContent(c.context, aiModel, genai.Text(prompt), nil)
	if err != nil {
		log.Fatal("Failed to generate content", err)
		return "", err
	}
	return obtainAnswerFromModel(response)
}

func obtainAnswerFromModel(result *genai.GenerateContentResponse) (string, error) {
	if len(result.Candidates) == 0 {
		return "", errors.New("no answer found")
	}
	answer := result.Candidates[0]

	for _, part := range answer.Content.Parts {
		if len(part.Text) > 0 {
			return part.Text, nil
		}
	}
	return "", errors.New("no answer found")
}
