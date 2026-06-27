package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const deepseekBaseURL = "https://api.deepseek.com/v1/chat/completions"

type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type chatRequest struct {
	Model     string    `json:"model"`
	Messages  []message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func NewClient() *Client {
	model := os.Getenv("DEEPSEEK_MODEL")
	if model == "" {
		model = "deepseek-chat"
	}
	return &Client{
		apiKey:     os.Getenv("DEEPSEEK_API_KEY"),
		model:      model,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *Client) chat(messages []message, maxTokens int) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("DEEPSEEK_API_KEY not set")
	}

	payload := chatRequest{
		Model:     c.model,
		Messages:  messages,
		MaxTokens: maxTokens,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", deepseekBaseURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// CheckWriting sends a business writing submission for AI feedback
func (c *Client) CheckWriting(exerciseType, context, prompt, userContent string) (map[string]interface{}, error) {
	systemPrompt := `You are a Business English coach specializing in IT/tech professional communication.
Analyze the writing submission and return a JSON object with these fields:
- score: integer 0-100
- overall: brief overall assessment (1-2 sentences)
- grammar: list of grammar issues found (array of strings, empty if none)
- tone: assessment of professional tone (string)
- vocabulary: suggestions for better business vocabulary (array of strings)
- improved_version: a short improved version of their text
Keep feedback constructive and specific to business/IT context.`

	userMsg := fmt.Sprintf(`Exercise type: %s
Context: %s
Prompt: %s

Student's submission:
%s

Return only valid JSON.`, exerciseType, context, prompt, userContent)

	response, err := c.chat([]message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMsg},
	}, 1000)
	if err != nil {
		return nil, err
	}

	// Strip markdown code fences if present
	cleaned := strings.TrimSpace(response)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	cleaned = strings.TrimSuffix(strings.TrimSpace(cleaned), "```")
	cleaned = strings.TrimSpace(cleaned)

	// Parse JSON response
	var feedback map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &feedback); err != nil {
		// Return raw response if not valid JSON
		return map[string]interface{}{
			"overall": response,
			"score":   50,
		}, nil
	}

	return feedback, nil
}

// GenerateSentence generates a random Indonesian business/technical sentence for the user to translate
func (c *Client) GenerateSentence() (map[string]interface{}, error) {
	systemPrompt := `You are a Business English coach for IT professionals and software engineers. Generate an Indonesian sentence commonly used in a professional software engineering or IT work environment.
The sentence should reflect real scenarios such as: sprint planning, pull request reviews, code reviews, system design discussions, incident response, deployment pipelines, API documentation, technical interviews, client demos, engineering meetings, architecture decisions, agile ceremonies, bug triage, or technical writing.
Use vocabulary and phrasing that software engineers actually use at work — including technical terms like: sprint, backlog, deployment, refactor, bottleneck, latency, scalability, trade-off, on-call, rollback, SLA, tech debt, MVP, stakeholder, etc.
Return ONLY a valid JSON object with these fields:
- indonesian_sentence: the Indonesian sentence (string, 1 sentence, professional software engineering/IT context)
- challenge_id: a short unique slug, e.g. "sprint-001" (string)
- topic: a very short technical topic label, e.g. "sprint planning", "code review", "incident response", "system design", "deployment", "API", "architecture", "agile", "client demo", "tech debt" (string)
- correct_answer: the correct natural Business English translation using appropriate technical vocabulary (string)
Do NOT include any explanation, markdown, or extra text. Return only the JSON object.`

	response, err := c.chat([]message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: "Generate a random Indonesian business/software engineering sentence for me to translate into Business English."},
	}, 300)
	if err != nil {
		return nil, err
	}

	// Strip markdown code fences if present
	cleaned := strings.TrimSpace(response)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	cleaned = strings.TrimSuffix(strings.TrimSpace(cleaned), "```")
	cleaned = strings.TrimSpace(cleaned)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse sentence generation response: %w", err)
	}
	return result, nil
}

// CheckTranslation checks a user's English translation of an Indonesian business/technical sentence
func (c *Client) CheckTranslation(indonesianSentence, correctAnswer, userAnswer string) (map[string]interface{}, error) {
	systemPrompt := `You are a Business English coach for IT professionals and software engineers evaluating a student's translation.
Given an Indonesian software engineering/IT sentence, the correct Business English translation, and the student's translation, provide a thorough evaluation.

Focus on:
1. Meaning accuracy — does the translation convey the correct technical meaning?
2. Business English tone — is the phrasing professional and appropriate for a workplace?
3. Technical vocabulary — are the correct technical terms used (e.g. "deploy" not "upload", "refactor" not "rewrite", "pull request" not "code submission")?
4. Natural phrasing — does it sound like something a native English-speaking engineer would actually say?

Return ONLY a valid JSON object with these fields:
- is_correct: true if the student's translation conveys the correct meaning with acceptable Business English (boolean)
- correct_answer: the best natural Business English translation (string)
- explanation: a clear and detailed explanation of why the correct answer is phrased that way — explain the business context, the technical term choice, and the tone. Make it educational so the student understands the reasoning, not just the answer. Use **double asterisks** to highlight key terms or phrases. Separate distinct points with newlines. (string, 2-4 sentences or bullet points starting with "- ")
- corrections: if the student's answer has issues, clearly state what was wrong and why. Use this format: start each issue with "- " on a new line, bold the wrong phrase with **wrong phrase**, then explain the better alternative with reasoning. Example: "- **deploy the code** is too informal; use **ship to production** or **deploy to production** instead." If fully correct, return empty string. (string)
Do NOT include any explanation, markdown, or extra text outside the JSON. Return only the JSON object.`

	userMsg := fmt.Sprintf(`Indonesian sentence: "%s"
Correct English translation: "%s"
Student's translation: "%s"

Evaluate the student's translation.`, indonesianSentence, correctAnswer, userAnswer)

	response, err := c.chat([]message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userMsg},
	}, 800)
	if err != nil {
		return nil, err
	}

	// Strip markdown code fences if present
	cleaned := strings.TrimSpace(response)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
	}
	cleaned = strings.TrimSuffix(strings.TrimSpace(cleaned), "```")
	cleaned = strings.TrimSpace(cleaned)

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return map[string]interface{}{
			"is_correct":     false,
			"correct_answer": correctAnswer,
			"explanation":    response,
			"corrections":    "",
		}, nil
	}
	return result, nil
}

// RoleplayOpen generates the AI's opening message for a scenario
func (c *Client) RoleplayOpen(systemPrompt, context, aiRole string) (string, error) {
	messages := []message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("Context: %s\nYou are playing: %s\nStart the conversation naturally.", context, aiRole)},
	}
	return c.chat(messages, 300)
}

// RoleplayReply generates the AI's reply in an ongoing roleplay session
func (c *Client) RoleplayReply(systemPrompt, context, aiRole string, history []map[string]string) (string, error) {
	messages := []message{
		{Role: "system", Content: fmt.Sprintf("%s\nContext: %s\nYou are playing: %s\nKeep responses realistic and professional, 2-4 sentences.", systemPrompt, context, aiRole)},
	}
	for _, h := range history {
		messages = append(messages, message{Role: h["role"], Content: h["content"]})
	}
	return c.chat(messages, 400)
}
