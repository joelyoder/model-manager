package api

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func DownloadFile(url, destDir, filename string) (string, error) {
	token := os.Getenv("CIVIT_API_KEY")
	if token != "" {
		if strings.Contains(url, "?") {
			url += "&token=" + token
		} else {
			url += "?token=" + token
		}
	}

	resp, err := http.Get(url)
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

	_, err = io.Copy(out, resp.Body)
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
