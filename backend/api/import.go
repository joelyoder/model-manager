package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var modelIDRegex = regexp.MustCompile(`models/(\d+)`)
var versionIDRegex = regexp.MustCompile(`modelVersionId=(\d+)`)

// ImportModels ingests a multipart upload named "file" that contains JSON
// export data and recreates the referenced models/versions in the database. If
// the optional fields query parameter is provided, each imported version will
// immediately be refreshed using the same semantics as the refresh API. The
// handler writes new records and may download additional assets during refresh.
func ImportModels(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read file"})
		return
	}

	fields := c.DefaultQuery("fields", "")

	var records []ImportRecord
	if err := json.Unmarshal(data, &records); err != nil {
		var rec ImportRecord
		if err2 := json.Unmarshal(data, &rec); err2 == nil {
			records = []ImportRecord{rec}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}
	}

	successCount := 0
	failures := make([]string, 0)

	for _, r := range records {
		modelName, verName := splitName(r.Name)
		modelID := extractID(modelIDRegex, r.URL)
		origVerID := extractID(versionIDRegex, r.URL)

		if origVerID == 0 && modelID != 0 {
			apiKey := getCivitaiAPIKey()
			modelData, err := FetchCivitModel(apiKey, modelID)
			if err == nil {
				if len(modelData.ModelVersions) == 1 {
					origVerID = modelData.ModelVersions[0].ID
				} else if len(modelData.ModelVersions) > 1 && verName != "" {
					for _, mv := range modelData.ModelVersions {
						if strings.EqualFold(strings.TrimSpace(mv.Name), strings.TrimSpace(verName)) {
							origVerID = mv.ID
							break
						}
					}
				}
			}
		}

		versionID := origVerID
		if versionID == 0 {
			versionID = -int(time.Now().UnixNano())
		}

		tags := strings.Join(r.Groups, ",")
		nsfw := false
		for _, g := range r.Groups {
			if strings.EqualFold(g, "nsfw") {
				nsfw = true
				break
			}
		}

		baseModel := resolveBaseModel(r.BaseModel, r.Groups)

		var model models.Model
		err = database.DB.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Where("civit_id = ?", modelID).First(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		if err != nil {
			log.Printf("failed to load model %s: %v", r.Name, err)
			failures = append(failures, fmt.Sprintf("%s: %v", r.Name, err))
			continue
		}
		if model.ID == 0 {
			model = models.Model{
				CivitID:     modelID,
				Name:        modelName,
				Type:        r.ModelType,
				Tags:        tags,
				Nsfw:        nsfw,
				Description: r.Description,
				Weight:      1,
			}
			if err = database.DB.Create(&model).Error; err != nil {
				log.Printf("failed to create model %s: %v", r.Name, err)
				failures = append(failures, fmt.Sprintf("%s: %v", r.Name, err))
				continue
			}
		} else {
			updated := false
			if model.Weight <= 0 {
				model.Weight = 1
				updated = true
			}
			if updated {
				if err = database.DB.Save(&model).Error; err != nil {
					log.Printf("failed to update model %s weight: %v", r.Name, err)
					failures = append(failures, fmt.Sprintf("%s: %v", r.Name, err))
					continue
				}
			}
		}

		var ver models.Version
		if origVerID != 0 {
			err = database.DB.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
				Unscoped().Where("version_id = ?", origVerID).First(&ver).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = nil
			}
			if err != nil {
				log.Printf("failed to check version for %s: %v", r.Name, err)
				failures = append(failures, fmt.Sprintf("%s: %v", r.Name, err))
				continue
			}
			if ver.ID != 0 {
				// Version already exists, skip
				continue
			}
		}

		createdStr := ""
		if r.CreatedAt > 0 {
			t := time.Unix(int64(r.CreatedAt), 0)
			createdStr = t.Format(time.RFC3339)
		}

		// Determine local file path under /downloads/<type>/
		var filePath string
		if r.Location != "" {
			baseName := filepath.Base(r.Location)
			cand := filepath.Join(database.GetModelPath(), r.ModelType, baseName)
			if info, err := os.Stat(cand); err == nil && !info.IsDir() {
				if abs, aerr := filepath.Abs(cand); aerr == nil {
					filePath = abs
				} else {
					filePath = cand
				}
			} else if r.Location != "" {
				log.Printf("model file not found: %s", cand)
			}
		}

		// Look for a preview image stored under /images/<type>/
		var imagePath string
		var imgW, imgH int
		if r.Location != "" {
			base := strings.TrimSuffix(filepath.Base(r.Location), filepath.Ext(r.Location))
			imgDir := filepath.Join(database.GetImagePath(), r.ModelType)
			for _, ext := range []string{".jpg", ".jpeg", ".png"} {
				cand := filepath.Join(imgDir, base+ext)
				if info, err := os.Stat(cand); err == nil && !info.IsDir() {
					if abs, aerr := filepath.Abs(cand); aerr == nil {
						imagePath = abs
					} else {
						imagePath = cand
					}
					w, h, _ := GetImageDimensions(imagePath)
					imgW = w
					imgH = h
					break
				}
			}
			if imagePath == "" {
				log.Printf("preview image not found for %s", r.Name)
			}
		}

		ver = models.Version{
			ModelID:        model.ID,
			VersionID:      versionID,
			Name:           verName,
			BaseModel:      baseModel,
			TrainedWords:   r.PositivePrompts,
			Nsfw:           nsfw,
			Type:           r.ModelType,
			Tags:           tags,
			Description:    r.Description,
			ModelURL:       r.URL,
			SHA256:         r.SHA256Hash,
			DownloadURL:    r.DownloadURL,
			FilePath:       MakeRelativePath(filePath, database.GetModelPath()),
			ImagePath:      MakeRelativePath(imagePath, database.GetImagePath()),
			CivitCreatedAt: createdStr,
		}
		if err = database.DB.Create(&ver).Error; err != nil {
			log.Printf("failed to create version for %s: %v", r.Name, err)
			failures = append(failures, fmt.Sprintf("%s: %v", r.Name, err))
			continue
		}

		updated := false
		if model.ImagePath == "" && imagePath != "" {
			model.ImagePath = imagePath
			model.ImageWidth = imgW
			model.ImageHeight = imgH
			updated = true
		}
		if model.FilePath == "" && filePath != "" {
			model.FilePath = filePath
			updated = true
		}
		if updated {
			database.DB.Save(&model)
		}
		successCount++

		if fields != "" && origVerID != 0 {
			_ = refreshVersionData(int(ver.ID), fields)
		}
	}

	log.Printf("Import complete: %d succeeded, %d failed", successCount, len(failures))
	if len(failures) > 0 {
		log.Printf("Failed models: %s", strings.Join(failures, ", "))
	}

	c.JSON(http.StatusOK, gin.H{"message": "import complete"})
}

func extractID(re *regexp.Regexp, url string) int {
	m := re.FindStringSubmatch(url)
	if len(m) > 1 {
		id, _ := strconv.Atoi(m[1])
		return id
	}
	return 0
}

func resolveBaseModel(base string, groups []string) string {
	base = strings.TrimSpace(base)
	if base != "" {
		return base
	}
	for _, g := range groups {
		if strings.EqualFold(g, "Illustrious") {
			return "Illustrious"
		}
		if strings.EqualFold(g, "Pony") {
			return "Pony"
		}
	}
	return "SD 1.5"
}

func splitName(name string) (string, string) {
	if strings.HasSuffix(name, "]") {
		if i := strings.LastIndex(name, "["); i != -1 {
			model := strings.TrimSpace(name[:i])
			ver := strings.TrimRight(strings.TrimSpace(name[i+1:]), "]")
			return model, ver
		}
	}
	if strings.HasSuffix(name, ")") {
		if i := strings.LastIndex(name, "("); i != -1 {
			model := strings.TrimSpace(name[:i])
			ver := strings.TrimRight(strings.TrimSpace(name[i+1:]), ")")
			return model, ver
		}
	}
	return name, name
}
