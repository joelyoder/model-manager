package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"bou.ke/monkey"
)

func TestDownloadFileAndHelpers(t *testing.T) {
	fileContent := []byte("sample file data")

	img := image.NewRGBA(image.Rect(0, 0, 4, 2))
	var imgBuf bytes.Buffer
	if err := png.Encode(&imgBuf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	imgBytes := imgBuf.Bytes()

	var fileAuth string
	mux := http.NewServeMux()
	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		fileAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Length", strconv.Itoa(len(fileContent)))
		w.Write(fileContent)
	})
	mux.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", strconv.Itoa(len(imgBytes)))
		w.Write(imgBytes)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	patch := monkey.Patch(getCivitaiAPIKey, func() string { return "token" })
	defer patch.Unpatch()

	destDir := filepath.Join(t.TempDir(), "downloads")

	filePath, size, err := DownloadFile(srv.URL+"/file", destDir, "test.txt")
	if err != nil {
		t.Fatalf("DownloadFile file: %v", err)
	}
	if size != int64(len(fileContent)) {
		t.Errorf("size = %d, want %d", size, len(fileContent))
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if !bytes.Equal(data, fileContent) {
		t.Errorf("downloaded content mismatch: %q != %q", data, fileContent)
	}
	if CurrentDownloadProgress != 100 {
		t.Errorf("progress = %d, want 100", CurrentDownloadProgress)
	}
	if fileAuth != "Bearer token" {
		t.Errorf("authorization header = %q", fileAuth)
	}
	hash, err := FileHash(filePath)
	if err != nil {
		t.Fatalf("FileHash: %v", err)
	}
	expHash := sha256.Sum256(fileContent)
	if hash != hex.EncodeToString(expHash[:]) {
		t.Errorf("hash = %s, want %s", hash, hex.EncodeToString(expHash[:]))
	}

	imgPath, imgSize, err := DownloadFile(srv.URL+"/image", destDir, "img.png")
	if err != nil {
		t.Fatalf("DownloadFile image: %v", err)
	}
	if imgSize != int64(len(imgBytes)) {
		t.Errorf("image size = %d, want %d", imgSize, len(imgBytes))
	}
	w, h, err := GetImageDimensions(imgPath)
	if err != nil {
		t.Fatalf("GetImageDimensions: %v", err)
	}
	if w != img.Bounds().Dx() || h != img.Bounds().Dy() {
		t.Errorf("dimensions = %dx%d, want %dx%d", w, h, img.Bounds().Dx(), img.Bounds().Dy())
	}
}

func TestIsVideoURL(t *testing.T) {
	cases := []struct {
		url  string
		want bool
	}{
		{"http://example.com/video.mp4", true},
		{"http://example.com/movie.webm?download=1", true},
		{"http://example.com/clip.MOV", true},
		{"http://example.com/image.png", false},
		{"http://example.com/download?file=movie.mp4", false},
		{"http://example.com/photo.jpg?format=mp4", false},
		{"http://example.com/anim.gif?x=1", true},
	}

	for _, c := range cases {
		if got := isVideoURL(c.url); got != c.want {
			t.Errorf("isVideoURL(%q)=%v, want %v", c.url, got, c.want)
		}
	}
}
