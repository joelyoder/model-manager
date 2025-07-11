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
	ID int `json:"id"`
}

type VersionResponse struct {
	ID                   int          `json:"id"`
	ModelID              int          `json:"modelId"`
	Name                 string       `json:"name"`
	BaseModel            string       `json:"baseModel"`
	Created              string       `json:"createdAt"`
	EarlyAccessTimeFrame int          `json:"earlyAccessTimeFrame"`
	TrainedWords         []string     `json:"trainedWords"`
	ModelFiles           []ModelFile  `json:"files"`
	Images               []ModelImage `json:"images"`
}

type ModelFile struct {
	Name        string  `json:"name"`
	DownloadURL string  `json:"downloadUrl"`
	SizeKB      float64 `json:"sizeKB"`
}

type ModelImage struct {
	URL      string `json:"url"`
	URLSmall string `json:"urlSmall"`
}
