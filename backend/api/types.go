package api

type CivitModel struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Type          string           `json:"type"`
	Description   string           `json:"description"`
	Nsfw          bool             `json:"nsfw"`
	Tags          []string         `json:"tags"`
	Mode          string           `json:"mode"`
	ModelVersions []VersionSummary `json:"modelVersions"`
	Created       string           `json:"createdAt"`
	Updated       string           `json:"updatedAt"`
}

type VersionSummary struct {
	ID                   int          `json:"id"`
	Name                 string       `json:"name"`
	BaseModel            string       `json:"baseModel"`
	EarlyAccessTimeFrame int          `json:"earlyAccessTimeFrame"`
	TrainedWords         []string     `json:"trainedWords"`
	Created              string       `json:"createdAt"`
	Updated              string       `json:"updatedAt"`
	Files                []ModelFile  `json:"files"`
	Images               []ModelImage `json:"images"`
}

type VersionResponse struct {
	ID                   int          `json:"id"`
	ModelID              int          `json:"modelId"`
	Name                 string       `json:"name"`
	BaseModel            string       `json:"baseModel"`
	Created              string       `json:"createdAt"`
	Updated              string       `json:"updatedAt"`
	EarlyAccessTimeFrame int          `json:"earlyAccessTimeFrame"`
	TrainedWords         []string     `json:"trainedWords"`
	ModelFiles           []ModelFile  `json:"files"`
	Images               []ModelImage `json:"images"`
}

type ModelFile struct {
	Name        string  `json:"name"`
	DownloadURL string  `json:"downloadUrl"`
	SizeKB      float64 `json:"sizeKB"`
	Hashes      struct {
		SHA256 string `json:"SHA256"`
	} `json:"hashes"`
}

type ModelImage struct {
	URL      string                 `json:"url"`
	URLSmall string                 `json:"urlSmall"`
	Width    int                    `json:"width"`
	Height   int                    `json:"height"`
	Hash     string                 `json:"hash"`
	Meta     map[string]interface{} `json:"meta"`
}

// VersionInfo represents a simplified view of a model version returned to the frontend.
// It contains the basic fields required for display and selection when downloading
// a specific model version.
type VersionInfo struct {
	ID                   int      `json:"id"`
	ModelID              int      `json:"modelId"`
	Name                 string   `json:"name"`
	BaseModel            string   `json:"baseModel"`
	SizeKB               float64  `json:"sizeKB"`
	TrainedWords         []string `json:"trainedWords"`
	EarlyAccessTimeFrame int      `json:"earlyAccessTimeFrame"`
	SHA256               string   `json:"sha256"`
	Created              string   `json:"createdAt"`
	Updated              string   `json:"updatedAt"`
}

// ImportRecord represents a single entry in the JSON import file.
// Only a subset of fields are mapped into the database.
type ImportRecord struct {
	Name            string   `json:"name"`
	BaseModel       string   `json:"base_model"`
	ModelType       string   `json:"model_type"`
	DownloadURL     string   `json:"download_url"`
	URL             string   `json:"url"`
	PreviewURL      string   `json:"preview_url"`
	Description     string   `json:"description"`
	PositivePrompts string   `json:"positive_prompts"`
	SHA256Hash      string   `json:"sha256_hash"`
	CreatedAt       float64  `json:"created_at"`
	Groups          []string `json:"groups"`
	Location        string   `json:"location"`
}
