package api

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	progressbar "github.com/schollz/progressbar/v3"
)

func DownloadFile(url, destDir, filename string) (string, error) {
	token := os.Getenv("CIVIT_API_KEY")
	log.Printf("Downloading %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	os.MkdirAll(destDir, os.ModePerm)
	fullPath := filepath.Join(destDir, filename)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}

	out, err := os.Create(absPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	bar := progressbar.DefaultBytes(resp.ContentLength)
	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	return absPath, err
}

func GetImageDimensions(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}
