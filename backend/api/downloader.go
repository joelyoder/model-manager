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

// DownloadFile streams the content at url into destDir/filename. The caller
// must supply a destination directory and filename; the handler ensures the
// directory exists, injects the CivitAI token when available, updates the
// package-level CurrentDownloadProgress, and returns the absolute path and
// number of bytes written. It performs filesystem writes as a side effect.
func DownloadFile(url, destDir, filename string) (string, int64, error) {
	token := getCivitaiAPIKey()
	log.Printf("Downloading %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	os.MkdirAll(destDir, os.ModePerm)
	fullPath := filepath.Join(destDir, filename)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", 0, err
	}

	out, err := os.Create(absPath)
	if err != nil {
		return "", 0, err
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
				return "", 0, werr
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
			return "", 0, err
		}
	}

	if strings.Contains(destDir, "downloads") {
		CurrentDownloadProgress = 100
	}

	return absPath, downloaded, nil
}

// GetImageDimensions opens the image at path and returns its width and height
// in pixels. The path argument must reference a readable image file; the
// function opens the file for inspection but does not modify it.
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

// FileHash computes the SHA-256 digest of the file at path. It requires read
// access to the file and returns the encoded hash without altering the source
// data.
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

// isVideoURL returns true if the provided URL points to a video file.
// It checks the file extension against a list of common video formats
// and is used to skip downloading preview videos from CivitAI.
func isVideoURL(u string) bool {
	// Strip query parameters before inspecting the extension
	if idx := strings.Index(u, "?"); idx != -1 {
		u = u[:idx]
	}
	ext := strings.ToLower(filepath.Ext(u))
	switch ext {
	case ".mp4", ".webm", ".avi", ".mov", ".mkv", ".flv", ".wmv", ".m4v", ".mpeg", ".mpg", ".gif":
		return true
	}
	return false
}
