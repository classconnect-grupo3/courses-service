package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"courses-service/src/config"
	"courses-service/src/model"
	"courses-service/src/schemas"

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

func generateSubmissionCorrectionPrompt(assignment *model.Assignment, submission *model.Submission) string {
	prompt := CorrectSubmissionPrompt

	// Add assignment info
	prompt += "ASSIGNMENT INFO:\n"
	prompt += fmt.Sprintf("Título: %s\n", assignment.Title)
	prompt += fmt.Sprintf("Puntaje Máximo: %.2f\n", assignment.TotalPoints)
	prompt += fmt.Sprintf("Tipo: %s\n\n", assignment.Type)

	// Create a map of assignment questions for easy lookup
	questionMap := make(map[string]model.Question)
	for _, question := range assignment.Questions {
		questionMap[question.ID] = question
	}

	// Add questions and student answers
	for _, answer := range submission.Answers {
		if question, exists := questionMap[answer.QuestionID]; exists {
			prompt += fmt.Sprintf("ID: %s\n", question.ID)
			prompt += fmt.Sprintf("Pregunta: %s\n", question.Text)
			prompt += fmt.Sprintf("Tipo: %s\n", question.Type)
			prompt += fmt.Sprintf("Puntaje: %.2f\n", question.Points)

			// Add correct answers
			if len(question.CorrectAnswers) > 0 {
				prompt += fmt.Sprintf("Respuestas Correctas: %v\n", question.CorrectAnswers)
			}

			// Add student answer
			prompt += fmt.Sprintf("Respuesta del Estudiante: %v\n", answer.Content)
			prompt += "\n---\n\n"
		}
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

func (c *AiClient) CorrectSubmission(assignment *model.Assignment, submission *model.Submission) (*schemas.AiCorrectionResponse, error) {
	if c.Client == nil {
		// Return a mock response for test environment
		return &schemas.AiCorrectionResponse{
			Score:             assignment.TotalPoints * 0.8, // Mock score
			Feedback:          "Corrección automática realizada en entorno de test",
			NeedsManualReview: false,
		}, nil
	}

	prompt := generateSubmissionCorrectionPrompt(assignment, submission)
	response, err := c.Client.Models.GenerateContent(c.context, aiModel, genai.Text(prompt), nil)
	if err != nil {
		log.Printf("Failed to generate correction content: %v", err)
		return nil, err
	}

	rawResponse, err := obtainAnswerFromModel(response)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response from AI
	var correctionResponse schemas.AiCorrectionResponse
	// Clean the response to extract only the JSON part
	cleanedResponse := strings.TrimSpace(rawResponse)
	if strings.Contains(cleanedResponse, "{") {
		// Find the first { and last }
		start := strings.Index(cleanedResponse, "{")
		end := strings.LastIndex(cleanedResponse, "}")
		if start != -1 && end != -1 && end > start {
			jsonStr := cleanedResponse[start : end+1]
			err = json.Unmarshal([]byte(jsonStr), &correctionResponse)
			if err != nil {
				log.Printf("Failed to parse AI correction response: %v", err)
				// Return a fallback response
				return &schemas.AiCorrectionResponse{
					Score:             0,
					Feedback:          "Error al procesar la corrección automática. Requiere revisión manual.",
					NeedsManualReview: true,
				}, nil
			}
		}
	}

	return &correctionResponse, nil
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
