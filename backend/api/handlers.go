package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func resolveNSFWFilter(c *gin.Context) (onlySafe bool, onlyNSFW bool) {
	filter := strings.ToLower(c.Query("nsfwFilter"))
	hideNsfw := c.Query("hideNsfw") == "1"

	switch filter {
	case "no":
		return true, false
	case "only":
		return false, true
	case "both":
		return false, false
	case "":
		if hideNsfw {
			return true, false
		}
	default:
		if hideNsfw {
			return true, false
		}
	}

	return false, false
}

// GetModels returns a paginated list of models, optionally filtered by query
// parameters. Supported query params include page, limit, search, baseModel,
// modelType, nsfwFilter, tags, and includeVersions. The handler performs read-only
// database queries and responds with JSON containing the matching models.
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
	baseModel := c.Query("baseModel")
	modelType := c.Query("modelType")
	onlySafe, onlyNSFW := resolveNSFWFilter(c)
	tags := c.Query("tags")
	synced := c.Query("synced") == "1"

	var modelsList []models.Model
	q := database.DB.Model(&models.Model{})

	// Filter models by versions when filters are provided or when searching
	needJoin := search != "" || baseModel != "" || modelType != "" || onlySafe || onlyNSFW || tags != "" || synced
	if needJoin {
		q = q.Joins("JOIN versions ON versions.model_id = models.id")
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(models.name) LIKE ? OR LOWER(versions.name) LIKE ? OR LOWER(versions.trained_words) LIKE ?", like, like, like)
	}
	if baseModel != "" {
		q = q.Where("versions.base_model = ?", baseModel)
	}
	if modelType != "" {
		q = q.Where("versions.type = ?", modelType)
	}
	if onlySafe {
		q = q.Where("versions.nsfw = 0")
	} else if onlyNSFW {
		q = q.Where("versions.nsfw = 1")
	}
	if tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(strings.ToLower(t))
			if t != "" {
				q = q.Where("LOWER(versions.tags) LIKE ?", "%"+t+"%")
			}
		}
	}
	if synced {
		q = q.Joins("JOIN client_files ON client_files.model_version_id = versions.id").Where("client_files.status = ?", "installed")
	}

	if c.DefaultQuery("includeVersions", "1") == "1" {
		q = q.Preload("Versions", func(db *gorm.DB) *gorm.DB {
			if baseModel != "" {
				db = db.Where("base_model = ?", baseModel)
			}
			if modelType != "" {
				db = db.Where("type = ?", modelType)
			}
			if onlySafe {
				db = db.Where("nsfw = 0")
			} else if onlyNSFW {
				db = db.Where("nsfw = 1")
			}
			if tags != "" {
				for _, t := range strings.Split(tags, ",") {
					t = strings.TrimSpace(strings.ToLower(t))
					if t != "" {
						db = db.Where("LOWER(tags) LIKE ?", "%"+t+"%")
					}
				}
			}
			return db.Order("id DESC")
		})
	}

	if needJoin {
		q = q.Group("models.id")
	}

	q.Order("models.id DESC").Limit(limit).Offset((page - 1) * limit).Find(&modelsList)
	populateClientStatus(modelsList)
	c.JSON(http.StatusOK, modelsList)
}

// GetModelsCount mirrors GetModels filtering logic but returns only the total
// count. It honors the search, baseModel, modelType, nsfwFilter, and tags query
// parameters and does not modify any database records.
func GetModelsCount(c *gin.Context) {
	search := c.Query("search")
	baseModel := c.Query("baseModel")
	modelType := c.Query("modelType")
	onlySafe, onlyNSFW := resolveNSFWFilter(c)
	tags := c.Query("tags")
	synced := c.Query("synced") == "1"

	var count int64
	q := database.DB.Model(&models.Model{})
	needJoin := search != "" || baseModel != "" || modelType != "" || onlySafe || onlyNSFW || tags != "" || synced
	if needJoin {
		q = q.Joins("JOIN versions ON versions.model_id = models.id")
	}
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(models.name) LIKE ? OR LOWER(versions.name) LIKE ? OR LOWER(versions.trained_words) LIKE ?", like, like, like)
	}
	if baseModel != "" {
		q = q.Where("versions.base_model = ?", baseModel)
	}
	if modelType != "" {
		q = q.Where("versions.type = ?", modelType)
	}
	if onlySafe {
		q = q.Where("versions.nsfw = 0")
	} else if onlyNSFW {
		q = q.Where("versions.nsfw = 1")
	}
	if tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(strings.ToLower(t))
			if t != "" {
				q = q.Where("LOWER(versions.tags) LIKE ?", "%"+t+"%")
			}
		}
	}
	if synced {
		q = q.Joins("JOIN client_files ON client_files.model_version_id = versions.id").Where("client_files.status = ?", "installed")
	}
	if needJoin {
		q = q.Group("models.id")
	}
	q.Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// GetBaseModels enumerates the distinct base_model values from stored versions.
// It accepts no parameters and simply reads the database to return an ordered
// list of strings.
func GetBaseModels(c *gin.Context) {
	var baseModels []string
	database.DB.Model(&models.Version{}).
		Distinct().
		Order("base_model").
		Pluck("base_model", &baseModels)
	c.JSON(http.StatusOK, baseModels)
}

// GetModel loads a specific model and its versions identified by the :id path
// parameter. The handler validates the ID, reads the database, and returns the
// model document without modifying persisted data.
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

	// Populate ClientStatus using helper
	tmp := []models.Model{model}
	populateClientStatus(tmp)
	model = tmp[0]

	c.JSON(http.StatusOK, model)
}

// SyncCivitModels pulls the latest models from CivitAI using the configured API
// token. The handler accepts no parameters, fetches the remote catalog, and
// calls processModels which creates database records and downloads assets to
// disk as a side effect. A JSON status message is returned when syncing
// completes.
func SyncCivitModels(c *gin.Context) {
	apiKey := getCivitaiAPIKey()
	items, err := FetchCivitModels(apiKey)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch models"})
		return
	}
	processModels(items, apiKey)
	c.JSON(200, gin.H{"message": "Models synced successfully."})
}

// SyncCivitModelByID refreshes a single CivitAI model identified by the :id
// path parameter. It fetches the remote model payload, processes the model into
// local database records, and may download associated files/images.
func SyncCivitModelByID(c *gin.Context) {
	log.Println("Hit /api/sync/:id")
	apiKey := getCivitaiAPIKey()
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

// GetModelVersions returns remote version metadata for the CivitAI model whose
// ID is provided in the :id path parameter. The handler fetches the model data
// from CivitAI, extracts the version summaries, and returns a simplified slice
// of version information without persisting it locally.
func GetModelVersions(c *gin.Context) {
	apiKey := getCivitaiAPIKey()
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

	versions := make([]VersionInfo, 0, len(model.ModelVersions))
	for _, ver := range model.ModelVersions {
		file := selectModelFile(ver.Files)

		versions = append(versions, VersionInfo{
			ID:                   ver.ID,
			ModelID:              model.ID,
			Name:                 ver.Name,
			BaseModel:            ver.BaseModel,
			SizeKB:               file.SizeKB,
			TrainedWords:         ver.TrainedWords,
			EarlyAccessTimeFrame: ver.EarlyAccessTimeFrame,
			SHA256:               file.Hashes.SHA256,
			Created:              ver.Created,
			Updated:              ver.Updated,
		})
	}

	c.JSON(200, versions)
}

var errVersionSummaryNotFound = errors.New("version not found in model summary")

// fetchVersionDetails retrieves the detailed version payload for the supplied
// versionID. It first attempts to query the dedicated version endpoint. When
// that fails and a fallback model ID is provided, it loads the model summary
// and extracts the matching version entry. The returned data mirrors the
// VersionResponse shape expected by SyncVersionByID.
func fetchVersionDetails(apiKey string, versionID int, fallbackModelID int) (VersionResponse, error) {
	version, err := FetchModelVersion(apiKey, versionID)
	if err == nil {
		return version, nil
	}

	if fallbackModelID == 0 {
		return VersionResponse{}, err
	}

	model, modelErr := FetchCivitModel(apiKey, fallbackModelID)
	if modelErr != nil {
		return VersionResponse{}, err
	}

	for _, summary := range model.ModelVersions {
		if summary.ID == versionID {
			return VersionResponse{
				ID:                   summary.ID,
				ModelID:              model.ID,
				Name:                 summary.Name,
				BaseModel:            summary.BaseModel,
				Created:              summary.Created,
				Updated:              summary.Updated,
				EarlyAccessTimeFrame: summary.EarlyAccessTimeFrame,
				TrainedWords:         summary.TrainedWords,
				ModelFiles:           summary.Files,
				Images:               summary.Images,
			}, nil
		}
	}

	return VersionResponse{}, errVersionSummaryNotFound
}

func collectVersionImages(_ string, verData VersionResponse) []ModelImage {
	return verData.Images
}

// SyncVersionByID imports a specific CivitAI model version identified by the
// :versionId path parameter. The optional download query parameter controls
// whether associated files are downloaded. The handler creates or updates local
// model/version records and may write downloaded assets to disk.
func SyncVersionByID(c *gin.Context) {
	apiKey := getCivitaiAPIKey()
	versionID := c.Param("versionId")
	downloadParam := c.DefaultQuery("download", "1")
	shouldDownload := downloadParam != "0"

	id, err := strconv.Atoi(versionID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid version ID"})
		return
	}

	var fallbackModelID int
	if modelParam := c.Query("modelId"); modelParam != "" {
		fallbackModelID, err = strconv.Atoi(modelParam)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid model ID"})
			return
		}
	}

	verData, err := fetchVersionDetails(apiKey, id, fallbackModelID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to fetch version"
		if errors.Is(err, errVersionSummaryNotFound) {
			status = http.StatusNotFound
			message = "Version not found"
		}
		c.JSON(status, gin.H{"error": message})
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
			Weight:  1,
		}
		database.DB.Create(&model)
	} else {
		updated := false
		if model.Type == "" {
			// ensure type is populated for older records
			model.Type = modelData.Type
			updated = true
		}
		if model.Weight <= 0 {
			model.Weight = 1
			updated = true
		}
		if updated {
			database.DB.Save(&model)
		}
	}

	var filePath, imagePath string
	var size int64
	var imgW, imgH int
	var fileSHA string
	var downloadURL string
	var selectedFile ModelFile
	modelType := model.Type
	if modelType == "" {
		modelType = modelData.Type
	}
	if len(verData.ModelFiles) > 0 {
		selectedFile = selectModelFile(verData.ModelFiles)
		downloadURL = selectedFile.DownloadURL
		fileName := selectedFile.Name
		destDir := filepath.Join(database.GetModelPath(), modelType)
		destPath := filepath.Join(destDir, fileName)
		if shouldDownload {
			if _, err := os.Stat(destPath); err == nil {
				ext := filepath.Ext(fileName)
				base := strings.TrimSuffix(fileName, ext)
				fileName = fmt.Sprintf("%s_%d%s", base, verData.ID, ext)
			}
			var err error
			filePath, size, err = DownloadFile(downloadURL, destDir, fileName)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					c.JSON(http.StatusConflict, gin.H{"error": "Download cancelled"})
				} else {
					log.Printf("failed to download file: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download file"})
				}
				return
			}
			if size < 110 {
				if filePath != "" {
					moveToTrash(filePath)
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Downloaded file too small"})
				return
			}
		}
		fileSHA = selectedFile.Hashes.SHA256
	}

	versionRecord := models.Version{
		ModelID:              model.ID,
		VersionID:            verData.ID,
		Name:                 verData.Name,
		BaseModel:            verData.BaseModel,
		EarlyAccessTimeFrame: verData.EarlyAccessTimeFrame,
		SizeKB:               selectedFile.SizeKB,
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
		FilePath:             MakeRelativePath(filePath, database.GetModelPath()),
	}
	// Archive images in description
	if newDesc, changed := ArchiveDescriptionImages(verData.ID, versionRecord.Description); changed {
		versionRecord.Description = newDesc
	}
	database.DB.Create(&versionRecord)

	images := collectVersionImages(apiKey, verData)
	for idx, img := range images {
		imageURL := img.URL
		if imageURL == "" {
			imageURL = img.URLSmall
		}
		if imageURL == "" {
			continue
		}
		if isVideoURL(imageURL) {
			continue
		}
		imgPath, _, _ := DownloadFile(imageURL, filepath.Join(database.GetImagePath(), modelType), fmt.Sprintf("%d_%d.jpg", verData.ID, idx))
		w, h, _ := GetImageDimensions(imgPath)
		hash, _ := FileHash(imgPath)
		metaBytes, _ := json.Marshal(img.Meta)
		database.DB.Create(&models.VersionImage{
			VersionID: versionRecord.ID,
			Path:      MakeRelativePath(imgPath, database.GetImagePath()),
			Width:     w,
			Height:    h,
			Hash:      hash,
			Meta:      string(metaBytes),
		})
		if imagePath == "" {
			imagePath = imgPath
			imgW = w
			imgH = h
		}
	}

	versionRecord.ImagePath = MakeRelativePath(imagePath, database.GetImagePath())
	database.DB.Save(&versionRecord)

	if model.ImagePath == "" && imagePath != "" {
		model.ImagePath = MakeRelativePath(imagePath, database.GetImagePath())
		model.ImageWidth = imgW
		model.ImageHeight = imgH
	}
	if model.FilePath == "" && filePath != "" {
		model.FilePath = MakeRelativePath(filePath, database.GetModelPath())
	}
	database.DB.Save(&model)

	c.JSON(200, gin.H{"message": "Version synced", "versionId": verData.ID})
}

func processModel(item CivitModel, apiKey string) {
	var existing models.Model
	database.DB.Where("civit_id = ?", item.ID).Find(&existing)
	if existing.ID == 0 {
		existing = models.Model{
			CivitID: item.ID,
			Name:    item.Name,
			Type:    item.Type,
			Weight:  1,
		}
		database.DB.Create(&existing)
	} else {
		updated := false
		if existing.Type == "" {
			// populate missing type on older records
			existing.Type = item.Type
			updated = true
		}
		if existing.Weight <= 0 {
			existing.Weight = 1
			updated = true
		}
		if updated {
			database.DB.Save(&existing)
		}
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
		var size int64
		var imgW, imgH int
		var fileSHA string
		var downloadURL string
		var selectedFile ModelFile
		if len(verData.ModelFiles) > 0 {
			selectedFile = selectModelFile(verData.ModelFiles)
			downloadURL = selectedFile.DownloadURL
			fileName := selectedFile.Name
			destDir := filepath.Join(database.GetModelPath(), item.Type)
			destPath := filepath.Join(destDir, fileName)
			if _, err := os.Stat(destPath); err == nil {
				ext := filepath.Ext(fileName)
				base := strings.TrimSuffix(fileName, ext)
				fileName = fmt.Sprintf("%s_%d%s", base, verData.ID, ext)
			}
			var err error
			filePath, size, err = DownloadFile(downloadURL, destDir, fileName)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					log.Printf("download cancelled for version %d", verData.ID)
					return
				}
				log.Printf("failed to download file: %v", err)
				continue
			}
			if size < 110 {
				if filePath != "" {
					moveToTrash(filePath)
				}
				log.Printf("downloaded %s is too small", fileName)
				continue
			}
			fileSHA = selectedFile.Hashes.SHA256
		}

		versionRec := models.Version{
			ModelID:              existing.ID,
			VersionID:            verData.ID,
			Name:                 verData.Name,
			BaseModel:            verData.BaseModel,
			EarlyAccessTimeFrame: verData.EarlyAccessTimeFrame,
			SizeKB:               selectedFile.SizeKB,
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
			FilePath:             MakeRelativePath(filePath, database.GetModelPath()),
		}
		// Archive images in description
		if newDesc, changed := ArchiveDescriptionImages(verData.ID, versionRec.Description); changed {
			versionRec.Description = newDesc
		}
		database.DB.Create(&versionRec)

		images := collectVersionImages(apiKey, verData)
		for idx, img := range images {
			imageURL := img.URL
			if imageURL == "" {
				imageURL = img.URLSmall
			}
			if imageURL == "" {
				continue
			}
			if isVideoURL(imageURL) {
				continue
			}
			imgPath, _, _ := DownloadFile(imageURL, database.GetImagePath(), fmt.Sprintf("%d_%d.jpg", verData.ID, idx))
			w, h, _ := GetImageDimensions(imgPath)
			hash, _ := FileHash(imgPath)
			metaBytes, _ := json.Marshal(img.Meta)
			database.DB.Create(&models.VersionImage{
				VersionID: versionRec.ID,
				Path:      MakeRelativePath(imgPath, database.GetImagePath()),
				Width:     w,
				Height:    h,
				Hash:      hash,
				Meta:      string(metaBytes),
			})
			if imagePath == "" {
				imagePath = imgPath
				imgW = w
				imgH = h
			}
		}

		versionRec.ImagePath = MakeRelativePath(imagePath, database.GetImagePath())
		database.DB.Save(&versionRec)

		if existing.ImagePath == "" && imagePath != "" {
			existing.ImagePath = MakeRelativePath(imagePath, database.GetImagePath())
			existing.ImageWidth = imgW
			existing.ImageHeight = imgH
		}
		if existing.FilePath == "" && filePath != "" {
			existing.FilePath = MakeRelativePath(filePath, database.GetModelPath())
		}
		database.DB.Save(&existing)
	}
}

func processModels(items []CivitModel, apiKey string) {
	sem := make(chan struct{}, 4)
	var wg sync.WaitGroup
	for _, item := range items {
		itemCopy := item
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			processModel(itemCopy, apiKey)
			<-sem
		}()
	}
	wg.Wait()
}

// DeleteModel removes the model identified by the :id path parameter along
// with all related versions, version images, and associated files. Files are
// moved to the trash directory when possible and the corresponding database
// rows are permanently deleted.
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
		moveToTrash(ResolveModelPath(model.FilePath))
	}
	if model.ImagePath != "" {
		moveToTrash(ResolveImagePath(model.ImagePath))
	}
	for _, v := range model.Versions {
		if v.FilePath != "" {
			moveToTrash(ResolveModelPath(v.FilePath))
		}
		if v.ImagePath != "" {
			moveToTrash(ResolveImagePath(v.ImagePath))
		}
		var imgs []models.VersionImage
		database.DB.Where("version_id = ?", v.ID).Find(&imgs)
		for _, img := range imgs {
			if img.Path != "" {
				moveToTrash(ResolveImagePath(img.Path))
			}
		}
		database.DB.Where("version_id = ?", v.ID).Delete(&models.VersionImage{})

		// Remove archived images directory
		archiveDir := filepath.Join(database.GetImagePath(), "archives", fmt.Sprintf("%d", v.VersionID))
		if _, err := os.Stat(archiveDir); err == nil {
			moveToTrash(archiveDir)
		}
	}

	database.DB.Unscoped().Where("model_id = ?", model.ID).Delete(&models.Version{})
	database.DB.Unscoped().Delete(&model)

	c.JSON(http.StatusOK, gin.H{"message": "Model deleted"})
}

// GetVersion retrieves the version referenced by the :id path parameter and
// returns it alongside its parent model. The handler performs read-only lookups
// and responds with JSON describing both records.
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

	// Populate ClientStatus for the single version
	var cf models.ClientFile
	if err := database.DB.Where("model_version_id = ?", version.ID).First(&cf).Error; err == nil {
		version.ClientStatus = cf.Status
	}

	c.JSON(http.StatusOK, gin.H{"model": model, "version": version})
}

// DeleteVersion removes the version addressed by the :id path parameter. The
// optional files query parameter (defaults to 1) controls whether referenced
// files and images are moved to the trash. Related database records are deleted
// as part of the operation.
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

	deleteFiles := c.DefaultQuery("files", "1") != "0"

	var imgs []models.VersionImage
	database.DB.Where("version_id = ?", version.ID).Find(&imgs)

	if deleteFiles {
		if version.FilePath != "" {
			moveToTrash(ResolveModelPath(version.FilePath))
		}
		if version.ImagePath != "" {
			moveToTrash(ResolveImagePath(version.ImagePath))
		}
		for _, img := range imgs {
			if img.Path != "" {
				moveToTrash(ResolveImagePath(img.Path))
			}
		}

		// Remove archived images directory
		archiveDir := filepath.Join(database.GetImagePath(), "archives", fmt.Sprintf("%d", version.VersionID))
		if _, err := os.Stat(archiveDir); err == nil {
			moveToTrash(archiveDir)
		}
	}

	database.DB.Where("version_id = ?", version.ID).Delete(&models.VersionImage{})

	database.DB.Unscoped().Delete(&models.Version{}, version.ID)

	if deleteFiles {
		var remaining int64
		database.DB.Model(&models.Version{}).Where("model_id = ?", version.ModelID).Count(&remaining)
		if remaining == 0 {
			var model models.Model
			if err := database.DB.First(&model, version.ModelID).Error; err == nil {
				if model.FilePath != "" && model.FilePath == version.FilePath {
					moveToTrash(model.FilePath)
					model.FilePath = ""
				}
				if model.ImagePath != "" && model.ImagePath == version.ImagePath {
					moveToTrash(model.ImagePath)
					model.ImagePath = ""
				}
				database.DB.Save(&model)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Version deleted"})
}

// CreateModel inserts a placeholder model and companion version when called.
// No request body is required; the handler generates temporary negative
// identifiers to avoid conflicts and returns the created database IDs. New
// model and version rows are written as a side effect.
func CreateModel(c *gin.Context) {
	civitID := -int(time.Now().UnixNano())
	model := models.Model{CivitID: civitID, Name: "New Model", Type: "Checkpoint", Weight: 1}
	if err := database.DB.Create(&model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create model"})
		return
	}

	verID := -int(time.Now().UnixNano())
	version := models.Version{ModelID: model.ID, VersionID: verID, Name: "v1", Type: model.Type}
	if err := database.DB.Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create version"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"modelId": model.ID, "versionId": version.ID})
}

// UpdateModel updates the model referenced by the :id path parameter using the
// JSON body of the request. All mutable fields on the model struct can be
// provided, and the handler persists the changes to the database.
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

	model.CivitID = input.CivitID
	model.Name = input.Name
	model.Type = input.Type
	model.Tags = input.Tags
	model.Nsfw = input.Nsfw
	model.Description = input.Description
	model.ImagePath = input.ImagePath
	model.FilePath = input.FilePath
	model.ImageWidth = input.ImageWidth
	model.ImageHeight = input.ImageHeight
	if input.Weight <= 0 {
		input.Weight = 1
	}
	model.Weight = input.Weight

	if err := database.DB.Save(&model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update model"})
		return
	}

	c.JSON(http.StatusOK, model)
}

// UpdateVersion applies the JSON payload to the version specified by the :id
// path parameter. It validates version ID uniqueness before saving and then
// updates the stored version record with the provided fields.
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

	if version.VersionID != input.VersionID {
		var existing models.Version
		database.DB.Where("version_id = ?", input.VersionID).Not("id = ?", version.ID).First(&existing)
		if existing.ID > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "version id already exists"})
			return
		}
	}

	version.VersionID = input.VersionID
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
	version.CivitCreatedAt = input.CivitCreatedAt
	version.CivitUpdatedAt = input.CivitUpdatedAt
	version.SHA256 = input.SHA256
	version.DownloadURL = input.DownloadURL
	version.ImagePath = input.ImagePath
	version.FilePath = input.FilePath

	if err := database.DB.Save(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update version"})
		return
	}

	c.JSON(http.StatusOK, version)
}

// RefreshVersion pulls updated data from CivitAI for the version identified by
// the :id path parameter. The optional fields query parameter (all, metadata,
// description, images) selects which sections to refresh. The handler updates
// the existing database record and may download new images when requested.
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

// SetVersionMainImage sets the main image for a version using the :id version
// path parameter and :imageId gallery image parameter. When the parent model is
// referencing the previous main image, it is updated to the new selection. The
// operation only touches database state.
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

// UploadVersionFile handles multipart uploads of primary files or preview
// images for the version specified by the :id path parameter. The "file" form
// field must contain the binary data. Query parameters "kind" (file or image)
// and "type" (model folder) determine where the file is stored. The handler
// writes the uploaded file to disk and updates the associated version/model
// paths before returning the absolute path as JSON.
func UploadVersionFile(c *gin.Context) {
	verIDStr := c.Param("id")
	verID, err := strconv.Atoi(verIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	kind := c.DefaultQuery("kind", "file")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	var version models.Version
	if err := database.DB.First(&version, verID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	var model models.Model
	if err := database.DB.First(&model, version.ModelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	modelType := c.Query("type")
	if modelType == "" {
		modelType = version.Type
	}
	if modelType == "" {
		modelType = model.Type
	}
	if modelType == "" {
		modelType = "Checkpoint"
	}

	destDir := filepath.Join(database.GetModelPath(), modelType)
	if kind == "image" {
		destDir = filepath.Join(database.GetImagePath(), modelType)
	}
	os.MkdirAll(destDir, os.ModePerm)
	filename := filepath.Base(header.Filename)
	destPath := filepath.Join(destDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	out.Close()

	absPath, _ := filepath.Abs(destPath)

	if kind == "image" {
		version.ImagePath = MakeRelativePath(absPath, database.GetImagePath())
		if model.ImagePath == "" {
			model.ImagePath = version.ImagePath
			if w, h, err := GetImageDimensions(absPath); err == nil {
				model.ImageWidth = w
				model.ImageHeight = h
			}
		}
	} else {
		version.FilePath = MakeRelativePath(absPath, database.GetModelPath())
		if model.FilePath == "" {
			model.FilePath = version.FilePath
		}
	}

	database.DB.Save(&version)
	database.DB.Save(&model)

	c.JSON(http.StatusOK, gin.H{"path": absPath})
}

// UploadVersionImage uploads a supplemental gallery image for the version given
// by the :id path parameter. The request must include a "file" form field and
// may specify a "type" query parameter to pick the destination folder. The
// image is written to disk, metadata is extracted, and a new VersionImage row
// is created.
func UploadVersionImage(c *gin.Context) {
	verIDStr := c.Param("id")
	verID, err := strconv.Atoi(verIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	var version models.Version
	if err := database.DB.First(&version, verID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	var model models.Model
	if err := database.DB.First(&model, version.ModelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	modelType := c.Query("type")
	if modelType == "" {
		modelType = version.Type
	}
	if modelType == "" {
		modelType = model.Type
	}
	if modelType == "" {
		modelType = "Checkpoint"
	}

	destDir := filepath.Join(database.GetImagePath(), modelType)
	os.MkdirAll(destDir, os.ModePerm)
	filename := filepath.Base(header.Filename)
	destPath := filepath.Join(destDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save"})
		return
	}
	out.Close()

	absPath, _ := filepath.Abs(destPath)
	w, h, _ := GetImageDimensions(absPath)
	hash, _ := FileHash(absPath)
	metaMap, _ := ExtractImageMetadata(absPath)
	metaBytes, _ := json.Marshal(metaMap)

	img := models.VersionImage{
		VersionID: version.ID,
		Path:      MakeRelativePath(absPath, database.GetImagePath()),
		Width:     w,
		Height:    h,
		Hash:      hash,
		Meta:      string(metaBytes),
	}
	database.DB.Create(&img)
	c.JSON(http.StatusOK, img)
}

// DeleteVersionImage removes the gallery image identified by the :imgId path
// parameter from the version :id. It deletes the database row and moves the
// underlying image file to the trash if present.
func DeleteVersionImage(c *gin.Context) {
	verIDStr := c.Param("id")
	verID, err := strconv.Atoi(verIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}
	imgIDStr := c.Param("imgId")
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

	if image.Path != "" {
		moveToTrash(image.Path)
	}
	database.DB.Delete(&image)

	if version.ImagePath == image.Path {
		oldPath := version.ImagePath
		version.ImagePath = ""
		database.DB.Save(&version)
		var model models.Model
		if err := database.DB.First(&model, version.ModelID).Error; err == nil {
			if model.ImagePath == oldPath {
				model.ImagePath = ""
				database.DB.Save(&model)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted"})
}

func populateClientStatus(modelsList []models.Model) {
	var versionIDs []uint
	for i := range modelsList {
		for j := range modelsList[i].Versions {
			versionIDs = append(versionIDs, modelsList[i].Versions[j].ID)
		}
	}
	if len(versionIDs) == 0 {
		return
	}

	var clientFiles []models.ClientFile
	database.DB.Where("model_version_id IN ?", versionIDs).Find(&clientFiles)

	statusMap := make(map[uint]string)
	for _, cf := range clientFiles {
		// Prioritize "installed" over "pending" if multiple clients exist
		if current, ok := statusMap[cf.ModelVersionID]; ok {
			if current != "installed" && cf.Status == "installed" {
				statusMap[cf.ModelVersionID] = "installed"
			}
		} else {
			statusMap[cf.ModelVersionID] = cf.Status
		}
	}

	for i := range modelsList {
		for j := range modelsList[i].Versions {
			if status, ok := statusMap[modelsList[i].Versions[j].ID]; ok {
				modelsList[i].Versions[j].ClientStatus = status
			}
		}
	}
}
