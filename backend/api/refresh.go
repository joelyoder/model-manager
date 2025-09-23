package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"
)

// refreshVersionData updates the specified version based on the comma separated
// fields string. Valid fields are "metadata", "description" and "images". The
// value "all" updates everything.
func refreshVersionData(id int, fields string) error {
	var version models.Version
	if err := database.DB.Preload("Images").First(&version, id).Error; err != nil {
		return err
	}

	var model models.Model
	if err := database.DB.First(&model, version.ModelID).Error; err != nil {
		return err
	}

	updateMeta := false
	updateDesc := false
	updateImages := false
	f := strings.ToLower(fields)
	if f == "" || f == "all" {
		updateMeta = true
		updateDesc = true
		updateImages = true
	} else {
		for _, p := range strings.Split(f, ",") {
			switch strings.TrimSpace(p) {
			case "metadata":
				updateMeta = true
			case "description":
				updateDesc = true
			case "images":
				updateImages = true
			}
		}
	}

	apiKey := getCivitaiAPIKey()
	modelData, err := FetchCivitModel(apiKey, model.CivitID)
	if err != nil {
		return err
	}
	verData, err := fetchVersionDetails(apiKey, version.VersionID, model.CivitID)
	if err != nil {
		return err
	}

	if updateMeta {
		version.Name = verData.Name
		version.BaseModel = verData.BaseModel
		version.EarlyAccessTimeFrame = verData.EarlyAccessTimeFrame
		if len(verData.ModelFiles) > 0 {
			file := selectModelFile(verData.ModelFiles)
			version.SizeKB = file.SizeKB
			version.SHA256 = file.Hashes.SHA256
			version.DownloadURL = file.DownloadURL
		}
		version.TrainedWords = strings.Join(verData.TrainedWords, ",")
		version.Nsfw = modelData.Nsfw
		version.Type = modelData.Type
		version.Tags = strings.Join(modelData.Tags, ",")
		version.Mode = modelData.Mode
		version.ModelURL = fmt.Sprintf("https://civitai.com/models/%d?modelVersionId=%d", verData.ModelID, verData.ID)
		version.CivitCreatedAt = verData.Created
		version.CivitUpdatedAt = verData.Updated

		model.Name = modelData.Name
		model.Type = modelData.Type
		model.Tags = strings.Join(modelData.Tags, ",")
		model.Nsfw = modelData.Nsfw
	}

	if updateDesc {
		model.Description = modelData.Description
		version.Description = modelData.Description
	}

	if updateImages {
		if version.ImagePath != "" {
			moveToTrash(version.ImagePath)
		}
		var imgs []models.VersionImage
		database.DB.Where("version_id = ?", version.ID).Find(&imgs)
		for _, img := range imgs {
			if img.Path != "" {
				moveToTrash(img.Path)
			}
		}
		database.DB.Where("version_id = ?", version.ID).Delete(&models.VersionImage{})

		var imagePath string
		var imgW, imgH int
		modelType := model.Type
		if modelType == "" {
			modelType = modelData.Type
		}
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
			imgPath, _, _ := DownloadFile(imageURL, "./backend/images/"+modelType, fmt.Sprintf("%d_%d.jpg", verData.ID, idx))
			w, h, _ := GetImageDimensions(imgPath)
			hash, _ := FileHash(imgPath)
			metaBytes, _ := json.Marshal(img.Meta)
			database.DB.Create(&models.VersionImage{
				VersionID: version.ID,
				Path:      imgPath,
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

		version.ImagePath = imagePath
		if model.ImagePath == version.ImagePath || model.ImagePath == "" {
			model.ImagePath = imagePath
			model.ImageWidth = imgW
			model.ImageHeight = imgH
		}
	}

	database.DB.Save(&model)
	database.DB.Save(&version)
	return nil
}
