package openaiextract

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"fmt"
	"time"
	"strings"
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

func CreateResponse(apiKey string, promptKey string, fileID string, incomeCategories []string, expenseCategories []string) (ExtractedResponse, error) {
	incomeStr := strings.Join(incomeCategories, ", ")
	expenseStr := strings.Join(expenseCategories, ", ")

	payload := map[string]interface{}{
	"model": "gpt-4.1",
	"store": false,
	"prompt": map[string]interface{}{
		"id":      promptKey,
		"version": "9",
	},
	"input": []map[string]interface{}{
		{
			"role": "user",
			"content": []map[string]interface{}{
				{
					"type": "input_text",
					"text": fmt.Sprintf(
						"Extract all transactions from the provided bank statement and return JSON.\n\n"+
							"STRICT CATEGORY RULE:\n"+
							"- The category MUST be exactly one of the predefined categories below.\n"+
							"- Do NOT invent new categories.\n"+
							"- If unsure, choose the closest matching category from the list.\n\n"+
							"Predefined income categories:\n%s\n\n"+
							"Predefined expense categories:\n%s\n\n"+
							"Return only valid JSON following the required schema.",
						incomeStr,
						expenseStr,
					),
				},
				{
					"type":    "input_file",
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
	"temperature":       0.0,
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return ExtractedResponse{}, fmt.Errorf("openai error: status=%d body=%s", resp.StatusCode, string(bodyBytes))
	}

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

	if len(apiResp.Output) == 0 {
		return ExtractedResponse{}, fmt.Errorf("no output returned from OpenAI")
	}

	if len(apiResp.Output[0].Content) == 0 {
		return ExtractedResponse{}, fmt.Errorf("no content in OpenAI output")
	}

	jsonString := apiResp.Output[0].Content[0].Text

	var parsed ExtractedResponse
	if err := json.Unmarshal([]byte(jsonString), &parsed); err != nil {
		return ExtractedResponse{}, err
	}

	err = DeleteFile(apiKey, fileID)
	if err != nil {
		time.Sleep(5 * time.Second)

		err = DeleteFile(apiKey, fileID)
		if err != nil {
			fmt.Printf("Couldn't delete file after retry: %v\n", err)
		}
	}

	return parsed, nil
}

func DeleteFile(apiKey, fileID string) error {
	req, err := http.NewRequest(
		"DELETE",
		"https://api.openai.com/v1/files/"+fileID,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete file: %s", string(body))
	}

	return nil
}