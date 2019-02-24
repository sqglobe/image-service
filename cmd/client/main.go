package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func attachFile(path, name string, writer *multipart.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(name, filepath.Base(path))
	if err != nil {
		return err
	}
	if _, err = io.Copy(part, file); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Needed 2 params: <image path> <watermark>")
		os.Exit(1)
	}

	imagePath := os.Args[1]
	watermarkPath := os.Args[2]

	fmt.Printf("Used image: %s, watermark: %s\n", imagePath, watermarkPath)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := attachFile(imagePath, "image", writer); err != nil {
		fmt.Printf("Failed create part for image. Error: %s\n", err)
		os.Exit(1)
	}

	if err := attachFile(watermarkPath, "watermark", writer); err != nil {
		fmt.Printf("Failed create part for watermark. Error: %s\n", err)
		os.Exit(1)
	}

	if err := writer.Close(); err != nil {
		fmt.Printf("Failed create message. Error: %s\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", "http://localhost:3210/watermark", body)
	if err != nil {
		fmt.Printf("Failed create request. Error: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Failed send request. Error: %s\n", err)
		os.Exit(1)
	} else {
		respBody := &bytes.Buffer{}
		_, err := respBody.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(respBody)
	}

}
