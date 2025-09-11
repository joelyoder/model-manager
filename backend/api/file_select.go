package api

import "strings"

// selectModelFile returns the preferred model file from the provided list.
// It chooses a file with a .safetensors extension when available,
// otherwise it falls back to the first file in the slice.
func selectModelFile(files []ModelFile) ModelFile {
	if len(files) == 0 {
		return ModelFile{}
	}
	for _, f := range files {
		if strings.HasSuffix(strings.ToLower(f.Name), ".safetensors") {
			return f
		}
	}
	return files[0]
}
