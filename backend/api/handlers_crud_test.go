package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"model-manager/backend/database"
	"model-manager/backend/models"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&models.Model{}, &models.Version{}, &models.VersionImage{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db
}

func TestCreateModel(t *testing.T) {
	setupTestDB(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/models", nil)

	CreateModel(c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		ModelID   uint `json:"modelId"`
		VersionID uint `json:"versionId"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.ModelID == 0 || resp.VersionID == 0 {
		t.Fatalf("ids not returned: %+v", resp)
	}
	var m models.Model
	if err := database.DB.First(&m, resp.ModelID).Error; err != nil {
		t.Fatalf("model not persisted: %v", err)
	}
	var v models.Version
	if err := database.DB.First(&v, resp.VersionID).Error; err != nil {
		t.Fatalf("version not persisted: %v", err)
	}
	if v.ModelID != m.ID {
		t.Fatalf("version model id = %d want %d", v.ModelID, m.ID)
	}
}

func TestUpdateModel(t *testing.T) {
	setupTestDB(t)
	m := models.Model{CivitID: 1, Name: "orig", Type: "Checkpoint"}
	database.DB.Create(&m)

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/models/abc", bytes.NewBuffer([]byte("{}")))
		UpdateModel(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/models/999", bytes.NewBuffer([]byte("{}")))
		UpdateModel(c)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(m.ID))}}
		c.Request = httptest.NewRequest(http.MethodPut, "/models/1", bytes.NewBufferString("{"))
		UpdateModel(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("updates fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(m.ID))}}
		body := models.Model{CivitID: 2, Name: "new", Type: "LORA", Tags: "t", Nsfw: true, Description: "d", ImagePath: "img", FilePath: "file", ImageWidth: 5, ImageHeight: 6}
		buf, _ := json.Marshal(body)
		c.Request = httptest.NewRequest(http.MethodPut, "/models/1", bytes.NewBuffer(buf))
		c.Request.Header.Set("Content-Type", "application/json")
		UpdateModel(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var got models.Model
		if err := database.DB.First(&got, m.ID).Error; err != nil {
			t.Fatalf("db: %v", err)
		}
		if got.CivitID != 2 || got.Name != "new" || got.Type != "LORA" || got.Tags != "t" || !got.Nsfw || got.Description != "d" || got.ImagePath != "img" || got.FilePath != "file" || got.ImageWidth != 5 || got.ImageHeight != 6 {
			t.Fatalf("model not updated: %+v", got)
		}
	})
}

func TestUpdateVersion(t *testing.T) {
	setupTestDB(t)
	m := models.Model{Name: "m", Type: "Checkpoint"}
	database.DB.Create(&m)
	v1 := models.Version{ModelID: m.ID, VersionID: 1, Name: "v1"}
	v2 := models.Version{ModelID: m.ID, VersionID: 2, Name: "v2"}
	database.DB.Create(&v1)
	database.DB.Create(&v2)

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "zzz"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/versions/zzz", bytes.NewBuffer([]byte("{}")))
		UpdateVersion(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request = httptest.NewRequest(http.MethodPut, "/versions/999", bytes.NewBuffer([]byte("{}")))
		UpdateVersion(c)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v1.ID))}}
		c.Request = httptest.NewRequest(http.MethodPut, "/versions/1", bytes.NewBufferString("{"))
		UpdateVersion(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("version id collision", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v1.ID))}}
		body := models.Version{VersionID: 2, Name: "x"}
		buf, _ := json.Marshal(body)
		c.Request = httptest.NewRequest(http.MethodPut, "/versions/1", bytes.NewBuffer(buf))
		c.Request.Header.Set("Content-Type", "application/json")
		UpdateVersion(c)
		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("updates fields", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v1.ID))}}
		body := models.Version{VersionID: 3, Name: "nv", BaseModel: "b", EarlyAccessTimeFrame: 5, SizeKB: 1.5, TrainedWords: "tw", Nsfw: true, Type: "LORA", Tags: "tg", Description: "d", Mode: "m", ModelURL: "u", CivitCreatedAt: "c", CivitUpdatedAt: "u2", SHA256: "h", DownloadURL: "du", ImagePath: "img", FilePath: "file"}
		buf, _ := json.Marshal(body)
		c.Request = httptest.NewRequest(http.MethodPut, "/versions/1", bytes.NewBuffer(buf))
		c.Request.Header.Set("Content-Type", "application/json")
		UpdateVersion(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var got models.Version
		if err := database.DB.First(&got, v1.ID).Error; err != nil {
			t.Fatalf("db: %v", err)
		}
		if got.VersionID != 3 || got.Name != "nv" || got.BaseModel != "b" || got.EarlyAccessTimeFrame != 5 || got.SizeKB != 1.5 || got.TrainedWords != "tw" || !got.Nsfw || got.Type != "LORA" || got.Tags != "tg" || got.Description != "d" || got.Mode != "m" || got.ModelURL != "u" || got.CivitCreatedAt != "c" || got.CivitUpdatedAt != "u2" || got.SHA256 != "h" || got.DownloadURL != "du" || got.ImagePath != "img" || got.FilePath != "file" {
			t.Fatalf("version not updated: %+v", got)
		}
	})
}

func TestDeleteModel(t *testing.T) {
	setupTestDB(t)
	m := models.Model{Name: "m", Type: "Checkpoint", FilePath: "m.ckpt", ImagePath: "m.png"}
	database.DB.Create(&m)
	v := models.Version{ModelID: m.ID, VersionID: 1, FilePath: "v.ckpt", ImagePath: "v.png"}
	database.DB.Create(&v)
	img := models.VersionImage{VersionID: v.ID, Path: "img.png"}
	database.DB.Create(&img)

	var trashed []string
	patch := monkey.Patch(moveToTrash, func(p string) error {
		trashed = append(trashed, p)
		return nil
	})
	defer patch.Unpatch()

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}
		DeleteModel(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		DeleteModel(c)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("deletes and trashes", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(m.ID))}}
		DeleteModel(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var count int64
		database.DB.Model(&models.Model{}).Where("id = ?", m.ID).Count(&count)
		if count != 0 {
			t.Fatalf("model not deleted")
		}
		database.DB.Model(&models.Version{}).Where("model_id = ?", m.ID).Count(&count)
		if count != 0 {
			t.Fatalf("versions not deleted")
		}
		if len(trashed) != 5 {
			t.Fatalf("moveToTrash called %d times, want 5", len(trashed))
		}
	})
}

func TestDeleteVersion(t *testing.T) {
	setupTestDB(t)
	m := models.Model{Name: "m", Type: "Checkpoint"}
	database.DB.Create(&m)
	v := models.Version{ModelID: m.ID, VersionID: 1, FilePath: "v.ckpt", ImagePath: "v.png"}
	database.DB.Create(&v)
	img := models.VersionImage{VersionID: v.ID, Path: "img.png"}
	database.DB.Create(&img)

	var trashed []string
	patch := monkey.Patch(moveToTrash, func(p string) error {
		trashed = append(trashed, p)
		return nil
	})
	defer patch.Unpatch()

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "nope"}}
		DeleteVersion(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		DeleteVersion(c)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("deletes and trashes", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v.ID))}}
		DeleteVersion(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var count int64
		database.DB.Model(&models.Version{}).Where("id = ?", v.ID).Count(&count)
		if count != 0 {
			t.Fatalf("version not deleted")
		}
		database.DB.Model(&models.VersionImage{}).Where("version_id = ?", v.ID).Count(&count)
		if count != 0 {
			t.Fatalf("images not deleted")
		}
		if len(trashed) != 3 {
			t.Fatalf("moveToTrash called %d times, want 3", len(trashed))
		}
	})
}

func TestSetVersionMainImage(t *testing.T) {
	setupTestDB(t)
	m := models.Model{Name: "m", Type: "Checkpoint", ImagePath: "old.png"}
	database.DB.Create(&m)
	v := models.Version{ModelID: m.ID, VersionID: 1, ImagePath: "old.png"}
	database.DB.Create(&v)
	img := models.VersionImage{VersionID: v.ID, Path: "new.png", Width: 10, Height: 20}
	database.DB.Create(&img)

	t.Run("invalid version id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "bad"}, {Key: "imageId", Value: "1"}}
		SetVersionMainImage(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("image not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v.ID))}, {Key: "imageId", Value: "999"}}
		SetVersionMainImage(c)
		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("updates version and model", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(v.ID))}, {Key: "imageId", Value: strconv.Itoa(int(img.ID))}}
		SetVersionMainImage(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
		var gotV models.Version
		database.DB.First(&gotV, v.ID)
		if gotV.ImagePath != "new.png" {
			t.Fatalf("version image path = %s", gotV.ImagePath)
		}
		var gotM models.Model
		database.DB.First(&gotM, m.ID)
		if gotM.ImagePath != "new.png" || gotM.ImageWidth != 10 || gotM.ImageHeight != 20 {
			t.Fatalf("model not updated: %+v", gotM)
		}
	})
}
