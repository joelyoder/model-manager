package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
)

// FetchCivitModels calls the CivitAI REST API using the provided apiKey and
// retrieves a paginated list of models. The function issues an outbound HTTP
// request, unmarshals the JSON response, and returns the parsed models or an
// error. Network I/O is the primary side effect.
func FetchCivitModels(apiKey string) ([]CivitModel, error) {
	var models []CivitModel
	url := "https://civitai.com/api/v1/models?limit=100"

	log.Printf("GET %s", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return models, fmt.Errorf("failed to fetch model list")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &models)
	return models, err
}

// FetchCivitModel retrieves details for a specific CivitAI model identified by
// modelID. The apiKey is injected as a bearer token. A single HTTP request is
// executed and the JSON body is unmarshaled into a CivitModel structure.
func FetchCivitModel(apiKey string, modelID int) (CivitModel, error) {
	var model CivitModel
	url := fmt.Sprintf("https://civitai.com/api/v1/models/%d", modelID)

	log.Printf("GET %s", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return model, fmt.Errorf("failed to fetch model %d", modelID)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &model)
	return model, err
}

// FetchModelVersion fetches metadata for a specific model version from CivitAI
// using the supplied apiKey and versionID. It performs an HTTP GET request and
// decodes the JSON response into a VersionResponse value.
func FetchModelVersion(apiKey string, versionID int) (VersionResponse, error) {
	var version VersionResponse
	url := fmt.Sprintf("https://civitai.com/api/v1/model-versions/%d", versionID)

	log.Printf("GET %s", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return version, fmt.Errorf("failed to fetch version %d", versionID)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &version)
	return version, err
}

// FetchVersionImages retrieves the images associated with a specific model version
// using the paginated CivitAI images endpoint. It aggregates all pages, filters the
// results to the supplied version/model type and returns the curated list. When
// filtering removes every entry, the original unfiltered set is returned to avoid
// silently dropping data.
func FetchVersionImages(apiKey string, versionID int, modelType string) ([]ModelImage, error) {
	var images []ModelImage
	cursor := ""
	seen := make(map[int]struct{})

	for {
		url := fmt.Sprintf("https://civitai.com/api/v1/images?modelVersionId=%d&limit=100&withMeta=true", versionID)
		if cursor != "" {
			url += "&cursor=" + neturl.QueryEscape(cursor)
		}

		log.Printf("GET %s", url)

		req, _ := http.NewRequest("GET", url, nil)
		if apiKey != "" {
			req.Header.Add("Authorization", "Bearer "+apiKey)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch images for version %d", versionID)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		var parsed imagesResponse
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, err
		}

		for _, item := range parsed.Items {
			if _, exists := seen[item.ID]; exists {
				continue
			}
			seen[item.ID] = struct{}{}
			images = append(images, item)
		}

		if parsed.Metadata.NextCursor == "" {
			break
		}
		cursor = parsed.Metadata.NextCursor
	}

	filtered := filterImagesForVersion(images, versionID, modelType)
	if len(filtered) == 0 {
		return images, nil
	}
	return filtered, nil
}

func filterImagesForVersion(images []ModelImage, versionID int, modelType string) []ModelImage {
	if versionID == 0 {
		return images
	}

	normType := normalizeModelType(modelType)
	var filtered []ModelImage

	for _, img := range images {
		if imageMatchesVersion(img, versionID, normType) {
			filtered = append(filtered, img)
		}
	}

	return filtered
}

func imageMatchesVersion(img ModelImage, versionID int, modelType string) bool {
	var hasVersionID bool
	for _, id := range img.ModelVersionIDs {
		if id == versionID {
			hasVersionID = true
			break
		}
	}

	if !hasVersionID {
		return metaReferencesVersion(img.Meta, versionID, modelType)
	}

	if modelType == "" {
		return true
	}

	if metaReferencesVersion(img.Meta, versionID, modelType) {
		return true
	}

	// When metadata is unavailable we keep the image to avoid missing previews.
	return img.Meta == nil
}

func metaReferencesVersion(meta map[string]interface{}, versionID int, modelType string) bool {
	if meta == nil {
		return false
	}

	resources, ok := meta["civitaiResources"]
	if !ok {
		return false
	}

	entries, ok := resources.([]interface{})
	if !ok {
		return false
	}

	for _, entry := range entries {
		resMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		id := numericToInt(resMap["modelVersionId"])
		if id != versionID {
			continue
		}

		if modelType == "" {
			return true
		}

		if typeVal, ok := resMap["type"].(string); ok {
			if strings.EqualFold(typeVal, modelType) {
				return true
			}
		} else {
			return true
		}
	}

	return false
}

func normalizeModelType(modelType string) string {
	return strings.ToLower(strings.TrimSpace(modelType))
}

func numericToInt(value interface{}) int {
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return int(i)
		}
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return 0
}
