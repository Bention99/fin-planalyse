package openaiextract

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type ExtractedResponse struct {
	Transactions []ExtractedTransaction `json:"transactions"`
}

type ExtractedTransaction struct {
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Category    string  `json:"category"`
	Necessity   string  `json:"necessity"`
}

func UploadFile(apiKey, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	if err := writer.WriteField("purpose", "assistants"); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/files", &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ID, nil
}

func CreateResponse(apiKey string, promptKey string, fileID string) (ExtractedResponse, error) {

	payload := map[string]interface{}{
		"model": "gpt-4.1",
		"prompt": map[string]interface{}{
			"id": promptKey,
			"version": "8",
		},
		"input": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "input_text",
						"text": "Extract all transactions from the provided bank statement and return the result as JSON.",
					},
					{
						"type": "input_file",
						"file_id": fileID,
					},
				},
			},
		},
		"text": map[string]interface{}{
			"format": map[string]interface{}{
				"type": "json_object",
			},
		},
		"temperature": 0.0,
		"max_output_tokens": 3000,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return ExtractedResponse{}, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(jsonBody))
	if err != nil {
		return ExtractedResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ExtractedResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)


	var apiResp struct {
		Output []struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return ExtractedResponse{}, err
	}

	jsonString := apiResp.Output[0].Content[0].Text

	var parsed ExtractedResponse
	if err := json.Unmarshal([]byte(jsonString), &parsed); err != nil {
		return ExtractedResponse{}, err
	}
	return parsed, nil
}