package uploader

import (
	"bot-manager/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

// ImgBB API response struct
type ImgBBResponse struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
	Success bool `json:"success"`
}

func AlreadyUploaded(url string) bool {
	resp, err := http.Head(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// UploadImage uploads an image to ImgBB and returns its URL
func UploadImage(fileName string, fileData io.Reader, expirationSeconds int) (string, error) {
	var checkExistsUrl string = fmt.Sprintf("https://i.ibb.co/RTKdyNHj/%s", fileName)
	if AlreadyUploaded(checkExistsUrl) {
		return checkExistsUrl, nil
	}

	apiKey := config.GetUploadKey()
	if apiKey == "" {
		return "", fmt.Errorf("IMGBB_API_KEY is not set")
	}

	// Create form-data payload
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create file field
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, fileData); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}

	// Add API key as form field
	_ = writer.WriteField("key", apiKey)

	// Optional: Set the image name
	_ = writer.WriteField("name", fileName)

	// Optional: Set expiration time if provided
	if expirationSeconds > 0 {
		_ = writer.WriteField("expiration", strconv.Itoa(expirationSeconds))
	}

	// Close writer
	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.imgbb.com/1/upload", &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON response
	var result ImgBBResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Check if upload was successful
	if !result.Success {
		return "", fmt.Errorf("upload failed: %s", string(body))
	}

	return result.Data.URL, nil
}
