package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	apiKey := database.GetAPIKey()
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
	apiKey := database.GetAPIKey()
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
	apiKey := database.GetAPIKey()
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
		sha := ""
		created := ver.Created
		updated := ver.Updated
		eaf := ver.EarlyAccessTimeFrame
		if len(ver.ModelFiles) > 0 {
			sizeKB = ver.ModelFiles[0].SizeKB
			sha = ver.ModelFiles[0].Hashes.SHA256
		}

		versions = append(versions, VersionInfo{
			ID:                   ver.ID,
			Name:                 ver.Name,
			BaseModel:            ver.BaseModel,
			SizeKB:               sizeKB,
			TrainedWords:         ver.TrainedWords,
			EarlyAccessTimeFrame: eaf,
			SHA256:               sha,
			Created:              created,
			Updated:              updated,
		})
	}

	c.JSON(200, versions)
}

func SyncVersionByID(c *gin.Context) {
	apiKey := database.GetAPIKey()
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

	modelData, _ := FetchCivitModel(apiKey, verData.ModelID)

	var model models.Model
	database.DB.Unscoped().Where("civit_id = ?", verData.ModelID).Find(&model)
	if model.ID == 0 {
		model = models.Model{
			CivitID: modelData.ID,
			Name:    modelData.Name,
			Type:    modelData.Type,
		}
		database.DB.Create(&model)
	} else if model.Type == "" {
		// ensure type is populated for older records
		model.Type = modelData.Type
		database.DB.Save(&model)
	}

	var filePath, imagePath string
	var imgW, imgH int
	var fileSHA string
	var downloadURL string
	modelType := model.Type
	if modelType == "" {
		modelType = modelData.Type
	}
	if len(verData.ModelFiles) > 0 {
		downloadURL = verData.ModelFiles[0].DownloadURL
		filePath, _ = DownloadFile(downloadURL, filepath.Join(database.GetModelsPath(), modelType), verData.ModelFiles[0].Name)
		rel, _ := filepath.Rel(database.GetModelsPath(), filePath)
		filePath = rel
		if info, err := os.Stat(filepath.Join(database.GetModelsPath(), filePath)); err == nil && info.Size() < 110 {
			os.Remove(filepath.Join(database.GetModelsPath(), filePath))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Downloaded file too small"})
			return
		}
		fileSHA = verData.ModelFiles[0].Hashes.SHA256
	}

	versionRecord := models.Version{
		ModelID:              model.ID,
		VersionID:            verData.ID,
		Name:                 verData.Name,
		BaseModel:            verData.BaseModel,
		EarlyAccessTimeFrame: verData.EarlyAccessTimeFrame,
		SizeKB:               verData.ModelFiles[0].SizeKB,
		TrainedWords:         strings.Join(verData.TrainedWords, ","),
		Nsfw:                 modelData.Nsfw,
		Type:                 modelData.Type,
		Tags:                 strings.Join(modelData.Tags, ","),
		Description:          modelData.Description,
		Mode:                 modelData.Mode,
		ModelURL:             fmt.Sprintf("https://civitai.com/models/%d?modelVersionId=%d", verData.ModelID, verData.ID),
		CivitCreatedAt:       verData.Created,
		CivitUpdatedAt:       verData.Updated,
		SHA256:               fileSHA,
		DownloadURL:          downloadURL,
		FilePath:             filePath,
	}
	database.DB.Create(&versionRecord)

	for idx, img := range verData.Images {
		imageURL := img.URL
		if imageURL == "" {
			imageURL = img.URLSmall
		}
		if imageURL == "" {
			continue
		}
		imgPath, _ := DownloadFile(imageURL, filepath.Join(database.GetImagesPath(), modelType), fmt.Sprintf("%d_%d.jpg", verData.ID, idx))
		relImg, _ := filepath.Rel(database.GetImagesPath(), imgPath)
		imgPath = relImg
		w, h, _ := GetImageDimensions(filepath.Join(database.GetImagesPath(), imgPath))
		hash, _ := FileHash(filepath.Join(database.GetImagesPath(), imgPath))
		metaBytes, _ := json.Marshal(img.Meta)
		database.DB.Create(&models.VersionImage{
			VersionID: versionRecord.ID,
			Path:      imgPath,
			Width:     w,
			Height:    h,
			Hash:      hash,
			Meta:      string(metaBytes),
		})
		if idx == 0 {
			imagePath = imgPath
			imgW = w
			imgH = h
		}
	}

	versionRecord.ImagePath = imagePath
	database.DB.Save(&versionRecord)

	if model.ImagePath == "" && imagePath != "" {
		model.ImagePath = imagePath
		model.ImageWidth = imgW
		model.ImageHeight = imgH
	}
	if model.FilePath == "" && filePath != "" {
		model.FilePath = filePath
	}
	database.DB.Save(&model)

	c.JSON(200, gin.H{"message": "Version synced", "versionId": verData.ID})
}

func processModels(items []CivitModel, apiKey string) {
	for _, item := range items {
		var existing models.Model
		database.DB.Where("civit_id = ?", item.ID).Find(&existing)
		if existing.ID == 0 {
			existing = models.Model{
				CivitID: item.ID,
				Name:    item.Name,
				Type:    item.Type,
			}
			database.DB.Create(&existing)
		} else if existing.Type == "" {
			// populate missing type on older records
			existing.Type = item.Type
			database.DB.Save(&existing)
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
			var fileSHA string
			var downloadURL string
			if len(verData.ModelFiles) > 0 {
				downloadURL = verData.ModelFiles[0].DownloadURL
				fileName := verData.ModelFiles[0].Name
				filePath, _ = DownloadFile(downloadURL, filepath.Join(database.GetModelsPath(), item.Type), fileName)
				rel, _ := filepath.Rel(database.GetModelsPath(), filePath)
				filePath = rel
				if info, err := os.Stat(filepath.Join(database.GetModelsPath(), filePath)); err == nil && info.Size() < 110 {
					os.Remove(filepath.Join(database.GetModelsPath(), filePath))
					log.Printf("downloaded %s is too small", fileName)
					continue
				}
				fileSHA = verData.ModelFiles[0].Hashes.SHA256
			}

			versionRec := models.Version{
				ModelID:              existing.ID,
				VersionID:            verData.ID,
				Name:                 verData.Name,
				BaseModel:            verData.BaseModel,
				EarlyAccessTimeFrame: verData.EarlyAccessTimeFrame,
				SizeKB:               verData.ModelFiles[0].SizeKB,
				TrainedWords:         strings.Join(verData.TrainedWords, ","),
				Nsfw:                 item.Nsfw,
				Type:                 item.Type,
				Tags:                 strings.Join(item.Tags, ","),
				Description:          item.Description,
				Mode:                 item.Mode,
				ModelURL:             fmt.Sprintf("https://civitai.com/models/%d?modelVersionId=%d", item.ID, verData.ID),
				CivitCreatedAt:       verData.Created,
				CivitUpdatedAt:       verData.Updated,
				SHA256:               fileSHA,
				DownloadURL:          downloadURL,
				FilePath:             filePath,
			}
			database.DB.Create(&versionRec)

			for idx, img := range verData.Images {
				imageURL := img.URL
				if imageURL == "" {
					imageURL = img.URLSmall
				}
				if imageURL == "" {
					continue
				}
				imgPath, _ := DownloadFile(imageURL, filepath.Join(database.GetImagesPath(), item.Type), fmt.Sprintf("%d_%d.jpg", verData.ID, idx))
				relImg, _ := filepath.Rel(database.GetImagesPath(), imgPath)
				imgPath = relImg
				w, h, _ := GetImageDimensions(filepath.Join(database.GetImagesPath(), imgPath))
				hash, _ := FileHash(filepath.Join(database.GetImagesPath(), imgPath))
				metaBytes, _ := json.Marshal(img.Meta)
				database.DB.Create(&models.VersionImage{
					VersionID: versionRec.ID,
					Path:      imgPath,
					Width:     w,
					Height:    h,
					Hash:      hash,
					Meta:      string(metaBytes),
				})
				if idx == 0 {
					imagePath = imgPath
					imgW = w
					imgH = h
				}
			}

			versionRec.ImagePath = imagePath
			database.DB.Save(&versionRec)

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
		os.Remove(filepath.Join(database.GetModelsPath(), model.FilePath))
	}
	if model.ImagePath != "" {
		os.Remove(filepath.Join(database.GetImagesPath(), model.ImagePath))
	}
	for _, v := range model.Versions {
		if v.FilePath != "" {
			os.Remove(filepath.Join(database.GetModelsPath(), v.FilePath))
		}
		if v.ImagePath != "" {
			os.Remove(filepath.Join(database.GetImagesPath(), v.ImagePath))
		}
		var imgs []models.VersionImage
		database.DB.Where("version_id = ?", v.ID).Find(&imgs)
		for _, img := range imgs {
			if img.Path != "" {
				os.Remove(filepath.Join(database.GetImagesPath(), img.Path))
			}
		}
		database.DB.Where("version_id = ?", v.ID).Delete(&models.VersionImage{})
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
	if err := database.DB.Preload("Images").First(&version, id).Error; err != nil {
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
		os.Remove(filepath.Join(database.GetModelsPath(), version.FilePath))
	}
	if version.ImagePath != "" {
		os.Remove(filepath.Join(database.GetImagesPath(), version.ImagePath))
	}
	var imgs []models.VersionImage
	database.DB.Where("version_id = ?", version.ID).Find(&imgs)
	for _, img := range imgs {
		if img.Path != "" {
			os.Remove(filepath.Join(database.GetImagesPath(), img.Path))
		}
	}
	database.DB.Where("version_id = ?", version.ID).Delete(&models.VersionImage{})

	database.DB.Unscoped().Delete(&models.Version{}, version.ID)

	var remaining int64
	database.DB.Model(&models.Version{}).Where("model_id = ?", version.ModelID).Count(&remaining)
	if remaining == 0 {
		var model models.Model
		if err := database.DB.First(&model, version.ModelID).Error; err == nil {
			if model.FilePath != "" && model.FilePath == version.FilePath {
				os.Remove(filepath.Join(database.GetModelsPath(), model.FilePath))
				model.FilePath = ""
			}
			if model.ImagePath != "" && model.ImagePath == version.ImagePath {
				os.Remove(filepath.Join(database.GetImagesPath(), model.ImagePath))
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
	version.EarlyAccessTimeFrame = input.EarlyAccessTimeFrame
	version.SizeKB = input.SizeKB
	version.TrainedWords = input.TrainedWords
	version.Nsfw = input.Nsfw
	version.Type = input.Type
	version.Tags = input.Tags
	version.Description = input.Description
	version.Mode = input.Mode
	version.ModelURL = input.ModelURL
	version.ImagePath = input.ImagePath
	version.FilePath = input.FilePath

	if err := database.DB.Save(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update version"})
		return
	}

	c.JSON(http.StatusOK, version)
}

// RefreshVersion pulls updated data from CivitAI for the specified version.
// The fields query parameter controls which parts to update. Acceptable values
// are comma-separated list of "metadata", "description", and "images". The
// value "all" updates everything.
func RefreshVersion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	fields := c.DefaultQuery("fields", "all")
	if err := refreshVersionData(id, fields); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Version refreshed"})
}

// SetVersionMainImage sets the ImagePath of a version to the specified image
// ID. If the model currently references the old ImagePath, it will be updated
// to the new one as well.
func SetVersionMainImage(c *gin.Context) {
	verIDStr := c.Param("id")
	verID, err := strconv.Atoi(verIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	imgIDStr := c.Param("imageId")
	imgID, err := strconv.Atoi(imgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	var version models.Version
	if err := database.DB.First(&version, verID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	var image models.VersionImage
	if err := database.DB.First(&image, imgID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}
	if image.VersionID != version.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image does not belong to this version"})
		return
	}

	oldPath := version.ImagePath
	version.ImagePath = image.Path
	if err := database.DB.Save(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update version"})
		return
	}

	var model models.Model
	if err := database.DB.First(&model, version.ModelID).Error; err == nil {
		if model.ImagePath == oldPath {
			model.ImagePath = image.Path
			model.ImageWidth = image.Width
			model.ImageHeight = image.Height
			database.DB.Save(&model)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Main image updated"})
}
