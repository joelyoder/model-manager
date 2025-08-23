package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	t.Setenv("MODELS_DB_PATH", "file::memory:?cache=shared")
	database.ConnectDatabase()

	m1 := models.Model{CivitID: 1, Name: "Alpha"}
	if err := database.DB.Create(&m1).Error; err != nil {
		t.Fatalf("create model1: %v", err)
	}
	if err := database.DB.Create(&models.Version{ModelID: m1.ID, VersionID: 11, BaseModel: "SD1", Type: "checkpoint", Nsfw: false, Tags: "tag1"}).Error; err != nil {
		t.Fatalf("create version1: %v", err)
	}

	m2 := models.Model{CivitID: 2, Name: "Beta"}
	if err := database.DB.Create(&m2).Error; err != nil {
		t.Fatalf("create model2: %v", err)
	}
	if err := database.DB.Create(&models.Version{ModelID: m2.ID, VersionID: 22, BaseModel: "SD1", Type: "lora", Nsfw: true, Tags: "tag2,tag3"}).Error; err != nil {
		t.Fatalf("create version2: %v", err)
	}

	m3 := models.Model{CivitID: 3, Name: "Gamma"}
	if err := database.DB.Create(&m3).Error; err != nil {
		t.Fatalf("create model3: %v", err)
	}
	if err := database.DB.Create(&models.Version{ModelID: m3.ID, VersionID: 33, BaseModel: "SD2", Type: "checkpoint", Nsfw: false, Tags: "tag1,tag2"}).Error; err != nil {
		t.Fatalf("create version3: %v", err)
	}

	m4 := models.Model{CivitID: 4, Name: "Delta"}
	if err := database.DB.Create(&m4).Error; err != nil {
		t.Fatalf("create model4: %v", err)
	}
	if err := database.DB.Create(&models.Version{ModelID: m4.ID, VersionID: 44, BaseModel: "SD2", Type: "lora", Nsfw: false, Tags: "tag3"}).Error; err != nil {
		t.Fatalf("create version4: %v", err)
	}
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/models", GetModels)
	r.GET("/models/count", GetModelsCount)
	return r
}

func TestGetModelsAndCountFilters(t *testing.T) {
	setupTestDB(t)
	r := newTestRouter()

	tests := []struct {
		name      string
		query     string
		wantNames []string
	}{
		{"search", "?search=alpha", []string{"Alpha"}},
		{"baseModel", "?baseModel=SD1", []string{"Beta", "Alpha"}},
		{"modelType", "?modelType=lora", []string{"Delta", "Beta"}},
		{"hideNsfw", "?hideNsfw=1", []string{"Delta", "Gamma", "Alpha"}},
		{"tags", "?tags=tag3", []string{"Delta", "Beta"}},
		{"multiTags", "?tags=tag1,tag2", []string{"Gamma"}},
		{"combo", "?baseModel=SD2&modelType=checkpoint", []string{"Gamma"}},
		{"comboNone", "?baseModel=SD1&modelType=lora&hideNsfw=1", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/models"+tt.query, nil)
			r.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("status = %d", w.Code)
			}
			var modelsResp []models.Model
			if err := json.Unmarshal(w.Body.Bytes(), &modelsResp); err != nil {
				t.Fatalf("unmarshal models: %v", err)
			}
			if len(modelsResp) != len(tt.wantNames) {
				t.Fatalf("got %d models, want %d", len(modelsResp), len(tt.wantNames))
			}
			for i, m := range modelsResp {
				if m.Name != tt.wantNames[i] {
					t.Errorf("model %d name = %s, want %s", i, m.Name, tt.wantNames[i])
				}
			}

			w = httptest.NewRecorder()
			req = httptest.NewRequest(http.MethodGet, "/models/count"+tt.query, nil)
			r.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Fatalf("count status = %d", w.Code)
			}
			var countResp struct {
				Count int64 `json:"count"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &countResp); err != nil {
				t.Fatalf("unmarshal count: %v", err)
			}
			if int(countResp.Count) != len(tt.wantNames) {
				t.Errorf("count = %d, want %d", countResp.Count, len(tt.wantNames))
			}
		})
	}

	t.Run("pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/models?limit=2&page=2", nil)
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var modelsResp []models.Model
		if err := json.Unmarshal(w.Body.Bytes(), &modelsResp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		want := []string{"Beta", "Alpha"}
		if len(modelsResp) != len(want) {
			t.Fatalf("got %d models, want %d", len(modelsResp), len(want))
		}
		for i, m := range modelsResp {
			if m.Name != want[i] {
				t.Errorf("model %d name = %s, want %s", i, m.Name, want[i])
			}
		}
	})
}
