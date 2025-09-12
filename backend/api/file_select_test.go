package api

import "testing"

// TestSelectModelFile verifies the selection logic for ModelFile slices.
func TestSelectModelFile(t *testing.T) {
	t.Run("prefers safetensors when available", func(t *testing.T) {
		files := []ModelFile{
			{Name: "model.ckpt"},
			{Name: "model.safetensors"},
			{Name: "model.pt"},
		}
		f := selectModelFile(files)
		if f.Name != "model.safetensors" {
			t.Fatalf("expected safetensors, got %s", f.Name)
		}
	})

	t.Run("falls back to first file when no safetensors", func(t *testing.T) {
		files := []ModelFile{
			{Name: "model.ckpt"},
			{Name: "model.pt"},
		}
		f := selectModelFile(files)
		if f.Name != "model.ckpt" {
			t.Fatalf("expected first file, got %s", f.Name)
		}
	})

	t.Run("returns zero value for empty slice", func(t *testing.T) {
		var files []ModelFile
		f := selectModelFile(files)
		if (f != ModelFile{}) {
			t.Fatalf("expected zero value, got %+v", f)
		}
	})

	t.Run("selects first safetensors deterministically", func(t *testing.T) {
		files := []ModelFile{
			{Name: "first.safetensors"},
			{Name: "second.safetensors"},
			{Name: "model.pt"},
		}
		f := selectModelFile(files)
		if f.Name != "first.safetensors" {
			t.Fatalf("expected first safetensors, got %s", f.Name)
		}
	})

	t.Run("unsupported extensions fall back to first file", func(t *testing.T) {
		files := []ModelFile{
			{Name: "model.txt"},
			{Name: "model.doc"},
		}
		f := selectModelFile(files)
		if f.Name != "model.txt" {
			t.Fatalf("expected first file, got %s", f.Name)
		}
	})
}
