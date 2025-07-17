package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

// ImportModels handles JSON file uploads and inserts records into the database.
// If the optional "fields" query parameter is provided, each imported version
// will immediately be refreshed using the same semantics as the refresh API.
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

	for _, r := range records {
		modelID := extractID(modelIDRegex, r.URL)
		versionID := extractID(versionIDRegex, r.URL)

		modelName, verName := splitName(r.Name)
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
			}
			database.DB.Create(&model)
		}

		var ver models.Version
		err = database.DB.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Unscoped().Where("version_id = ?", versionID).First(&ver).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		if err != nil {
			continue
		}
		if ver.ID != 0 {
			continue
		}

		createdStr := ""
		if r.CreatedAt > 0 {
			t := time.Unix(int64(r.CreatedAt), 0)
			createdStr = t.Format(time.RFC3339)
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
			FilePath:       r.Location,
			CivitCreatedAt: createdStr,
		}
		database.DB.Create(&ver)
		if fields != "" && versionID != 0 {
			_ = refreshVersionData(int(ver.ID), fields)
		}
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
