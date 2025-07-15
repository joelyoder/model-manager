# Model Manager
This is a simple desktop application to download and manage models for local image generation.

## Setup
1. If you're on Windows, create a `.env` file and place it in the main directory
    - Set `CGO_ENABLED` to `1`
2. Run `go install` in \backend\
3. Run `npm i` in \frontend\
4. In the main directory, run `go run .\backend\main.go`
5. In the \frontend\ directory, run `npm run dev`

## TODO
- Data Import/Export
    - Export all model data as JSON
- Data management
    - Add a tag to multiple models at once
    - Remove a tag from multiple models at once
    - Add/remove images associated with the model version using the in-app UI
