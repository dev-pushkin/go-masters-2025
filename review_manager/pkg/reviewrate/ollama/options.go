package ollama

import "net/http"

type options struct {
	endpoint    string
	model       string
	systemPromt string
	httpClient  *http.Client
}
type Option func(*options)

func WithEndpoint(endpoint string) Option {
	return func(o *options) {
		o.endpoint = endpoint
	}
}
func WithModel(model string) Option {
	return func(o *options) {
		o.model = model
	}
}
func WithSystemPrompt(systemPrompt string) Option {
	return func(o *options) {
		o.systemPromt = systemPrompt
	}
}
func WithHTTPClient(client *http.Client) Option {
	return func(o *options) {
		o.httpClient = client
	}
}
