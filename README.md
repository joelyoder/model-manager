# Model Manager
This is a desktop application to download and manage models for local image generation.

![Application preview screenshot](homepage-preview.png)

## Setup
1. Run `go install` in \backend\
2. Run `npm i` in \frontend\
3. In the main directory, run `go run .\backend\main.go`
4. In the \frontend\ directory, run `npm run dev`

## Todo
- Library management
    - Add a tag to multiple models at once
    - Remove a tag from multiple models at once
    - A way to see all current tags in the library and see all models with that tag
    - Add/remove gallery images associated with the model version
- Civitai Import
    - Add check to ensure that model files with the same name don't overrite existing models in the folder
- Settings
    - Add a setting to change the models and image folder locations
- Utilities
    - Statistics about the library size broken down by several different factors
        - Model Type
        - Base Model
        - Models vs images
    - A tool that will scan the library folder and find any models that are missing from the library
