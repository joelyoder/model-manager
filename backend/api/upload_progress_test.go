package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"

	"model-manager/backend/database"
	"model-manager/backend/models"
)

// setupUploadTest initializes test environment and database.
func setupUploadTest(t *testing.T) string {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	t.Setenv("MODELS_DB_PATH", filepath.Join(dir, "test.db"))
	database.ConnectDatabase()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(cwd) })
	return dir
}

func TestUploadVersionFile(t *testing.T) {
	dir := setupUploadTest(t)

	m := models.Model{Name: "test", Type: "Checkpoint", Weight: 1}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("create model: %v", err)
	}
	v := models.Version{ModelID: m.ID, VersionID: 1, Name: "v1", Type: "Checkpoint"}
	if err := database.DB.Create(&v).Error; err != nil {
		t.Fatalf("create version: %v", err)
	}

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("file", "file.txt")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write([]byte("content")); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("writer close: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/versions/"+strconv.Itoa(int(v.ID))+"/upload", body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v.ID))}}
	c.Request = req

	UploadVersionFile(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var resp struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	expected := filepath.Join(dir, "backend", "downloads", "Checkpoint", "file.txt")
	if resp.Path != expected {
		t.Fatalf("path = %s, want %s", resp.Path, expected)
	}
	if _, err := os.Stat(resp.Path); err != nil {
		t.Fatalf("uploaded file missing: %v", err)
	}
	var vdb models.Version
	if err := database.DB.First(&vdb, v.ID).Error; err != nil {
		t.Fatalf("version from db: %v", err)
	}
	if vdb.FilePath != expected {
		t.Errorf("version filepath = %s, want %s", vdb.FilePath, expected)
	}
	var mdb models.Model
	if err := database.DB.First(&mdb, m.ID).Error; err != nil {
		t.Fatalf("model from db: %v", err)
	}
	if mdb.FilePath != expected {
		t.Errorf("model filepath = %s, want %s", mdb.FilePath, expected)
	}
}

// requestProgress issues a GET request to the progress endpoint and returns the progress value.
func requestProgress() (int64, error) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/progress", nil)

	GetDownloadProgress(c)
	if rec.Code != http.StatusOK {
		return 0, fmt.Errorf("status = %d", rec.Code)
	}
	var resp struct {
		Progress int64 `json:"progress"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		return 0, err
	}
	return resp.Progress, nil
}

// getProgress is a helper for tests that must fail the test on error.
func getProgress(t *testing.T) int64 {
	t.Helper()
	p, err := requestProgress()
	if err != nil {
		t.Fatalf("request progress: %v", err)
	}
	return p
}

func TestGetDownloadProgress(t *testing.T) {
	gin.SetMode(gin.TestMode)
	CurrentDownloadProgress = 42
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/progress", nil)

	GetDownloadProgress(c)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var resp struct {
		Progress int64 `json:"progress"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Progress != 42 {
		t.Errorf("progress = %d, want 42", resp.Progress)
	}
}

// TestGetDownloadProgressUpdates verifies the endpoint reflects sequential updates
// to CurrentDownloadProgress and that values persist or reset appropriately.
func TestGetDownloadProgressUpdates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	CurrentDownloadProgress = 0
	if p := getProgress(t); p != 0 {
		t.Errorf("progress = %d, want 0", p)
	}

	CurrentDownloadProgress = 50
	if p := getProgress(t); p != 50 {
		t.Errorf("progress = %d, want 50", p)
	}
	// Ensure value persists without another update.
	if p := getProgress(t); p != 50 {
		t.Errorf("progress persisted = %d, want 50", p)
	}

	CurrentDownloadProgress = 100
	if p := getProgress(t); p != 100 {
		t.Errorf("progress = %d, want 100", p)
	}

	// Reset back to 0 and confirm.
	CurrentDownloadProgress = 0
	if p := getProgress(t); p != 0 {
		t.Errorf("progress after reset = %d, want 0", p)
	}
}

// TestGetDownloadProgressOutOfRange verifies negative or >100 values are returned as-is.
func TestGetDownloadProgressOutOfRange(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, v := range []int64{-10, 150} {
		CurrentDownloadProgress = v
		if p := getProgress(t); p != v {
			t.Errorf("progress = %d, want %d", p, v)
		}
	}
}

// TestGetDownloadProgressConcurrent ensures simultaneous requests return the current value.
func TestGetDownloadProgressConcurrent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	CurrentDownloadProgress = 66

	const n = 10
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p, err := requestProgress()
			if err != nil {
				t.Errorf("request: %v", err)
				return
			}
			if p != 66 {
				t.Errorf("progress = %d, want 66", p)
			}
		}()
	}
	wg.Wait()
}
