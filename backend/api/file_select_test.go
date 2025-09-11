package api

import "testing"

func TestSelectModelFile(t *testing.T) {
	files := []ModelFile{
		{Name: "model.ckpt"},
		{Name: "model.safetensors"},
		{Name: "model.pt"},
	}
	f := selectModelFile(files)
	if f.Name != "model.safetensors" {
		t.Fatalf("expected safetensors, got %s", f.Name)
	}

	files = []ModelFile{
		{Name: "model.ckpt"},
		{Name: "model.pt"},
	}
	f = selectModelFile(files)
	if f.Name != "model.ckpt" {
		t.Fatalf("expected first file, got %s", f.Name)
	}
}
