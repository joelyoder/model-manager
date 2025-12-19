package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

var imgTagRegex = regexp.MustCompile(`<img[^>]+src="([^">]+)"`)

// ArchiveDescriptionImages parses the HTML description, finds external images,
// downloads them to a local archive directory, rewrites the HTML to point to
// the local files, and returns the modified description.
// It returns true if changes were made.
func ArchiveDescriptionImages(versionID int, description string) (string, bool) {
	if description == "" {
		return "", false
	}

	changed := false
	matches := imgTagRegex.FindAllStringSubmatch(description, -1)
	if len(matches) == 0 {
		return description, false
	}

	// Prepare archive directory: <ImagesPath>/archives/<VersionID>/
	archivesDir := filepath.Join(database.GetImagePath(), "archives", fmt.Sprintf("%d", versionID))

	for _, match := range matches {
		src := match[1]

		// Skip if already local or relative
		if !strings.HasPrefix(src, "http") {
			continue
		}

		// Ensure archive dir exists
		if _, err := os.Stat(archivesDir); os.IsNotExist(err) {
			if err := os.MkdirAll(archivesDir, os.ModePerm); err != nil {
				log.Printf("Failed to create archive dir for version %d: %v", versionID, err)
				continue
			}
		}

		// Generate filename from URL
		// Simple strategy: hash or just base name. Let's use base name but handle query params.
		filename := filepath.Base(src)
		if idx := strings.Index(filename, "?"); idx != -1 {
			filename = filename[:idx]
		}
		// If filename is empty or weird (e.g. ends in /), generate one
		if filename == "" || filename == "." || filename == "/" {
			filename = fmt.Sprintf("img_%d.jpg", len(strings.Split(description, src))) // simple collision avoidance attempt
		}

		// Clean filename
		filename = strings.ReplaceAll(filename, "%20", "_")

		destPath := filepath.Join(archivesDir, filename)

		// Download if not exists
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			_, _, err := DownloadFile(src, archivesDir, filename)
			if err != nil {
				log.Printf("Failed to archive image %s: %v", src, err)
				continue
			}
			log.Printf("Archived image for version %d: %s -> %s", versionID, src, filename)
		}

		// Calculate relative path for the src attribute
		// The frontend serves /images/* from GetImagePath()
		// So we want /images/archives/<VersionID>/<filename> (or just archives/... if using relative to where it's served?)
		// The app serves /images/ mapped to the image root.
		// So the src should be "/images/archives/<VersionID>/<filename>"

		// Wait, we need to be careful about how the frontend displays this.
		// If the frontend renders the description HTML directly, relative paths might be relative to the page URL.
		// So we likely need an absolute path from the domain root, e.g. "/images/..."

		// IMPORTANT: backend/main.go serves:
		// r.GET("/images/*filepath", ... api.ResolveImagePath(relativePath) ...)
		// So a request to /images/archives/123/foo.jpg resolves to ResolveImagePath("archives/123/foo.jpg")
		// = <ImagesPath>/archives/123/foo.jpg. This is correct.

		newSrc := fmt.Sprintf("/images/archives/%d/%s", versionID, filename)

		// Replace in description
		// We use strings.Replace to replace the specific src
		// NOTE: This might be risky if the same URL appears twice but we want to replace all occurrences anyway.
		if src != newSrc {
			description = strings.ReplaceAll(description, src, newSrc)
			changed = true
		}
	}

	return description, changed
}

// ArchiveImages is a handler that iterates over all versions and archives their description images.
func ArchiveImages(c *gin.Context) {
	var versions []models.Version
	if err := database.DB.Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions"})
		return
	}

	count := 0
	updatedCount := 0

	for _, v := range versions {
		newDesc, changed := ArchiveDescriptionImages(v.VersionID, v.Description)
		if changed {
			v.Description = newDesc
			if err := database.DB.Save(&v).Error; err != nil {
				log.Printf("Failed to save version %d: %v", v.VersionID, err)
			} else {
				updatedCount++
			}
		}
		count++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Archive complete",
		"processed": count,
		"updated":   updatedCount,
	})
}
