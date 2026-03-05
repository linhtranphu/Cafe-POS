package printbridge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// HTTPClient handles communication with Print Bridge via HTTP
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// PrintRequest represents a print request to Print Bridge
type PrintRequest struct {
	JobID       string `json:"jobId"`
	Content     string `json:"content"`
	PrinterIP   string `json:"printerIP"`
	PrinterPort int    `json:"printerPort"`
	Type        string `json:"type"` // "bill" or "label"
}

// PrintResponse represents response from Print Bridge
type PrintResponse struct {
	Success   bool   `json:"success"`
	JobID     string `json:"jobId"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
}

// NewHTTPClient creates a new Print Bridge HTTP client
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendPrintJob sends a print job to Print Bridge via HTTP
func (c *HTTPClient) SendPrintJob(req PrintRequest) (*PrintResponse, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("print bridge URL not configured")
	}

	// Marshal request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/print", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// Send request
	log.Printf("[PrintBridge] Sending print job %s to %s", req.JobID, url)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var printResp PrintResponse
	if err := json.Unmarshal(body, &printResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return &printResp, fmt.Errorf("print bridge returned error: %s", printResp.Error)
	}

	if !printResp.Success {
		return &printResp, fmt.Errorf("print job failed: %s", printResp.Message)
	}

	log.Printf("[PrintBridge] Print job %s sent successfully", req.JobID)
	return &printResp, nil
}

// TestConnection tests connection to Print Bridge
func (c *HTTPClient) TestConnection() error {
	if c.baseURL == "" {
		return fmt.Errorf("print bridge URL not configured")
	}

	url := fmt.Sprintf("%s/health", c.baseURL)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	log.Printf("[PrintBridge] Connection test successful: %s", c.baseURL)
	return nil
}
