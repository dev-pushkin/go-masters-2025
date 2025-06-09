package ollama

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/go_course_master/review_manager/pkg/reviewrate"
)

// mockRoundTripper заменяет http.Client.Transport для подмены ответов
type mockRoundTripper struct {
	response string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body := `{"response": "` + m.response + `", "done": true}`
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}, nil
}

func TestReviewRate_Rate(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     reviewrate.RateLevel
		wantErr  bool
	}{
		{"simple digit", "4", 4, false},
		{"prefixed", "3", 3, false},
		{"textual", "I'd rate this a 5 because it's awesome.", 0, true},
		{"empty", "", 0, true},
		{"invalid", "meh", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rater := &ReviewRate{
				c:           &http.Client{Transport: &mockRoundTripper{response: tt.response}},
				endpoint:    "http://mocked/endpoint",
				model:       "mock-model",
				systemPromt: "",
			}

			got, err := rater.Rate(context.Background(), "This is a test review.")
			if (err != nil) != tt.wantErr {
				t.Errorf("Rate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Rate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
