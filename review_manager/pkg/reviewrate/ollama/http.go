package ollama

import (
	"net/http"
	"time"
)

var (
	defaultTimeout = time.Duration(30) * time.Second // Default timeout for HTTP requests

	defaultHttpClient *http.Client = &http.Client{
		Timeout: defaultTimeout, // Default timeout, can be set to a specific value if needed
	}
)
