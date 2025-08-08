package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Config struct {
	BootstrapIP   string `json:"bootstrap_ip"`
	BootstrapPort int    `json:"bootstrap_port"`
	BootstrapRoom int64  `json:"bootstrap_room"`
	PrefillAddr   string `json:"prefill_addr"`
}

func main() {
	a := "eyJib290c3RyYXBfaXAiOiIyNi41LjI3LjIzOCIsImJvb3RzdHJhcF9wb3J0Ijo4OTk4LCJib290c3RyYXBfcm9vbSI6MTg3NDkzNTUyODg4NTg0NjAxNywicHJlZmlsbF9hZGRyIjoiMjYuNS4yNy4yMzg6MzAwMDAifQ=="
	a = "eyJib290c3RyYXBfaXAiOiIyNi41LjI3LjIzOCIsImJvb3RzdHJhcF9wb3J0Ijo4OTk4LCJib290c3RyYXBfcm9vbSI6MTg3NDkzNTUyODg4NTg0NjAxNiwicHJlZmlsbF9hZGRyIjoiMjYuNS4yNy4yMzg6MzAwMDAifQ=="
	// 解码 base64
	decoded, err := base64.StdEncoding.DecodeString(a)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)
		return
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(decoded, &config); err != nil {
		fmt.Printf("JSON unmarshal error: %v\n", err)
		return
	}

	// 打印解析结果
	fmt.Printf("解析结果:\n")
	fmt.Printf("Bootstrap IP: %s\n", config.BootstrapIP)
	fmt.Printf("Bootstrap Port: %d\n", config.BootstrapPort)
	fmt.Printf("Bootstrap Room: %d\n", config.BootstrapRoom)
	fmt.Printf("Prefill Addr: %s\n", config.PrefillAddr)

	// 打印原始 JSON
	fmt.Printf("\n原始 JSON:\n%s\n", string(decoded))
}
