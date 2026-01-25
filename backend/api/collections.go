package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"model-manager/backend/database"
	"model-manager/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCollections returns a list of collections, optionally filtered by name
func GetCollections(c *gin.Context) {
	search := c.Query("search")
	collections := make([]models.Collection, 0)

	q := database.DB.Model(&models.Collection{})
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Where("LOWER(name) LIKE ?", like)
	}

	if err := q.Order("LOWER(name) ASC").Preload("Versions").Find(&collections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch collections"})
		return
	}

	c.JSON(http.StatusOK, collections)
}

// GetCollection returns a single collection by ID
func GetCollection(c *gin.Context) {
	id := c.Param("id")
	var collection models.Collection
	if err := database.DB.Preload("Versions").First(&collection, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}
	c.JSON(http.StatusOK, collection)
}

// CreateCollection creates a new collection
func CreateCollection(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection := models.Collection{
		Name:        input.Name,
		Description: input.Description,
	}

	if err := database.DB.Create(&collection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create collection"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// UpdateCollection updates an existing collection
func UpdateCollection(c *gin.Context) {
	id := c.Param("id")
	var collection models.Collection

	if err := database.DB.First(&collection, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection.Name = input.Name
	collection.Description = input.Description

	if err := database.DB.Save(&collection).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update collection"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

// DeleteCollection deletes a collection
func DeleteCollection(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Collection{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete collection"})
		return
	}

	// Association in many-to-many is handled by GORM (records in join table removed),
	// but we should verify if cascading delete is configured or if we need to clean up manual join table.
	// With gorm:"many2many:collection_versions", GORM handles deletion from the join table.

	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted"})
}

// AddVersionToCollection adds a model version to a collection
func AddVersionToCollection(c *gin.Context) {
	collectionID := c.Param("id")
	var input struct {
		VersionID uint `json:"versionId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var collection models.Collection
	if err := database.DB.First(&collection, collectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	var version models.Version
	if err := database.DB.First(&version, input.VersionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	if err := database.DB.Model(&collection).Association("Versions").Append(&version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add version to collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to collection"})
}

// RemoveVersionFromCollection removes a model version from a collection
func RemoveVersionFromCollection(c *gin.Context) {
	collectionID := c.Param("id")
	versionID := c.Param("versionId")

	var collection models.Collection
	if err := database.DB.First(&collection, collectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	var version models.Version
	if err := database.DB.First(&version, versionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	if err := database.DB.Model(&collection).Association("Versions").Delete(&version); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove version from collection"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from collection"})
}

// GetCollectionVersions returns versions in a collection, with full filtering support (similar to GetModels)
func GetCollectionVersions(c *gin.Context) {
	collectionID := c.Param("id")

	// Basic validation
	var collection models.Collection
	if err := database.DB.First(&collection, collectionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	// Parse filters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	search := c.Query("search")
	baseModel := c.Query("baseModel")
	modelType := c.Query("modelType")
	onlySafe, onlyNSFW := resolveNSFWFilter(c)
	tags := c.Query("tags")
	synced := c.Query("synced") == "1"

	// Start query on Versions joined with Models
	// We need Versions that are in the collection

	versions := make([]models.Version, 0)
	q := database.DB.Model(&models.Version{}).
		Joins("JOIN collection_versions ON collection_versions.version_id = versions.id").
		Where("collection_versions.collection_id = ?", collectionID).
		Preload("ParentModel").
		Preload("Images").
		Preload("Collections")

	// Apply filters to the versions or parent model
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		q = q.Joins("JOIN models ON models.id = versions.model_id").
			Where("LOWER(models.name) LIKE ? OR LOWER(versions.name) LIKE ? OR LOWER(versions.trained_words) LIKE ?", like, like, like)
	}

	if baseModel != "" {
		q = q.Where("versions.base_model = ?", baseModel)
	}
	if modelType != "" {
		q = q.Where("versions.type = ?", modelType)
	}
	if onlySafe {
		q = q.Where("versions.nsfw = 0")
	} else if onlyNSFW {
		q = q.Where("versions.nsfw = 1")
	}
	if tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(strings.ToLower(t))
			if t != "" {
				q = q.Where("LOWER(versions.tags) LIKE ?", "%"+t+"%")
			}
		}
	}
	if synced {
		q = q.Joins("JOIN client_files ON client_files.model_version_id = versions.id").Where("client_files.status = ?", "installed")
	}

	// Count total results before pagination
	var total int64
	countQ := database.DB.Model(&models.Version{}).
		Joins("JOIN collection_versions ON collection_versions.version_id = versions.id").
		Where("collection_versions.collection_id = ?", collectionID)

	// Apply same filters for count
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		countQ = countQ.Joins("JOIN models ON models.id = versions.model_id").
			Where("LOWER(models.name) LIKE ? OR LOWER(versions.name) LIKE ? OR LOWER(versions.trained_words) LIKE ?", like, like, like)
	}
	if baseModel != "" {
		countQ = countQ.Where("versions.base_model = ?", baseModel)
	}
	if modelType != "" {
		countQ = countQ.Where("versions.type = ?", modelType)
	}
	if onlySafe {
		countQ = countQ.Where("versions.nsfw = 0")
	} else if onlyNSFW {
		countQ = countQ.Where("versions.nsfw = 1")
	}
	if tags != "" {
		for _, t := range strings.Split(tags, ",") {
			t = strings.TrimSpace(strings.ToLower(t))
			if t != "" {
				countQ = countQ.Where("LOWER(versions.tags) LIKE ?", "%"+t+"%")
			}
		}
	}
	if synced {
		countQ = countQ.Joins("JOIN client_files ON client_files.model_version_id = versions.id").Where("client_files.status = ?", "installed")
	}
	countQ.Count(&total)
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))

	q.Order("versions.id DESC").Limit(limit).Offset((page - 1) * limit).Find(&versions)

	// Populate ClientStatus for versions
	// Handlers.go has populateClientStatus for models. We need one for versions.
	// Since populateClientStatus logic is internal to handlers.go and depends on ClientFile,
	// I'll inline a simplified version here or refactor.
	// For now, I'll inline the relevant part: checking ClientFile status.

	// Fetch client files for these versions
	if len(versions) > 0 {
		versionIDs := make([]uint, len(versions))
		for i, v := range versions {
			versionIDs[i] = v.ID
		}
		var clientFiles []models.ClientFile
		database.DB.Where("model_version_id IN ?", versionIDs).Find(&clientFiles)

		statusMap := make(map[uint]string)
		for _, cf := range clientFiles {
			statusMap[cf.ModelVersionID] = cf.Status
		}

		for i := range versions {
			if s, ok := statusMap[versions[i].ID]; ok {
				versions[i].ClientStatus = s
			}
		}
	}

	c.JSON(http.StatusOK, versions)
}

// GetVersionCollections returns the collections a version belongs to
func GetVersionCollections(c *gin.Context) {
	versionID := c.Param("id")

	var version models.Version
	if err := database.DB.Preload("Collections", func(db *gorm.DB) *gorm.DB {
		return db.Order("LOWER(collections.name) ASC")
	}).First(&version, versionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}

	if version.Collections == nil {
		version.Collections = make([]models.Collection, 0)
	}
	c.JSON(http.StatusOK, version.Collections)
}

// BulkAddVersions adds all model versions matching criteria to a collection
func BulkAddVersions(c *gin.Context) {
	id := c.Param("id")
	collectionID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var input struct {
		Query              string `json:"query"`
		SearchTags         bool   `json:"searchTags"`
		SearchModelName    bool   `json:"searchModelName"`
		SearchTrainedWords bool   `json:"searchTrainedWords"`
		ExactMatch         bool   `json:"exactMatch"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[BulkAdd] Received: Query='%s', Tags=%v, Name=%v, Words=%v, Exact=%v", input.Query, input.SearchTags, input.SearchModelName, input.SearchTrainedWords, input.ExactMatch)

	queryStr := strings.ToLower(strings.TrimSpace(input.Query))
	if queryStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query cannot be empty"})
		return
	}

	if !input.SearchTags && !input.SearchModelName && !input.SearchTrainedWords {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one search criteria must be selected"})
		return
	}

	// Build the WHERE clause
	var conditions []string
	var args []interface{}

	// Helper to generate SQL for a specific column
	addCondition := func(col string, isList bool) {
		if input.ExactMatch {
			if isList {
				// Comma separated list (Tags, TrainedWords)
				// We need to match the token exactly, bounded by start/end or commas.
				// We support optional spaces after commas.
				// 1. Exact match: col = 'q'
				// 2. Start: col LIKE 'q,%' OR col LIKE 'q, %'
				// 3. End: col LIKE '%,q' OR col LIKE '%, q'
				// 4. Middle: col LIKE '%,q,%' OR col LIKE '%, q,%' OR col LIKE '%,q, %' OR col LIKE '%, q, %'

				// Simplified approach:
				// "q"
				// "q,%"
				// "%,q"
				// "%,q,%"
				// "%, q"
				// "%, q,%"

				cond := fmt.Sprintf(`(
					LOWER(%s) = ? 
					OR LOWER(%s) LIKE ? 
					OR LOWER(%s) LIKE ? 
					OR LOWER(%s) LIKE ? 
					OR LOWER(%s) LIKE ? 
					OR LOWER(%s) LIKE ?
				)`, col, col, col, col, col, col)

				conditions = append(conditions, cond)
				args = append(args,
					queryStr,
					queryStr+",%",
					"%, "+queryStr,
					"%, "+queryStr+",%",
					"%,"+queryStr,      // No space after comma (end)
					"%,"+queryStr+",%", // No space after comma (middle)
				)
			} else {
				// Space separated (Model Name)
				cond := fmt.Sprintf("(LOWER(%s) = ? OR LOWER(%s) LIKE ? OR LOWER(%s) LIKE ? OR LOWER(%s) LIKE ?)", col, col, col, col)
				conditions = append(conditions, cond)
				args = append(args, queryStr, queryStr+" %", "% "+queryStr, "% "+queryStr+" %")
			}
		} else {
			// Partial match
			conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE ?", col))
			args = append(args, "%"+queryStr+"%")
		}
	}

	if input.SearchTags {
		addCondition("versions.tags", true)
		addCondition("models.tags", true)
	}

	if input.SearchModelName {
		addCondition("models.name", false)
		addCondition("versions.name", false)
	}

	if input.SearchTrainedWords {
		addCondition("versions.trained_words", true)
	}

	whereClause := strings.Join(conditions, " OR ")

	query := `
		INSERT OR IGNORE INTO collection_versions (collection_id, version_id)
		SELECT ?, versions.id
		FROM versions
		JOIN models ON models.id = versions.model_id
		WHERE ` + whereClause

	// Prepend collectionID to args
	finalArgs := append([]interface{}{collectionID}, args...)

	log.Printf("[BulkAdd] Executing Query: %s", query)
	log.Printf("[BulkAdd] Args: %v", args) // Don't log collectionID for brevity

	result := database.DB.Exec(query, finalArgs...)
	if result.Error != nil {
		log.Printf("[BulkAdd] Error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add versions to collection"})
		return
	}

	log.Printf("[BulkAdd] Rows Affected: %d", result.RowsAffected)
	c.JSON(http.StatusOK, gin.H{"added": result.RowsAffected})
}
