package api

import (
	"os"
	"path/filepath"

	"model-manager/backend/database"
)

// ResolveModelPath resolves a given path against the configured model storage
// directory. If the path is absolute and exists, it returns it as-is for
// backward compatibility. Otherwise, it returns the path joined with the configured
// root.
func ResolveModelPath(path string) string {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	root := database.GetModelPath()
	return filepath.Join(root, path)
}

// ResolveImagePath resolves a given path against the configured images storage
// directory. If the path is absolute and exists, it returns it as-is for
// backward compatibility. Otherwise, it returns the path joined with the configured
// root.
func ResolveImagePath(path string) string {
	if filepath.IsAbs(path) {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	root := database.GetImagePath()
	return filepath.Join(root, path)
}

// MakeRelativePath converts an absolute path to be relative to the provided base
// path if possible. It normalizes separators to forward slashes for cross-platform
// consistency in the database.
func MakeRelativePath(absPath, basePath string) string {
	absBase, err := filepath.Abs(basePath)
	if err != nil {
		absBase = basePath
	}
	rel, err := filepath.Rel(absBase, absPath)
	if err == nil && !filepath.IsAbs(rel) && rel != ".." && !isParentDir(rel) {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(absPath)
}

func isParentDir(path string) bool {
	// Simple check if path starts with .. or contains ..
	// filepath.Rel usually handles the ".." prefix well, but let's be safe
	return len(path) >= 2 && path[:2] == ".."
}
