package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	apiKey = os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}
}

func main() {
	repoRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatalf("Gitリポジトリのルートディレクトリを取得できませんでした: %v", err)
		return
	}

	repoRootPath := strings.TrimSpace(string(repoRoot))
	if err := os.Chdir(repoRootPath); err != nil {
		log.Fatalf("Gitリポジトリのルートディレクトリに移動できませんでした: %v", err)
		return
	}

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
	prompt := fmt.Sprintf("以下のGitの変更差分を見て、適切なコミットメッセージを簡潔な形で、また日本語で提案してください。\n\n%s", diff)
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

	cmd := exec.Command("curl", "-s", "https://api.openai.com/v1/chat/completions",
		"-H", "Content-Type: application/json",
		"-H", fmt.Sprintf("Authorization: Bearer %s", apiKey),
		"-d", string(payload))

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatalf("OpenAI APIのリクエストに失敗しました: %v", err)
		return
	}

	var response OpenAIResponse
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		log.Fatalf("レスポンスのJSONデコードに失敗しました: %v", err)
		return
	}

	if len(response.Choices) == 0 {
		fmt.Println("No response from OpenAI")
		return
	}

	commitMessage := response.Choices[0].Message.Content
	fmt.Println(commitMessage)
}
