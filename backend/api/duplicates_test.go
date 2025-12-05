package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"model-manager/backend/database"
	"model-manager/backend/models"
)

func setupDuplicatesTest(t *testing.T) {
	t.Helper()
	dbPath := "file:" + t.Name() + "?mode=memory&cache=shared"
	t.Setenv("MODELS_DB_PATH", dbPath)
	database.ConnectDatabase()
}

func TestGetDuplicateFilePaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupDuplicatesTest(t)

	m1 := models.Model{CivitID: 1, Name: "m1", Weight: 1}
	m2 := models.Model{CivitID: 2, Name: "m2", Weight: 1}
	if err := database.DB.Create(&[]models.Model{m1, m2}).Error; err != nil {
		t.Fatalf("create models: %v", err)
	}
	path := "/tmp/dup.pt"
	versions := []models.Version{
		{ModelID: 1, VersionID: 1, Name: "v1", FilePath: path},
		{ModelID: 2, VersionID: 2, Name: "v2", FilePath: path},
		{ModelID: 2, VersionID: 3, Name: "v3", FilePath: "/tmp/unique.pt"},
	}
	if err := database.DB.Create(&versions).Error; err != nil {
		t.Fatalf("create versions: %v", err)
	}

	r := gin.New()
	r.GET("/api/duplicate-file-paths", GetDuplicateFilePaths)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/duplicate-file-paths", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Duplicates []struct {
			Path     string `json:"path"`
			Versions []struct {
				ModelName   string `json:"modelName"`
				VersionName string `json:"versionName"`
			} `json:"versions"`
		} `json:"duplicates"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Duplicates) != 1 {
		t.Fatalf("got %d duplicates, want 1", len(resp.Duplicates))
	}
	if resp.Duplicates[0].Path != path {
		t.Errorf("path = %s, want %s", resp.Duplicates[0].Path, path)
	}
	if len(resp.Duplicates[0].Versions) != 2 {
		t.Fatalf("got %d versions, want 2", len(resp.Duplicates[0].Versions))
	}
}

func TestGetDuplicateFilePathsNone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupDuplicatesTest(t)
	m := models.Model{CivitID: 1, Name: "m1", Weight: 1}
	if err := database.DB.Create(&m).Error; err != nil {
		t.Fatalf("create model: %v", err)
	}
	versions := []models.Version{
		{ModelID: m.ID, VersionID: 1, Name: "v1", FilePath: "/tmp/a.pt"},
		{ModelID: m.ID, VersionID: 2, Name: "v2", FilePath: "/tmp/b.pt"},
	}
	if err := database.DB.Create(&versions).Error; err != nil {
		t.Fatalf("create versions: %v", err)
	}

	r := gin.New()
	r.GET("/api/duplicate-file-paths", GetDuplicateFilePaths)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/duplicate-file-paths", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Duplicates []interface{} `json:"duplicates"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Duplicates) != 0 {
		t.Fatalf("got %d duplicates, want 0", len(resp.Duplicates))
	}
}
