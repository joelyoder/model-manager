package api

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/chai2010/webp"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

// EnsureVersionThumbnail creates a thumbnail for the given version
// Format: thumbnails/v_<VersionID>.webp
func EnsureVersionThumbnail(versionID uint, sourcePath string) error {
	if sourcePath == "" {
		return nil
	}

	// Use filepath.Join for OS-correct separators
	thumbnailRelPath := filepath.Join("thumbnails", fmt.Sprintf("v_%d.webp", versionID))
	fullThumbnailPath := ResolveImagePath(thumbnailRelPath)

	// Resolve full path to source image
	fullSourcePath := ResolveImagePath(sourcePath)

	return generateThumbnail(fullSourcePath, fullThumbnailPath)
}

// Private helper to do the actual resizing/encoding
func generateThumbnail(fullSourcePath, fullThumbnailPath string) error {
	// Create thumbnails directory if it doesn't exist
	thumbDir := filepath.Dir(fullThumbnailPath)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return fmt.Errorf("failed to create thumbnail dir: %w", err)
	}

	log.Printf("Generating thumbnail from %s to %s", fullSourcePath, fullThumbnailPath)

	// Open source
	file, err := os.Open(fullSourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source image: %w", err)
	}
	defer file.Close()

	// Decode
	img, _, err := image.Decode(file)
	if err != nil {
		// Try WebP explicitly if generic decode fails
		if _, errSeek := file.Seek(0, 0); errSeek == nil {
			if imgWebp, errWebp := webp.Decode(file); errWebp == nil {
				img = imgWebp
				err = nil // success
			} else {
				return fmt.Errorf("failed to decode image (std: %v, webp: %v)", err, errWebp)
			}
		} else {
			return fmt.Errorf("failed to decode image: %w", err)
		}
	}

	// Resize if necessary
	bounds := img.Bounds()
	if bounds.Dy() > 450 {
		// Resize to height 450, preserving aspect ratio
		img = resize.Resize(0, 450, img, resize.Lanczos3)
	}

	// Create thumbnail file
	out, err := os.Create(fullThumbnailPath)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail file: %w", err)
	}
	defer out.Close()

	// Encode WebP
	if err := webp.Encode(out, img, &webp.Options{Lossless: false, Quality: 80}); err != nil {
		return fmt.Errorf("failed to encode webp: %w", err)
	}

	return nil
}

// DeleteVersionThumbnail removes the thumbnail for a version
func DeleteVersionThumbnail(versionID uint) error {
	thumbnailRelPath := filepath.Join("thumbnails", fmt.Sprintf("v_%d.webp", versionID))
	fullThumbnailPath := ResolveImagePath(thumbnailRelPath)

	if _, err := os.Stat(fullThumbnailPath); err == nil {
		log.Printf("Deleting version thumbnail: %s", fullThumbnailPath)
		return os.Remove(fullThumbnailPath)
	}
	return nil
}

// GenerateMissingThumbnails API handler to generate missing thumbnails
func GenerateMissingThumbnails(c *gin.Context) {
	count := 0
	errors := 0

	go func() {
		// Generate version thumbnails
		var allVersions []models.Version
		database.DB.Find(&allVersions)

		for _, v := range allVersions {
			if v.ImagePath == "" {
				continue
			}

			// Use filepath.Join for OS-correct separators
			thumbPath := ResolveImagePath(filepath.Join("thumbnails", fmt.Sprintf("v_%d.webp", v.ID)))

			if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
				if err := EnsureVersionThumbnail(v.ID, v.ImagePath); err != nil {
					log.Printf("Error generating thumbnail for version %d: %v", v.ID, err)
					errors++
				} else {
					count++
				}
			}
		}

		log.Printf("Finished generating thumbnails. Created: %d, Errors: %d", count, errors)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Thumbnail generation started in background"})
}
