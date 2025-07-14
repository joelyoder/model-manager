package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

func GetModels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	search := c.Query("search")

	var modelsList []models.Model
	q := database.DB.Preload("Versions")
	if search != "" {
		q = q.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(search)+"%")
	}
	q.Limit(limit).Offset((page - 1) * limit).Find(&modelsList)
	c.JSON(http.StatusOK, modelsList)
}

// GetModel returns a single model by ID with its versions
func GetModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	var model models.Model
	if err := database.DB.Preload("Versions").First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, model)
}

func SyncCivitModels(c *gin.Context) {
	apiKey := os.Getenv("CIVIT_API_KEY")
	items, err := FetchCivitModels(apiKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch models"})
		return
	}
	processModels(items, apiKey)
	c.JSON(200, gin.H{"message": "Models synced successfully."})
}

func SyncCivitModelByID(c *gin.Context) {
	log.Println("Hit /api/sync/:id")
	apiKey := os.Getenv("CIVIT_API_KEY")
	id := c.Param("id")

	modelID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid model ID"})
		return
	}

	url := fmt.Sprintf("https://civitai.com/api/v1/models/%d", modelID)
	log.Printf("GET %s", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		c.JSON(500, gin.H{"error": "Failed to fetch model"})
		return
	}
	defer resp.Body.Close()

	var model CivitModel
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &model)

	processModels([]CivitModel{model}, apiKey)
	c.JSON(200, gin.H{"message": "Model synced successfully", "modelId": modelID})
}

func GetModelVersions(c *gin.Context) {
	apiKey := os.Getenv("CIVIT_API_KEY")
	modelID := c.Param("id")

	id, err := strconv.Atoi(modelID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid model ID"})
		return
	}

	model, err := FetchCivitModel(apiKey, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch model"})
		return
	}

	var versions []VersionInfo
	for _, vs := range model.ModelVersions {
		ver, err := FetchModelVersion(apiKey, vs.ID)
		if err != nil {
			continue
		}

		sizeKB := 0.0
		if len(ver.ModelFiles) > 0 {
			sizeKB = ver.ModelFiles[0].SizeKB
		}

		versions = append(versions, VersionInfo{
			ID:           ver.ID,
			Name:         ver.Name,
			BaseModel:    ver.BaseModel,
			SizeKB:       sizeKB,
			TrainedWords: ver.TrainedWords,
		})
	}

	c.JSON(200, versions)
}

func SyncVersionByID(c *gin.Context) {
	apiKey := os.Getenv("CIVIT_API_KEY")
	versionID := c.Param("versionId")

	id, err := strconv.Atoi(versionID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid version ID"})
		return
	}

	verData, err := FetchModelVersion(apiKey, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch version"})
		return
	}

	var existingVersion models.Version
	database.DB.Unscoped().Where("version_id = ?", id).Find(&existingVersion)
	if existingVersion.ID > 0 {
		log.Printf("Skipping download: version %d already exists", id)
		c.JSON(http.StatusConflict, gin.H{"error": "Version already exists"})
		return
	}

	var model models.Model
	database.DB.Unscoped().Where("civit_id = ?", verData.ModelID).Find(&model)
	if model.ID == 0 {
		modelData, _ := FetchCivitModel(apiKey, verData.ModelID)
		model = models.Model{
			CivitID:     modelData.ID,
			Name:        modelData.Name,
			Nsfw:        modelData.Nsfw,
			Type:        modelData.Type,
			Tags:        strings.Join(modelData.Tags, ","),
			Description: modelData.Description,
			CreatedAt:   modelData.Created,
			UpdatedAt:   modelData.Updated,
		}
		database.DB.Create(&model)
	}

	var filePath, imagePath string
	var imgW, imgH int
	if len(verData.ModelFiles) > 0 {
		filePath, _ = DownloadFile(verData.ModelFiles[0].DownloadURL, "./backend/downloads/"+model.Type, verData.ModelFiles[0].Name)
	}

	if len(verData.Images) > 0 {
		imageURL := verData.Images[0].URL
		if imageURL == "" {
			imageURL = verData.Images[0].URLSmall
		}
		if imageURL != "" {
			imagePath, _ = DownloadFile(imageURL, "./backend/images/"+model.Type, fmt.Sprintf("%d.jpg", verData.ID))
			imgW, imgH, _ = GetImageDimensions(imagePath)
		}
	}

	if model.ImagePath == "" && imagePath != "" {
		model.ImagePath = imagePath
		model.ImageWidth = imgW
		model.ImageHeight = imgH
	}
	if model.FilePath == "" && filePath != "" {
		model.FilePath = filePath
	}
	database.DB.Save(&model)

	database.DB.Create(&models.Version{
		ModelID:              model.ID,
		VersionID:            verData.ID,
		Name:                 verData.Name,
		BaseModel:            verData.BaseModel,
		CreatedAt:            verData.Created,
		EarlyAccessTimeFrame: verData.EarlyAccessTimeFrame,
		SizeKB:               verData.ModelFiles[0].SizeKB,
		TrainedWords:         strings.Join(verData.TrainedWords, ","),
		ModelURL:             fmt.Sprintf("https://civitai.com/models/%d?modelVersionId=%d", verData.ModelID, verData.ID),
		ImagePath:            imagePath,
		FilePath:             filePath,
	})

	c.JSON(200, gin.H{"message": "Version synced", "versionId": verData.ID})
}

func processModels(items []CivitModel, apiKey string) {
	for _, item := range items {
		var existing models.Model
		database.DB.Where("civit_id = ?", item.ID).Find(&existing)
		if existing.ID == 0 {
			existing = models.Model{
				CivitID:     item.ID,
				Name:        item.Name,
				Nsfw:        item.Nsfw,
				Type:        item.Type,
				Tags:        strings.Join(item.Tags, ","),
				Description: item.Description,
				CreatedAt:   item.Created,
				UpdatedAt:   item.Updated,
			}
			database.DB.Create(&existing)
		}

		for _, version := range item.ModelVersions {
			verData, err := FetchModelVersion(apiKey, version.ID)
			if err != nil {
				continue
			}

			var versionExists models.Version
			database.DB.Unscoped().Where("version_id = ?", verData.ID).Find(&versionExists)
			if versionExists.ID > 0 {
				log.Printf("Skipping download: version %d already exists", verData.ID)
				continue
			}

			var filePath, imagePath string
			var imgW, imgH int
			if len(verData.ModelFiles) > 0 {
				fileURL := verData.ModelFiles[0].DownloadURL
				fileName := verData.ModelFiles[0].Name
				filePath, _ = DownloadFile(fileURL, "./backend/downloads/"+item.Type, fileName)
			}
			if len(verData.Images) > 0 {
				imageURL := verData.Images[0].URL
				if imageURL == "" {
					imageURL = verData.Images[0].URLSmall
				}
				if imageURL != "" {
					imagePath, _ = DownloadFile(imageURL, "./backend/images/"+item.Type, fmt.Sprintf("%d.jpg", verData.ID))
					imgW, imgH, _ = GetImageDimensions(imagePath)
				}
			}

			database.DB.Create(&models.Version{
				ModelID:              existing.ID,
				VersionID:            verData.ID,
				Name:                 verData.Name,
				BaseModel:            verData.BaseModel,
				CreatedAt:            verData.Created,
				EarlyAccessTimeFrame: verData.EarlyAccessTimeFrame,
				SizeKB:               verData.ModelFiles[0].SizeKB,
				TrainedWords:         strings.Join(verData.TrainedWords, ","),
				ModelURL:             fmt.Sprintf("https://civitai.com/models/%d?modelVersionId=%d", item.ID, verData.ID),
				ImagePath:            imagePath,
				FilePath:             filePath,
			})

			if existing.ImagePath == "" && imagePath != "" {
				existing.ImagePath = imagePath
				existing.ImageWidth = imgW
				existing.ImageHeight = imgH
			}
			if existing.FilePath == "" && filePath != "" {
				existing.FilePath = filePath
			}
			database.DB.Save(&existing)
		}
	}
}

// DeleteModel removes a model and its versions from the database. It also deletes any associated files and images stored on disk.
func DeleteModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	var model models.Model
	if err := database.DB.Preload("Versions").First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	if model.FilePath != "" {
		os.Remove(model.FilePath)
	}
	if model.ImagePath != "" {
		os.Remove(model.ImagePath)
	}
	for _, v := range model.Versions {
		if v.FilePath != "" {
			os.Remove(v.FilePath)
		}
		if v.ImagePath != "" {
			os.Remove(v.ImagePath)
		}
	}

	database.DB.Unscoped().Where("model_id = ?", model.ID).Delete(&models.Version{})
	database.DB.Unscoped().Delete(&model)

	c.JSON(http.StatusOK, gin.H{"message": "Model deleted"})
}

// GetVersion returns a single model version along with its parent model
func GetVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	var version models.Version
	if err := database.DB.First(&version, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	var model models.Model
	if err := database.DB.First(&model, version.ModelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": model, "version": version})
}

// DeleteVersion removes a single model version and associated files from the database.
func DeleteVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	var version models.Version
	if err := database.DB.First(&version, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	if version.FilePath != "" {
		os.Remove(version.FilePath)
	}
	if version.ImagePath != "" {
		os.Remove(version.ImagePath)
	}

	database.DB.Unscoped().Delete(&models.Version{}, version.ID)

	var remaining int64
	database.DB.Model(&models.Version{}).Where("model_id = ?", version.ModelID).Count(&remaining)
	if remaining == 0 {
		var model models.Model
		if err := database.DB.First(&model, version.ModelID).Error; err == nil {
			if model.FilePath != "" && model.FilePath == version.FilePath {
				os.Remove(model.FilePath)
				model.FilePath = ""
			}
			if model.ImagePath != "" && model.ImagePath == version.ImagePath {
				os.Remove(model.ImagePath)
				model.ImagePath = ""
			}
			database.DB.Save(&model)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Version deleted"})
}

// UpdateModel updates an existing model with new values
func UpdateModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	var model models.Model
	if err := database.DB.First(&model, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	var input models.Model
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	model.Name = input.Name
	model.Type = input.Type
	model.Tags = input.Tags
	model.Nsfw = input.Nsfw
	model.Description = input.Description
	model.CreatedAt = input.CreatedAt
	model.UpdatedAt = input.UpdatedAt
	model.ImagePath = input.ImagePath
	model.FilePath = input.FilePath
	model.ImageWidth = input.ImageWidth
	model.ImageHeight = input.ImageHeight

	if err := database.DB.Save(&model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model"})
		return
	}

	c.JSON(http.StatusOK, model)
}

// UpdateVersion updates an existing version with new values
func UpdateVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	var version models.Version
	if err := database.DB.First(&version, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	var input models.Version
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	version.Name = input.Name
	version.BaseModel = input.BaseModel
	version.CreatedAt = input.CreatedAt
	version.EarlyAccessTimeFrame = input.EarlyAccessTimeFrame
	version.SizeKB = input.SizeKB
	version.TrainedWords = input.TrainedWords
	version.ModelURL = input.ModelURL
	version.ImagePath = input.ImagePath
	version.FilePath = input.FilePath

	if err := database.DB.Save(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update version"})
		return
	}

	c.JSON(http.StatusOK, version)
}
