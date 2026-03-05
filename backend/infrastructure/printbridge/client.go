package printbridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with local print bridge
type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// RenderRequest represents the request to render HTML
type RenderRequest struct {
	HTML         string                 `json:"html"`
	Width        int                    `json:"width,omitempty"`
	ShopSettings map[string]interface{} `json:"shopSettings,omitempty"`
}

// RenderResponse represents the response from print bridge
type RenderResponse struct {
	Success      bool   `json:"success"`
	ESCPOSData   string `json:"escposData"`   // base64 encoded
	PreviewImage string `json:"previewImage"` // base64 encoded
	Stats        struct {
		Duration   int `json:"duration"`
		ImageSize  int `json:"imageSize"`
		ESCPOSSize int `json:"escposSize"`
	} `json:"stats"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
}

// NewClient creates a new print bridge client
func NewClient(baseURL string, timeout time.Duration) *Client {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// RenderHTMLToESCPOS sends HTML to print bridge for rendering
func (c *Client) RenderHTMLToESCPOS(ctx context.Context, html string, width int) ([]byte, error) {
	// Prepare request
	reqBody := RenderRequest{
		HTML:  html,
		Width: width,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/render-html", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to print bridge: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var renderResp RenderResponse
	if err := json.Unmarshal(body, &renderResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if !renderResp.Success {
		return nil, fmt.Errorf("print bridge error: %s", renderResp.Error)
	}

	// Decode base64 ESC/POS data
	escposData, err := base64.StdEncoding.DecodeString(renderResp.ESCPOSData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ESC/POS data: %w", err)
	}

	return escposData, nil
}

// TestConnection tests connection to print bridge
func (c *Client) TestConnection(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to print bridge: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("print bridge returned status %d", resp.StatusCode)
	}

	return nil
}

// IsAvailable checks if print bridge is available
func (c *Client) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := c.TestConnection(ctx)
	return err == nil
}

// RenderAndPrint sends HTML to print bridge for rendering and printing
func (c *Client) RenderAndPrint(ctx context.Context, html string, printerIP string, printerPort int, width int) error {
	if printerPort == 0 {
		printerPort = 9100
	}
	if width == 0 {
		width = 576
	}

	// Prepare request
	reqBody := map[string]interface{}{
		"html":        html,
		"width":       width,
		"printerIP":   printerIP,
		"printerPort": printerPort,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/render-and-print", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to print bridge: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var printResp struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
		Message string `json:"message,omitempty"`
	}
	if err := json.Unmarshal(body, &printResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if !printResp.Success {
		return fmt.Errorf("print bridge error: %s", printResp.Error)
	}

	return nil
}

// PrintESCPOS sends ESC/POS data to print bridge for printing
func (c *Client) PrintESCPOS(ctx context.Context, escposData []byte, printerIP string, printerPort int) error {
	if printerPort == 0 {
		printerPort = 9100
	}

	// Prepare request
	reqBody := map[string]interface{}{
		"content":     base64.StdEncoding.EncodeToString(escposData),
		"printerIP":   printerIP,
		"printerPort": printerPort,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/print", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to print bridge: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var printResp struct {
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
		Message string `json:"message,omitempty"`
	}
	if err := json.Unmarshal(body, &printResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if !printResp.Success {
		return fmt.Errorf("print bridge error: %s", printResp.Error)
	}

	return nil
}
