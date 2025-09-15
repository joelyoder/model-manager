# Model Manager
This is a desktop application to download and manage models for local image generation.

![Application preview screenshot](homepage-preview.png)

## Setup
1. Run `go install` in \backend\
2. Run `npm i` in \frontend\
3. In the main directory, run `go run .\backend\main.go`
4. In the \frontend\ directory, run `npm run dev`

## Tests
### Backend
Run all Go tests from the backend directory:

```sh
go test ./...
```

### Frontend
Unit tests are executed with [Vitest](https://vitest.dev/):

```sh
npm test
```

End-to-end tests use [Playwright](https://playwright.dev/):

```sh
npm run test:e2e
```

## Gallery Management
Use the model detail page to upload additional images or remove existing gallery
images from a version. The uploaded image will be scanned for embedded metadata
and displayed alongside the image.

### API Endpoints
- `POST /api/versions/:id/images` – upload a gallery image. Form field `file`
  should contain the image. Returns the created `VersionImage` record.
- `DELETE /api/versions/:id/images/:imgId` – remove a gallery image.

## Todo
- Library management
    - Add a tag to multiple models at once
    - Remove a tag from multiple models at once
    - A way to see all current tags in the library and see all models with that tag
- Settings
    - Add a setting to change the models and image folder locations
