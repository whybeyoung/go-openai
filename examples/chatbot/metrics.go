package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	// ...
)

// Metric 表示一个完整的指标，包括名称、类型及内容
type Metric struct {
	Name    string
	Type    string
	Content []ContentItem
}
type ContentItem struct {
	Key   string
	Value float64 // 改成 float64
}

func ParseMetrics(text string) map[string]*Metric {
	lines := strings.Split(text, "\n")
	metrics := make(map[string]*Metric, len(lines))
	var currentMetric *Metric

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "# HELP"):
			continue

		case strings.HasPrefix(line, "# TYPE"):
			parts := strings.Fields(line)
			if len(parts) < 4 {
				continue
			}

			currentMetric = &Metric{
				Name:    parts[2],
				Type:    parts[3],
				Content: []ContentItem{},
			}
			metrics[parts[2]] = currentMetric

		default:
			if currentMetric == nil {
				continue
			}
			keyValue := strings.SplitN(line, " ", 2)
			if len(keyValue) == 2 {
				val, err := strconv.ParseFloat(keyValue[1], 64)
				if err != nil {
					// 解析失败你可以选择跳过或赋0
					val = 0
				}
				currentMetric.Content = append(currentMetric.Content, ContentItem{
					Key:   keyValue[0],
					Value: val,
				})
			}
		}
	}

	return metrics
}

func main() {
	sampleText := `
# HELP sglang:num_decode_prealloc_queue_reqs The number of requests in the decode prealloc queue.
# TYPE sglang:num_decode_prealloc_queue_reqs gauge
sglang:num_decode_prealloc_queue_reqs{engine_type="unified",model_name="/work/models/"} 0.0
# HELP sglang:num_decode_transfer_queue_reqs The number of requests in the decode transfer queue.
# TYPE sglang:num_decode_transfer_queue_reqs gauge
sglang:num_decode_transfer_queue_reqs{engine_type="unified",model_name="/work/models/"} 0.0
# HELP sglang:prompt_tokens_total Number of prefill tokens processed.
# TYPE sglang:prompt_tokens_total counter
sglang:prompt_tokens_total{model_name="/work/models/"} 1.282888e+06
# HELP sglang:generation_tokens_total Number of generation tokens processed.
# TYPE sglang:generation_tokens_total counter
sglang:generation_tokens_total{model_name="/work/models/"} 337.0
# HELP sglang:num_requests_total Number of requests processed.
# TYPE sglang:num_requests_total counter
sglang:num_requests_total{model_name="/work/models/"} 337.0
# HELP sglang:num_so_requests_total Number of structured output requests processed.
# TYPE sglang:num_so_requests_total counter
sglang:num_so_requests_total{model_name="/work/models/"} 5.0
# HELP sglang:time_to_first_token_seconds Histogram of time to first token in seconds.
# TYPE sglang:time_to_first_token_seconds histogram
`

	got := ParseMetrics(sampleText)

	want := map[string]*Metric{
		"sglang_prompt_tokens_total": {
			Name: "sglang_prompt_tokens_total",
			Type: "counter",
			Content: []ContentItem{
				{Key: `sglang_prompt_tokens_total{model="llama"}`, Value: 123456},
				{Key: `sglang_prompt_tokens_total{model="mistral"}`, Value: 654321},
			},
		},
		"http_requests_total": {
			Name: "http_requests_total",
			Type: "counter",
			Content: []ContentItem{
				{Key: `http_requests_total{method="get"}`, Value: 100},
				{Key: `http_requests_total{method="post"}`, Value: 200},
			},
		},
	}
	fmt.Println((got["sglang:prompt_tokens_total"].Content[0].Value - 100000) / 3)

	if !reflect.DeepEqual(got, want) {
		fmt.Println("err")
	}
}
