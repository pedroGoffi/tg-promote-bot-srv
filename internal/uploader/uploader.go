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

// PrepareUploadResponse represents the response from UploadThing's prepareUpload endpoint.
type PrepareUploadResponse struct {
	Files []struct {
		Key          string `json:"key"`
		UploadURL    string `json:"uploadUrl"`
		UploadedFile string `json:"uploadedFile"`
	} `json:"files"`
}

// UploadFile uploads a file to UploadThing and returns the file URL.
func UploadFile(fileName string, fileData io.Reader, fileSize int64) (string, error) {
	apiKey := config.GetUT_KEY()
	if apiKey == "" {
		return "", fmt.Errorf("UPLOADTHING_API_KEY is not set")
	}

	// Step 1: Prepare the Upload
	prepareUploadURL := "https://api.uploadthing.com/v6/prepareUpload"
	payload := map[string]interface{}{
		"callbackUrl":  "",
		"callbackSlug": "",
		"files": []map[string]interface{}{
			{
				"name":     fileName,
				"size":     fileSize,
				"customId": nil,
			},
		},
		"routeConfig": []string{"image"}, // Adjust based on your needs
		"metadata":    nil,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", prepareUploadURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("X-Uploadthing-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to prepare upload: %s", string(body))
	}

	var prepareResponse PrepareUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&prepareResponse); err != nil {
		return "", fmt.Errorf("failed to parse prepareUpload response: %w", err)
	}

	if len(prepareResponse.Files) == 0 {
		return "", fmt.Errorf("no upload URL returned")
	}

	uploadURL := prepareResponse.Files[0].UploadURL

	// Step 2: Upload the File
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, fileData); err != nil {
		return "", fmt.Errorf("failed to copy file data: %w", err)
	}
	writer.Close()

	req, err = http.NewRequest("PUT", uploadURL, &requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("file upload failed: %s", string(body))
	}

	return prepareResponse.Files[0].UploadedFile, nil
}
