package api

import (
	"os"

	"model-manager/backend/database"
)

func getCivitaiAPIKey() string {
	key := database.GetSettingValue("civitai_api_key")
	if key == "" {
		// fallback to env var for compatibility
		return os.Getenv("CIVIT_API_KEY")
	}
	return key
}
