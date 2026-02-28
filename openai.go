package main

import(
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/Bention99/fin-planalyse/internal/openaiextract"
)

func uploadFile(catsIncome []string, catsExpense []string) error {
	godotenv.Load(".env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	promptKey := os.Getenv("PROMPTID")

	uploadPath, err := getSingleUploadFilePath("./uploads")
	if err != nil {
		return err
	}

	fileID, err := openaiextract.UploadFile(apiKey, uploadPath)
	if err != nil {
		return err
	}

	response, err := openaiextract.CreateResponse(apiKey, promptKey, fileID, catsIncome, catsExpense)
	if err != nil {
		return err
	}

	fmt.Println(response)

	//os.Remove(uploadPath)
	return nil
}

func getSingleUploadFilePath(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var files []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, e)
		}
	}

	if len(files) == 0 {
		return "", fmt.Errorf("no files found in %s", dir)
	}
	if len(files) > 1 {
		return "", fmt.Errorf("more than one file found in %s", dir)
	}

	return filepath.Join(dir, files[0].Name()), nil
}