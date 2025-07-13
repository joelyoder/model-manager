package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
