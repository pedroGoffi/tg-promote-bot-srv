package uploader

import (
	"bot-manager/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// UploadResponse represents the response from UploadThing
type UploadResponse struct {
	Files []struct {
		Url string `json:"url"`
	} `json:"files"`
}

// UploadFile uploads a file to UploadThing and returns the file URL.
func UploadFile(fileName string, fileData io.Reader) (string, error) {
	apiKey := config.GetUT_KEY() // Keep using your function to get API key
	if apiKey == "" {
		return "", fmt.Errorf("UPLOADTHING_API_KEY is not set")
	}

	// Prepare multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Attach file
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, fileData); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	// Close writer to finalize the form
	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://uploadthing.com/api/upload", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upload failed: %s", string(body))
	}

	// Parse JSON response
	var result UploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Ensure file URL is returned
	if len(result.Files) == 0 {
		return "", fmt.Errorf("no files returned from UploadThing")
	}

	return result.Files[0].Url, nil
}
