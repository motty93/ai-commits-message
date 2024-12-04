package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

var apiKey string

func init() {
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}
}

func main() {
	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		log.Fatalf("ステージングエリアの差分を取得できませんでした: %v", err)
		return
	}

	if len(diff) == 0 {
		fmt.Println("No changes to commit")
		return
	}

	content := "あなたは優れたソフトウェアエンジニアです。"
	prompt := fmt.Sprintf(`以下の Git の変更差分を見て、以下の形式で適切なコミットメッセージを提案してください:
		1. プレフィックス (例: "feat:", "fix:", "refactor:", など) を含める。
		2. メッセージは簡潔でわかりやすい形にする。
		3. 可能な限り日本語で記述する。
		Git の変更差分:\n\n%s`, diff)
	request := OpenAIRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "system", Content: content},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.7,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		log.Fatalf("リクエストのJSONエンコードに失敗しました: %v", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("HTTPリクエストの作成に失敗しました：%v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTPリクエストの送信に失敗しました: %v", err)
		return
	}
	defer resp.Body.Close()

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf("レスポンスのJSONデコードに失敗しました: %v", err)
		return
	}

	if len(response.Choices) == 0 {
		fmt.Println("No response from OpenAI")
		return
	}

	commitMessage := response.Choices[0].Message.Content
	fmt.Println(strings.Trim(commitMessage, "\""))
}
