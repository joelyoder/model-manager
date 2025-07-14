package api

import (
	"crypto/sha256"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var CurrentDownloadProgress int64

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

	if strings.Contains(destDir, "downloads") {
		CurrentDownloadProgress = 0
	}

	total := resp.ContentLength
	buf := make([]byte, 32*1024)
	var downloaded int64
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return "", werr
			}
			downloaded += int64(n)
			if strings.Contains(destDir, "downloads") && total > 0 {
				CurrentDownloadProgress = downloaded * 100 / total
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	if strings.Contains(destDir, "downloads") {
		CurrentDownloadProgress = 100
	}

	return absPath, nil
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

func FileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
