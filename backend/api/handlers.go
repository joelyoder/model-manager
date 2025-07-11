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
	var modelsList []models.Model
	database.DB.Preload("Versions").Find(&modelsList)
	c.JSON(200, modelsList)
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
	database.DB.Where("version_id = ?", id).First(&existingVersion)
	if existingVersion.ID > 0 {
		c.JSON(200, gin.H{"message": "Version already exists"})
		return
	}

	var model models.Model
	database.DB.Where("civit_id = ?", verData.ModelID).First(&model)
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
		ImagePath:            imagePath,
		FilePath:             filePath,
	})

	c.JSON(200, gin.H{"message": "Version synced", "versionId": verData.ID})
}

func processModels(items []CivitModel, apiKey string) {
	for _, item := range items {
		var existing models.Model
		database.DB.Where("civit_id = ?", item.ID).First(&existing)
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
			database.DB.Where("version_id = ?", verData.ID).First(&versionExists)
			if versionExists.ID > 0 {
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
