package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
