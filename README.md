# Model Manager
This is a simple desktop application to download and manage models for local image generation.

## Setup
1. Create a `.env` file and place it in the main directory
    - Set `CIVIT_API_KEY` to your API key
    - If you're on Windows, set `CGO_ENABLED` to `1`
2. Run `go install` in \backend\
3. Run `npm i` in \frontend\
4. In the main directory, run `go run .\backend\main.go`

## TODO
- Data Import/Export
    - Import models as JSON
    - Export all model data as JSON
- Data management
    - Add a tag to multiple models at once
    - Remove a tag from multiple models at once
    - Add/remove images associated with the model version
    - Select a different cover image for the model version
- Settings
    - Move Civitai API key to setting stored in the database
    - Add setting for the models folder location
