package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mh-go/internal/models"
	"net/http"
	"strings"
	"time"
)

type OpenAIService struct {
	apiKey string
	client *http.Client
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func NewOpenAIService(apiKey string) *OpenAIService {
	return &OpenAIService{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateQuestionFromContent generates a thoughtful question based on the content
func (s *OpenAIService) GenerateQuestionFromContent(content, title string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	// Create a prompt that encourages meaningful questions
	prompt := fmt.Sprintf(`Based on the following content, generate a single, thoughtful question that would help someone learn and understand the key concepts. The question should:

1. Test understanding of the main concepts
2. Be clear and specific
3. Encourage deep thinking
4. Be appropriate for a learning/study context

Title: %s
Content: %s

Generate only the question, nothing else:`, title, content)

	request := OpenAIRequest{
		Model: "gpt-4o-mini-2024-07-18",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   100,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI")
	}

	question := strings.TrimSpace(response.Choices[0].Message.Content)
	return question, nil
}

// GenerateLearningContent generates AI content for learning sessions
func (s *OpenAIService) GenerateLearningContent(ctx context.Context, aiRequest models.AIRequest) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	messages := s.buildPromptMessages(aiRequest)

	request := OpenAIRequest{
		Model:       "gpt-4o-mini-2024-07-18",
		Messages:    messages,
		MaxTokens:   1500,
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI")
	}

	content := strings.TrimSpace(response.Choices[0].Message.Content)
	return content, nil
}

// buildPromptMessages constructs the message array for OpenAI API
func (s *OpenAIService) buildPromptMessages(aiRequest models.AIRequest) []Message {
	messages := []Message{
		{Role: "system", Content: aiRequest.SystemPrompt},
	}

	// Add session summary as system message if available
	if aiRequest.SessionSummary != "" {
		messages = append(messages, Message{
			Role:    "system",
			Content: fmt.Sprintf("Session Progress Summary:\n%s", aiRequest.SessionSummary),
		})
	}

	// Add recent interactions (sliding window)
	for _, interaction := range aiRequest.RecentInteractions {
		messages = append(messages, Message{
			Role:    interaction.Role,
			Content: interaction.Content,
		})
	}

	// Add current user input
	if aiRequest.CurrentInput != "" {
		messages = append(messages, Message{
			Role:    "user",
			Content: aiRequest.CurrentInput,
		})
	}

	return messages
}

// GenerateSessionSummary creates a summary of a learning session's progress
func (s *OpenAIService) GenerateSessionSummary(ctx context.Context, interactions []models.LearningInteraction, conceptName string) (string, error) {
	if s.apiKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	// Build summary prompt
	summaryPrompt := fmt.Sprintf(
		"Summarize the user's progress on the concept '%s'. Include: key points covered, assessments passed/failed, areas of strength, areas needing review. Be concise and specific.",
		conceptName,
	)

	messages := []Message{
		{Role: "system", Content: "You are a learning progress summarizer. Create concise, helpful summaries of learning sessions."},
	}

	// Add interactions for context
	for _, interaction := range interactions {
		messages = append(messages, Message{
			Role:    interaction.Role,
			Content: interaction.Content,
		})
	}

	// Add summary request
	messages = append(messages, Message{
		Role:    "user",
		Content: summaryPrompt,
	})

	request := OpenAIRequest{
		Model:       "gpt-4o-mini-2024-07-18", // Use mini for summaries to save costs
		Messages:    messages,
		MaxTokens:   300,
		Temperature: 0.3, // Lower temperature for more focused summaries
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI")
	}

	summary := strings.TrimSpace(response.Choices[0].Message.Content)
	return summary, nil
}