package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/motty93/ai-commits-message/i18n"
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
var apiUrl = "https://api.openai.com/v1/chat/completions"

func init() {
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		return
	}

	i18n.Init()
}

func main() {
	diff, err := exec.Command("git", "diff", "--cached").Output()
	if err != nil {
		// log.Fatalf("ステージングエリアの差分を取得できませんでした: %v", err)
		fmt.Println("") // 空文字を返して終了
		return
	}

	if len(diff) == 0 {
		fmt.Println("") // 空文字を返して終了
		return
	}

	content := i18n.GetText("content")
	prompt := fmt.Sprintf(`%s%s`, i18n.GetText("prompt"), diff)
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
		log.Fatalf(i18n.GetText("encode_error"), err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf(i18n.GetText("create_request_failed"), err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf(i18n.GetText("post_request_failed"), err)
		return
	}
	defer resp.Body.Close()

	var response OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatalf(i18n.GetText("response_decode_error"), err)
		return
	}

	if len(response.Choices) == 0 {
		fmt.Println(i18n.GetText("no_response"))
		return
	}

	commitMessage := response.Choices[0].Message.Content
	fmt.Println(strings.Trim(commitMessage, "\""))
}
