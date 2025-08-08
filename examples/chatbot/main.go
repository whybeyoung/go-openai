package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
				Role:    openai.ChatMessageRoleUser,
				Content: "you are a helpful chatbot 你是谁",
			},
		},
		ExtraBody: map[string]interface{}{
			"bootstrap_host": "0.0.0.1",
			"bootstrap_port": 23,
			"bootstrap_room": 1111,
			"prefill_addr":   "0.0.0.1",
			"session_params": map[string]interface{}{
				"rid": "mgn00010282@dx196c8dc6f88ba44182",
			},
		},
	}
	stream, err := client.CreateChatCompletionStream(context.Background(), req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stream.Close()
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
			}
			return
		default:
		}
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if len(response.Choices) > 0 && response.Choices[0].Delta.Content != "" {
			fmt.Println(response.Choices[0].Delta.Content)
		}
	}
	//req.TopP = 0.1
	//fmt.Println("Conversation")
	//fmt.Println("---------------------")
	//fmt.Print("> ")
	//s := bufio.NewScanner(os.Stdin)
	//for s.Scan() {
	//	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
	//		Role:    openai.ChatMessageRoleUser,
	//		Content: s.Text(),
	//	})
	//	resp, err := client.CreateChatCompletion(context.Background(), req)
	//	if err != nil {
	//		fmt.Printf("ChatCompletion error: %v\n", err)
	//		continue
	//	}
	//	fmt.Printf("%s\n\n", resp.Choices[0].Message.Content)
	//	req.Messages = append(req.Messages, resp.Choices[0].Message)
	//	fmt.Print("> ")
	//}
}
