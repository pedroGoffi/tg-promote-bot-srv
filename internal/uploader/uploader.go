package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadFile uploads a file to UploadThing using an io.Reader and returns the file URL.
func UploadFile(fileName string, fileData io.Reader) (string, error) {
	apiKey := os.Getenv("UPLOADTHING_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("UPLOADTHING_API_KEY is not set")
	}

	// Create a buffer for multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create the file field in the form
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	// Copy fileData to the multipart field
	if _, err = io.Copy(part, fileData); err != nil {
		return "", err
	}

	// Close the writer to finalize the form data
	writer.Close()

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://uploadthing.com/api/upload", &requestBody)
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload file: %s", string(body))
	}

	// Parse the response JSON
	var result struct {
		Url string `json:"fileUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Url, nil
}
