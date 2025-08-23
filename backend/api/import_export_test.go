package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"

	"model-manager/backend/database"
	"model-manager/backend/models"
)

func initTestDB(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dir := t.TempDir()
	t.Setenv("MODELS_DB_PATH", filepath.Join(dir, "test.db"))
	database.ConnectDatabase()
}

func uploadRequest(t *testing.T, url string, data []byte) *http.Request {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	if data != nil {
		part, err := w.CreateFormFile("file", "data.json")
		if err != nil {
			t.Fatalf("CreateFormFile: %v", err)
		}
		if _, err := part.Write(data); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestImportModels(t *testing.T) {
	t.Run("single record", func(t *testing.T) {
		initTestDB(t)
		rec := ImportRecord{Name: "Foo [v1]", ModelType: "Checkpoint", Groups: []string{"tag1"}}
		buf, _ := json.Marshal(rec)
		req := uploadRequest(t, "/import", buf)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportModels(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var m models.Model
		if err := database.DB.Preload("Versions").First(&m).Error; err != nil {
			t.Fatalf("db: %v", err)
		}
		if m.Name != "Foo" || m.Type != "Checkpoint" || m.Tags != "tag1" {
			t.Fatalf("model fields: %+v", m)
		}
		if len(m.Versions) != 1 || m.Versions[0].Name != "v1" {
			t.Fatalf("version: %+v", m.Versions)
		}
	})

	t.Run("multi record", func(t *testing.T) {
		initTestDB(t)
		recs := []ImportRecord{
			{Name: "Bar [v1]", ModelType: "Checkpoint", URL: "http://example.com/models/1?modelVersionId=11"},
			{Name: "Baz [v2]", ModelType: "LORA", URL: "http://example.com/models/2?modelVersionId=22"},
		}
		buf, _ := json.Marshal(recs)
		req := uploadRequest(t, "/import", buf)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportModels(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var count int64
		database.DB.Model(&models.Model{}).Count(&count)
		if count != 2 {
			t.Fatalf("model count = %d", count)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		initTestDB(t)
		req := httptest.NewRequest(http.MethodPost, "/import", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportModels(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		initTestDB(t)
		req := uploadRequest(t, "/import", []byte("{"))
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportModels(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})
}

func TestImportDatabase(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		initTestDB(t)
		modelsList := []models.Model{
			{
				CivitID: 1,
				Name:    "M1",
				Type:    "Checkpoint",
				Versions: []models.Version{{
					VersionID: 101,
					Name:      "v1",
					Images:    []models.VersionImage{{Path: "img1"}},
				}},
			},
			{
				CivitID: 2,
				Name:    "M2",
				Type:    "LORA",
			},
		}
		buf, _ := json.Marshal(modelsList)
		req := uploadRequest(t, "/db/import", buf)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportDatabase(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var count int64
		database.DB.Model(&models.Model{}).Count(&count)
		if count != 2 {
			t.Fatalf("model count = %d", count)
		}
		var m models.Model
		if err := database.DB.Preload("Versions.Images").First(&m, "name = ?", "M1").Error; err != nil {
			t.Fatalf("db: %v", err)
		}
		if len(m.Versions) != 1 || len(m.Versions[0].Images) != 1 {
			t.Fatalf("nested not imported: %+v", m)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		initTestDB(t)
		req := httptest.NewRequest(http.MethodPost, "/db/import", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportDatabase(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		initTestDB(t)
		req := uploadRequest(t, "/db/import", []byte("{"))
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		ImportDatabase(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})
}

func TestExportModels(t *testing.T) {
	initTestDB(t)
	m := models.Model{CivitID: 1, Name: "M", Type: "Checkpoint", Tags: "t", Nsfw: true, Description: "d"}
	database.DB.Create(&m)
	v := models.Version{ModelID: m.ID, VersionID: 1, Name: "v"}
	database.DB.Create(&v)
	img := models.VersionImage{VersionID: v.ID, Path: "p", Width: 1, Height: 2, Hash: "h"}
	database.DB.Create(&img)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/export", nil)
	ExportModels(c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var got []models.Model
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got) != 1 || len(got[0].Versions) != 1 || len(got[0].Versions[0].Images) != 1 {
		t.Fatalf("export mismatch: %v", got)
	}
	if got[0].Name != "M" || got[0].Versions[0].Name != "v" || got[0].Versions[0].Images[0].Path != "p" {
		t.Fatalf("data mismatch: %+v", got[0])
	}
}

func TestRefreshVersionData(t *testing.T) {
	newModel := CivitModel{
		ID:          1,
		Name:        "NewModel",
		Type:        "Checkpoint",
		Description: "NewDesc",
		Nsfw:        true,
		Tags:        []string{"tag1"},
		Mode:        "mode",
	}
	newVer := VersionResponse{
		ID:        10,
		ModelID:   1,
		Name:      "NewVer",
		BaseModel: "NewBase",
		Created:   "2020-01-01",
		Updated:   "2020-01-02",
		ModelFiles: []ModelFile{{
			SizeKB:      100,
			DownloadURL: "u",
			Hashes: struct {
				SHA256 string `json:"SHA256"`
			}{SHA256: "hash"},
		}},
		TrainedWords: []string{"word"},
	}

	patchModel := monkey.Patch(FetchCivitModel, func(string, int) (CivitModel, error) { return newModel, nil })
	defer patchModel.Unpatch()
	patchVer := monkey.Patch(FetchModelVersion, func(string, int) (VersionResponse, error) { return newVer, nil })
	defer patchVer.Unpatch()
	patchKey := monkey.Patch(getCivitaiAPIKey, func() string { return "" })
	defer patchKey.Unpatch()

	t.Run("metadata only", func(t *testing.T) {
		initTestDB(t)
		m := models.Model{CivitID: 1, Name: "Old", Description: "OldDesc"}
		database.DB.Create(&m)
		v := models.Version{ModelID: m.ID, VersionID: 10, Name: "OldVer", Description: "OldDesc"}
		database.DB.Create(&v)

		if err := refreshVersionData(int(v.ID), "metadata"); err != nil {
			t.Fatalf("refresh: %v", err)
		}
		var gotM models.Model
		var gotV models.Version
		database.DB.First(&gotM, m.ID)
		database.DB.First(&gotV, v.ID)
		if gotM.Name != "NewModel" || gotV.Name != "NewVer" || gotV.BaseModel != "NewBase" {
			t.Fatalf("metadata not updated: %+v %+v", gotM, gotV)
		}
		if gotM.Description != "OldDesc" || gotV.Description != "OldDesc" {
			t.Fatalf("description changed: %+v %+v", gotM, gotV)
		}
	})

	t.Run("description only", func(t *testing.T) {
		initTestDB(t)
		m := models.Model{CivitID: 1, Name: "Old", Description: "OldDesc"}
		database.DB.Create(&m)
		v := models.Version{ModelID: m.ID, VersionID: 10, Name: "OldVer", Description: "OldDesc"}
		database.DB.Create(&v)

		if err := refreshVersionData(int(v.ID), "description"); err != nil {
			t.Fatalf("refresh: %v", err)
		}
		var gotM models.Model
		var gotV models.Version
		database.DB.First(&gotM, m.ID)
		database.DB.First(&gotV, v.ID)
		if gotM.Description != "NewDesc" || gotV.Description != "NewDesc" {
			t.Fatalf("description not updated: %+v %+v", gotM, gotV)
		}
		if gotM.Name != "Old" || gotV.Name != "OldVer" {
			t.Fatalf("metadata changed: %+v %+v", gotM, gotV)
		}
	})

	t.Run("all fields", func(t *testing.T) {
		initTestDB(t)
		m := models.Model{CivitID: 1, Name: "Old", Description: "OldDesc"}
		database.DB.Create(&m)
		v := models.Version{ModelID: m.ID, VersionID: 10, Name: "OldVer", Description: "OldDesc"}
		database.DB.Create(&v)

		if err := refreshVersionData(int(v.ID), "all"); err != nil {
			t.Fatalf("refresh: %v", err)
		}
		var gotM models.Model
		var gotV models.Version
		database.DB.First(&gotM, m.ID)
		database.DB.First(&gotV, v.ID)
		if gotM.Name != "NewModel" || gotM.Description != "NewDesc" || gotV.Name != "NewVer" || gotV.Description != "NewDesc" {
			t.Fatalf("fields not updated: %+v %+v", gotM, gotV)
		}
	})
}
