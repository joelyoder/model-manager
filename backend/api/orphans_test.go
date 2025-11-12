package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"model-manager/backend/database"
	"model-manager/backend/models"
)

// setupOrphansTest prepares a temporary downloads directory and in-memory database.
func setupOrphansTest(t *testing.T) string {
	t.Helper()

	// Use a unique in-memory database for each test
	dbPath := "file:" + t.Name() + "?mode=memory&cache=shared"
	t.Setenv("MODELS_DB_PATH", dbPath)
	database.ConnectDatabase()

	// Switch to repository root so the handler's relative path is correct
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir("../.."); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(wd)
		os.RemoveAll("backend/downloads")
		os.RemoveAll("backend/target")
	})

	os.RemoveAll("backend/downloads")
	os.RemoveAll("backend/target")

	if err := os.MkdirAll("backend/downloads/sub", 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	files := []string{
		"backend/downloads/a.pt",
		"backend/downloads/b.pt",
		"backend/downloads/sub/c.pt",
	}
	for _, f := range files {
		if err := os.WriteFile(f, []byte("test"), 0o644); err != nil {
			t.Fatalf("write file %s: %v", f, err)
		}
	}
	return "backend/downloads/sub/c.pt"
}

func TestGetOrphanedFiles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	orphan := setupOrphansTest(t)

	// Insert referenced files into DB
	absA, _ := filepath.Abs("backend/downloads/a.pt")
	absB, _ := filepath.Abs("backend/downloads/b.pt")
	m := models.Model{CivitID: 1, Name: "m1", FilePath: absA, Weight: 1}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("create model: %v", err)
	}
	if err := database.DB.Create(&models.Version{ModelID: m.ID, VersionID: 1, FilePath: absB}).Error; err != nil {
		t.Fatalf("create version: %v", err)
	}

	r := gin.New()
	r.GET("/api/orphaned-files", GetOrphanedFiles)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orphaned-files", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Orphans []string `json:"orphans"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Orphans) != 1 {
		t.Fatalf("got %d orphans, want 1", len(resp.Orphans))
	}
	absOrphan, _ := filepath.Abs(orphan)
	if resp.Orphans[0] != absOrphan {
		t.Errorf("orphan = %s, want %s", resp.Orphans[0], absOrphan)
	}
}

func TestGetOrphanedFilesNone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupOrphansTest(t)

	absA, _ := filepath.Abs("backend/downloads/a.pt")
	absB, _ := filepath.Abs("backend/downloads/b.pt")
	absC, _ := filepath.Abs("backend/downloads/sub/c.pt")

	m := models.Model{CivitID: 1, Name: "m1", FilePath: absA, Weight: 1}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("create model: %v", err)
	}
	versions := []models.Version{
		{ModelID: m.ID, VersionID: 1, FilePath: absB},
		{ModelID: m.ID, VersionID: 2, FilePath: absC},
	}
	if err := database.DB.Create(&versions).Error; err != nil {
		t.Fatalf("create versions: %v", err)
	}

	r := gin.New()
	r.GET("/api/orphaned-files", GetOrphanedFiles)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orphaned-files", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Orphans []string `json:"orphans"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Orphans) != 0 {
		t.Fatalf("got %d orphans, want 0", len(resp.Orphans))
	}
}

func TestGetOrphanedFilesSymlinkDir(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupOrphansTest(t)

	absA, _ := filepath.Abs("backend/downloads/a.pt")
	absB, _ := filepath.Abs("backend/downloads/b.pt")
	absC, _ := filepath.Abs("backend/downloads/sub/c.pt")

	m := models.Model{CivitID: 1, Name: "m1", FilePath: absA, Weight: 1}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("create model: %v", err)
	}
	versions := []models.Version{
		{ModelID: m.ID, VersionID: 1, FilePath: absB},
		{ModelID: m.ID, VersionID: 2, FilePath: absC},
	}
	if err := database.DB.Create(&versions).Error; err != nil {
		t.Fatalf("create versions: %v", err)
	}

	if err := os.MkdirAll("backend/target", 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	if err := os.WriteFile("backend/target/orphan.pt", []byte("test"), 0o644); err != nil {
		t.Fatalf("write orphan: %v", err)
	}
	if err := os.Symlink("../target", "backend/downloads/link"); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	r := gin.New()
	r.GET("/api/orphaned-files", GetOrphanedFiles)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orphaned-files", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Orphans []string `json:"orphans"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Orphans) != 1 {
		t.Fatalf("got %d orphans, want 1", len(resp.Orphans))
	}
	absOrphan, _ := filepath.Abs("backend/target/orphan.pt")
	if resp.Orphans[0] != absOrphan {
		t.Errorf("orphan = %s, want %s", resp.Orphans[0], absOrphan)
	}
}

func TestGetOrphanedFilesDBSymlinkPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	orphan := setupOrphansTest(t)

	if err := os.Symlink("downloads", "backend/link"); err != nil {
		t.Fatalf("symlink: %v", err)
	}
	t.Cleanup(func() { os.Remove("backend/link") })

	absSymlinkA, _ := filepath.Abs("backend/link/a.pt")
	absB, _ := filepath.Abs("backend/downloads/b.pt")

	m := models.Model{CivitID: 1, Name: "m1", FilePath: absSymlinkA, Weight: 1}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("create model: %v", err)
	}
	if err := database.DB.Create(&models.Version{ModelID: m.ID, VersionID: 1, FilePath: absB}).Error; err != nil {
		t.Fatalf("create version: %v", err)
	}

	r := gin.New()
	r.GET("/api/orphaned-files", GetOrphanedFiles)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/orphaned-files", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Orphans []string `json:"orphans"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Orphans) != 1 {
		t.Fatalf("got %d orphans, want 1", len(resp.Orphans))
	}
	absOrphan, _ := filepath.Abs(orphan)
	if resp.Orphans[0] != absOrphan {
		t.Errorf("orphan = %s, want %s", resp.Orphans[0], absOrphan)
	}
}
