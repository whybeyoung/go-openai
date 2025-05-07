package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/whybeyoung/go-openai"
)

func main() {
	config := openai.DefaultConfig("sk-Kj7yJbMsBJXWIUxBD1C5B7446bD34a99912702662eB31131")
	//config.BaseURL = "http://maas-api.cn-huabei-1.xf-yun.com/v1"
	config.BaseURL = "http://36.138.165.10:30822/v1"

	config.HTTPClient = &http.Client{
		Timeout: 2 * time.Minute,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}
	client := openai.NewClientWithConfig(config)

	req := openai.ChatCompletionRequest{
		Model: "xdeepseekv3",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "you are a helpful chatbot",
			},
		},
		ExtraBody: map[string]interface{}{
			"bootstrap_host": "0.0.0.1",
			"bootstrap_port": 23,
			"bootstrap_room": 1111,
			"prefill_addr":   "0.0.0.1",
		},
	}

	req.TopP = 0.1
	fmt.Println("Conversation")
	fmt.Println("---------------------")
	fmt.Print("> ")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: s.Text(),
		})
		resp, err := client.CreateChatCompletion(context.Background(), req)
		if err != nil {
			fmt.Printf("ChatCompletion error: %v\n", err)
			continue
		}
		fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
		req.Messages = append(req.Messages, resp.Choices[0].Message)
		fmt.Print("> ")
	}
}
