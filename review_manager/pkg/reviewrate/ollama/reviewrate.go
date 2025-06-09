package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go_course_master/review_manager/pkg/reviewrate"
)

type ReviewRate struct {
	c *http.Client

	endpoint    string
	model       string
	systemPromt string
}

func New(opts ...Option) *ReviewRate {
	options := &options{
		endpoint:    "http://localhost:11434/api/generate",
		model:       "qwen2.5:3b",
		systemPromt: "",
		httpClient:  defaultHttpClient,
	}
	for _, opt := range opts {
		opt(options)
	}

	return &ReviewRate{
		c:           options.httpClient,
		endpoint:    options.endpoint,
		model:       options.model,
		systemPromt: options.systemPromt,
	}
}

func (r *ReviewRate) Rate(ctx context.Context, review string) (rate reviewrate.RateLevel, err error) {
	reqData := struct {
		Model  string `json:"model"`
		Prompt string `json:"prompt"`
		System string `json:"system,omitempty"`
	}{
		Model:  r.model,
		Prompt: review,
		System: r.systemPromt,
	}

	body, err := json.Marshal(reqData)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.endpoint, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	resp, err := r.c.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("LLM API returned status %d", resp.StatusCode)
	}

	response := bytes.NewBuffer(nil)
	err = parseAndWriteLlmStream(resp, response)
	if err != nil {
		return 0, fmt.Errorf("parse and write LLM stream: %w", err)
	}

	// For simplicity, let's assume the response is a single rate level
	str := response.String()
	if str == "" {
		return 0, fmt.Errorf("empty response from LLM")
	}

	var level int
	// Convert the response string to an integer rate level
	_, err = fmt.Sscanf(str, "%d", &level)
	if err != nil {
		return 0, fmt.Errorf("failed to parse rate level from response: %w", err)
	}
	return reviewrate.RateLevel(level), nil
}

type ResponseChunk struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func parseAndWriteLlmStream(in *http.Response, out io.Writer) error {
	defer in.Body.Close()
	resBody := in.Body

	decoder := json.NewDecoder(resBody)
	for {
		var chunk ResponseChunk
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("decode error: %w", err)
		}
		if chunk.Response != "" {
			// Write the response to the output writer
			if _, err := fmt.Fprint(out, chunk.Response); err != nil {
				return fmt.Errorf("write error: %w", err)
			}
		}

		if chunk.Done {
			break
		}
	}

	return nil
}

func (r *ReviewRate) Stop(ctx context.Context) error {
	r.c.CloseIdleConnections()
	return nil
}
