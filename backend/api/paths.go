package api

import (
	"os"
	"path/filepath"
	"strings"

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
// It also handles Windows paths when running on Linux and vice versa.
func MakeRelativePath(absPath, basePath string) string {
	if absPath == "" {
		return ""
	}

	// Normalize both paths to forward slashes for comparison
	// Note: filepath.ToSlash only converts the current OS separator, so on Linux
	// it won't convert Windows backslashes. We must explicitly replace them.
	normalizedAbs := strings.ReplaceAll(absPath, "\\", "/")
	normalizedBase := strings.ReplaceAll(basePath, "\\", "/")

	// Try standard filepath.Rel first (works when both are on same OS format)
	absBase, err := filepath.Abs(basePath)
	if err == nil {
		rel, err := filepath.Rel(absBase, absPath)
		if err == nil && !filepath.IsAbs(rel) && !isParentDir(rel) {
			return filepath.ToSlash(rel)
		}
	}

	// If that failed, try to find common suffix patterns
	// This handles cases like Windows path "C:/Users/.../backend/images/LORA/file.jpg"
	// when basePath is "/mnt/tank/ai/images" - we extract "LORA/file.jpg"

	// Look for common directory markers in the path
	markers := []string{"/images/", "/downloads/", "\\images\\", "\\downloads\\"}
	for _, marker := range markers {
		if idx := strings.LastIndex(normalizedAbs, strings.ReplaceAll(marker, "\\", "/")); idx != -1 {
			// Extract everything after the marker (e.g., "LORA/file.jpg")
			suffix := normalizedAbs[idx+len(strings.ReplaceAll(marker, "\\", "/")):]
			if suffix != "" && !strings.HasPrefix(suffix, "/") && !strings.HasPrefix(suffix, "..") {
				return suffix
			}
		}
	}

	// If the path contains the base path's last component, try to extract from there
	baseDir := filepath.Base(normalizedBase)
	if baseDir != "" && baseDir != "." && baseDir != "/" {
		searchPattern := "/" + baseDir + "/"
		if idx := strings.LastIndex(normalizedAbs, searchPattern); idx != -1 {
			suffix := normalizedAbs[idx+len(searchPattern):]
			if suffix != "" && !strings.HasPrefix(suffix, "..") {
				return suffix
			}
		}
	}

	// Fallback: return the original path normalized
	return normalizedAbs
}

func isParentDir(path string) bool {
	// Simple check if path starts with .. or contains ..
	// filepath.Rel usually handles the ".." prefix well, but let's be safe
	return len(path) >= 2 && path[:2] == ".."
}
