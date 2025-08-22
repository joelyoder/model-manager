package api

import (
	"testing"

	"bou.ke/monkey"
	"model-manager/backend/database"
)

func TestGetCivitaiAPIKey(t *testing.T) {
	t.Run("db overrides env", func(t *testing.T) {
		t.Setenv("CIVIT_API_KEY", "envVal")
		patch := monkey.Patch(database.GetSettingValue, func(key string) string {
			if key != "civitai_api_key" {
				t.Fatalf("unexpected key: %s", key)
			}
			return "dbVal"
		})
		defer patch.Unpatch()

		if got := getCivitaiAPIKey(); got != "dbVal" {
			t.Errorf("expected dbVal, got %s", got)
		}
	})

	t.Run("env when db empty", func(t *testing.T) {
		t.Setenv("CIVIT_API_KEY", "envVal")
		patch := monkey.Patch(database.GetSettingValue, func(key string) string {
			if key != "civitai_api_key" {
				t.Fatalf("unexpected key: %s", key)
			}
			return ""
		})
		defer patch.Unpatch()

		if got := getCivitaiAPIKey(); got != "envVal" {
			t.Errorf("expected envVal, got %s", got)
		}
	})
}
