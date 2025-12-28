package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"model-manager/backend/database"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var CurrentDownloadProgress int64

var (
	downloadMu            sync.Mutex
	currentDownloadCancel context.CancelFunc
	currentDownloadPath   string
	currentDownloadToken  *struct{}
)

// DownloadFile streams the content at url into destDir/filename. The caller
// must supply a destination directory and filename; the handler ensures the
// directory exists, injects the CivitAI token when available, updates the
// package-level CurrentDownloadProgress, and returns the absolute path and
// number of bytes written. It performs filesystem writes as a side effect.
func DownloadFile(url, destDir, filename string) (string, int64, error) {
	apiToken := getCivitaiAPIKey()
	log.Printf("Downloading %s", url)

	fullPath := filepath.Join(destDir, filename)
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", 0, err
	}

	// Determine if this is a model download by checking if it matches the configured model path
	modelPath, _ := filepath.Abs(database.GetModelPath())
	destAbs, _ := filepath.Abs(destDir)
	// Simple check: is destAbs inside modelPath?
	// or just check if it contains the model path root
	isModelDownload := strings.HasPrefix(strings.ToLower(destAbs), strings.ToLower(modelPath))

	// Fallback to "downloads" check if path resolution fails or is weird
	if !isModelDownload {
		isModelDownload = strings.Contains(strings.ToLower(destDir), "downloads")
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if isModelDownload {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", 0, err
	}
	var downloadToken *struct{}
	if isModelDownload {
		downloadToken = &struct{}{}
		downloadMu.Lock()
		currentDownloadCancel = cancel
		currentDownloadPath = absPath
		currentDownloadToken = downloadToken
		CurrentDownloadProgress = 0
		downloadMu.Unlock()
	}
	if apiToken != "" {
		req.Header.Add("Authorization", "Bearer "+apiToken)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	os.MkdirAll(destDir, os.ModePerm)

	out, err := os.Create(absPath)
	if err != nil {
		return "", 0, err
	}
	defer out.Close()

	if isModelDownload {
		defer func() {
			downloadMu.Lock()
			if currentDownloadToken == downloadToken {
				currentDownloadCancel = nil
				currentDownloadPath = ""
				currentDownloadToken = nil
			}
			downloadMu.Unlock()
			if cancel != nil {
				cancel()
			}
		}()
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
			if isModelDownload && total > 0 {
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

	if isModelDownload {
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

// CancelActiveDownload cancels the current model download if one is in progress.
// It stops the HTTP request, resets the progress indicator, and removes any
// partially written file from disk.
func CancelActiveDownload() (bool, error) {
	downloadMu.Lock()
	cancel := currentDownloadCancel
	path := currentDownloadPath
	currentDownloadCancel = nil
	currentDownloadPath = ""
	currentDownloadToken = nil
	downloadMu.Unlock()

	if cancel == nil {
		return false, nil
	}

	cancel()
	CurrentDownloadProgress = 0
	if path != "" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return true, err
		}
	}
	return true, nil
}
