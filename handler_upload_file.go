package main

import(
	"fmt"
	"os"
	"path/filepath"
	"log"

	"github.com/joho/godotenv"

	"github.com/Bention99/fin-planalyse/internal/openaiextract"
)

func uploadFile() {
	godotenv.Load(".env")

	apiKey := os.Getenv("OPENAI_API_KEY")
	promptKey := os.Getenv("PROMPTID")

	uploadPath, err := getSingleUploadFilePath("./uploads")
	if err != nil {
		log.Fatal(err)
	}

	fileID, err := openaiextract.UploadFile(apiKey, uploadPath)
	if err != nil {
		panic(err)
	}

	response, err := openaiextract.CreateResponse(apiKey, promptKey, fileID)
	if err != nil {
		panic(err)
	}

	fmt.Println(response)

	//os.Remove(uploadPath)
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