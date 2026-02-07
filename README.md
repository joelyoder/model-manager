# Model Manager

This is a web application to download and manage models for local image generation that can be run as a standalone application or server with desktop client.

![Application preview screenshot](homepage-preview.png)

## Overview

Model Manager combines a Go backend and a Vue 3 single-page application to provide a cohesive experience:

- The Go backend (Gin) exposes REST APIs under `/api`, synchronizes metadata from Civitai, manages downloads, and persists data in a SQLite database.
- The Vue frontend, built with Vite, consumes those APIs to present search, filtering, and management workflows. During development it runs with Vite's dev server, while in production the backend serves the built SPA from `frontend/dist` alongside static assets.

## Features

- Powerful model search and filtering across base models, types, and tags.
- Version download management with progress tracking and orphan detection for unused files.
- Gallery management, including uploading supplemental images, setting main previews, and pruning outdated assets.
- Metadata import/export utilities and library statistics to keep local collections organized.
- Integration with the Civitai API for synchronizing models and refreshing version information.
- **Smart Remote Management**: Interactive controls for seamless downloading, syncing, and removal of models on remote clients with real-time status updates.
- **Image Archiving**: Automated tools to archive external descriptive images to local storage, ensuring content remains available even if the source disappears.
- **Custom Collections**: Effortlessly organize models into custom groups with support for bulk adding versions by tag and quick-access management.

## Prerequisites

Install or provision the following before working with the project:

- **Go**: version 1.24.4 or newer (matches the module's `go` directive).
- **Node.js**: version 18.17 or newer.
- **Package manager**: npm 9+ or Yarn 1.22+ (examples below use npm).
- **External services**: a Civitai account and API token for features that synchronize or download content from Civitai.

## Setup

1. Install Go dependencies from the project root:
   ```sh
   go mod download
   ```
2. Install frontend dependencies (run this from the repository root):
   ```sh
   npm --prefix frontend install
   ```
3. Run the Go backend from the repository root (it will load variables from a `.env` file if present):
   ```sh
   go run ./backend
   ```
4. In a separate terminal, start the Vue development server:
   ```sh
   npm --prefix frontend run dev
   ```
5. Open the URL printed by Vite (typically http://localhost:5173). The dev server proxies `/api` and `/images` to `http://localhost:8080`; if you run the Go backend on a different port, update the proxy targets in `frontend/vite.config.js` to match.

## Desktop Client

The project includes a Python-based desktop client (located in the `client` directory) that runs in the system tray. This client allows the backend to remotely manage files on another machine (e.g., your primary workstation) via WebSockets.

### Prerequisites

- **Python**: version 3.8 or newer.
- **Python Packages**: install required dependencies:
  ```sh
  pip install -r client/requirements.txt
  ```

### Configuration

1. Copy `client/config.json.example` to `client/config.json`.
2. Edit `client/config.json` with your settings:
   - `server_url`: The WebSocket URL of your backend (e.g., `ws://localhost:8080/ws`).
   - `api_key`: A secret key used for authentication. This **must match** the `CLIENT_SECRET` environment variable set on the backend.
   - `root_path`: The local directory where models will be managed.
   - `client_id`: A unique name for this client.

### Running & Building

To run the client during development:
```sh
python client/client_tray.py
```

To build a standalone executable for Windows:
```sh
cd client
build_client.bat
```
The executable will be generated in `client/dist/ModelManagerClient.exe`.

## Environment Variables

Set variables in your shell or a `.env` file located at the repository root:

| Variable | Description | Default |
| --- | --- | --- |
| `PORT` | Port the Go HTTP server listens on. | `8080` |
| `MODELS_DB_PATH` | Filesystem path to the SQLite database used by GORM. Relative paths resolve from the server's working directory. | `backend/models.db` |
| `CIVIT_API_KEY` | Personal access token for authenticating requests to the Civitai API (required for syncing and downloads). | _unset_ |
| `CLIENT_SECRET` | Secret key for authenticating the desktop client WebSocket connection (must match `api_key` in client config). | _unset_ |

Create a `.env` file in the repository root to persist these variables locally. Generate a Civitai token from <https://civitai.com/user/account/api> and assign it to `CIVIT_API_KEY` to enable synchronization features.

The backend serves user-managed assets from `./backend/images` at `/images` and downloaded model files from `./backend/downloads` at `/downloads`. It creates those directories if they do not exist, but the process must have permission to create and write to them.

## Tests

### Backend
Run all Go tests from the repository root:

```sh
go test ./...
```

### Frontend
Unit tests are executed with [Vitest](https://vitest.dev/):

```sh
npm --prefix frontend run test
```

End-to-end tests use [Playwright](https://playwright.dev/):

```sh
npm --prefix frontend run test:e2e
```

Before the first Playwright run, install the browser binaries:

```sh
npm --prefix frontend exec playwright install
```

## Deployment

Build the production assets and compile a standalone backend binary:

```sh
# From the repository root
npm --prefix frontend install
npm --prefix frontend run build

go build -o model-manager ./backend
```

The compiled binary expects the `frontend/dist` directory (created by `npm run build`) and the `backend/images` and `backend/downloads` folders to exist relative to its working directory. To run the production build:

```sh
./model-manager
```

Set environment variables (such as `PORT`, `MODELS_DB_PATH`, or `CIVIT_API_KEY`) before launching if you need non-default values.

## Gallery Management

Use the model detail page to upload additional images or remove existing gallery images from a version. The uploaded image will be scanned for embedded metadata and displayed alongside the image.

### API Endpoints
- `POST /api/versions/:id/images` – upload a gallery image. Form field `file` should contain the image. Returns the created `VersionImage` record.
- `DELETE /api/versions/:id/images/:imgId` – remove a gallery image.

## Known Issues
- Model images are not delivered to the desktop client

## Todo
- Library management
    - Add a tag to multiple models at once
    - Remove a tag from multiple models at once
    - A way to see all current tags in the library and see all models with that tag
