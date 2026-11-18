package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type AIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func AskAIBrain(prompt string) (string, error) {
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("AI_API_KEY environment variable set nahi hai")
	}

	url := "https://api.openai.com/v1/chat/completions"

	reqBody := AIRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{Role: "system", Content: "You are QYVORA AI, an expert cybersecurity and offensive security assistant."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	var aiResp AIResponse
	err = json.NewDecoder(resp.Body).Decode(&aiResp)
	if err != nil {
		return "", err
	}

	if len(aiResp.Choices) > 0 {
		return aiResp.Choices[0].Message.Content, nil
	}

	return "Koi response nahi mila.", nil
}
