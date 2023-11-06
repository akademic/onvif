package networking

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

// SendSoap send soap message
func SendSoap(ctx context.Context, httpClient *http.Client, endpoint, message string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBufferString(message))
	if err != nil {
		return nil, fmt.Errorf("NewRequestWithContext: %w", err)
	}

	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")
	resp, err := httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("Post: %w", err)
	}

	return resp, nil
}
