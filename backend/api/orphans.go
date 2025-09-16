package api

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
)

// GetOrphanedFiles scans the backend/downloads directory tree and returns any
// SAFETensors/PT files not referenced by models or versions in the database.
// The handler accepts no parameters, walks the filesystem (following directory
// symlinks), and reports orphaned absolute paths as JSON. It does not modify
// database records but logs extensively and may resolve symlinks during the
// scan.
func GetOrphanedFiles(c *gin.Context) {
	// Collect file paths from models and versions
	var modelPaths []string
	database.DB.Model(&models.Model{}).Where("file_path <> ''").Pluck("file_path", &modelPaths)
	var versionPaths []string
	database.DB.Model(&models.Version{}).Where("file_path <> ''").Pluck("file_path", &versionPaths)

	dbFiles := make(map[string]struct{})
	for _, p := range append(modelPaths, versionPaths...) {
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Printf("failed to get abs path for %s: %v", p, err)
			continue
		}
		resolved, err := filepath.EvalSymlinks(abs)
		if err != nil {
			log.Printf("failed to eval symlink for %s: %v", abs, err)
			resolved = abs
		}
		key := resolved
		if runtime.GOOS == "windows" {
			key = strings.ToLower(key)
		}
		dbFiles[key] = struct{}{}
		log.Printf("db file: %s", key)
	}

	var orphans []string
	root := filepath.Join("backend", "downloads")
	absRoot, err := filepath.Abs(root)
	if err != nil {
		log.Printf("failed to resolve downloads directory %s: %v", root, err)
	}
	log.Printf("walking downloads directory: %s", absRoot)

	visited := make(map[string]struct{})
	var walkFn func(string, fs.DirEntry, error) error

	walkFn = func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("walk error on %s: %v", path, err)
			return nil
		}

		if d.Type()&fs.ModeSymlink != 0 {
			target, err := filepath.EvalSymlinks(path)
			if err != nil {
				log.Printf("failed to eval symlink %s: %v", path, err)
				return nil
			}
			info, err := os.Stat(target)
			if err != nil {
				log.Printf("stat symlink target %s: %v", target, err)
				return nil
			}
			if info.IsDir() {
				absTarget, err := filepath.Abs(target)
				if err == nil {
					if _, ok := visited[absTarget]; ok {
						return nil
					}
					visited[absTarget] = struct{}{}
				}
				log.Printf("following symlink dir: %s -> %s", path, target)
				return filepath.WalkDir(target, walkFn)
			}
			d = fs.FileInfoToDirEntry(info)
			path = target
		}

		if d.IsDir() {
			log.Printf("scanning dir: %s", path)
			return nil
		}

		log.Printf("found file: %s", path)
		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext != ".safetensors" && ext != ".pt" {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			log.Printf("failed to get abs path for %s: %v", path, err)
			return nil
		}
		key := abs
		if runtime.GOOS == "windows" {
			key = strings.ToLower(key)
		}
		if _, exists := dbFiles[key]; !exists {
			log.Printf("orphaned file: %s", abs)
			orphans = append(orphans, abs)
		} else {
			log.Printf("file referenced in db: %s", abs)
		}
		return nil
	}

	filepath.WalkDir(absRoot, walkFn)

	c.JSON(http.StatusOK, gin.H{"orphans": orphans})
}
