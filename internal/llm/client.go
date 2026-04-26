package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	url   string
	model string
	http  *http.Client
}

func NewClient(url, model string) *Client {
	return &Client{
		url: url,
		model: model,
		http: &http.Client{},
	}
}


func (c *Client) Generate(prompt string) (string, error) {
	req := GenerateRequest {
		Model: c.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Post(c.url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("LLM не отвечает: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM вернула статус: %d", resp.StatusCode)
	}

	var res GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.Response, nil
}