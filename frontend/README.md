# Model Manager Frontend

## Purpose
The Model Manager frontend is a Vue 3 + Vite single-page application that lets you browse, organize, and enrich a library of generative AI models that are tracked by the Go backend. It surfaces the models stored in the Model Manager database, provides tools for importing releases from Civitai, and gives curators quick access to metadata, gallery images, download status, and cleanup utilities.

### Feature highlights
- **Library browsing & filtering:** Search models by name or tags, filter by category, base model, and type, hide NSFW entries, and jump directly to detailed views.
- **Version management:** Inspect an individual model version, edit metadata, upload gallery images, trigger metadata refreshes, and remove unwanted versions or their files.
- **Civitai integration:** Save a personal API key, fetch version metadata by pasting a Civitai model URL, and download the associated files directly into your library.
- **Utilities dashboard:** Review collection statistics, import or export data from JSON, and find orphaned or duplicate files for cleanup.

## Project structure
The frontend lives entirely in the `frontend/` directory. The most important files and folders are:

### Source code (`src/`)
- `main.js` – Boots the Vue app, wires up the router, and loads global styles.
- `App.vue` – Hosts the top-level layout, shared navigation, toast container, confirmation modals, and the back-to-top control.
- `router.js` – Defines the routes that point to the model library, model detail, settings, and utilities views.
- `components/`
  - `ModelList.vue` – The landing view that lists models, applies filters, loads Civitai versions, and starts downloads.
  - `ModelDetail.vue` – Displays a single model version with editable metadata, gallery management, and backend sync actions.
  - `AppSettings.vue` – Simple settings panel for persisting the Civitai API key through the backend.
  - `UtilitiesPage.vue` – Statistics, import/export helpers, and cleanup tools (orphaned files and duplicate paths).
  - `BackToTop.vue` – Floating button that reveals itself after scrolling to quickly return to the top of long lists.
- `utils/` – Small helper modules (`debounce.js`, `ui.js`) and their Vitest suites for reusable UI behaviours.
- `index.css` – Global stylesheet that augments Bootstrap defaults.

### Tooling & configuration
- `package.json` – npm scripts for local development, testing, linting, and building the production bundle.
- `vite.config.js` – Vite configuration shared across dev and production builds.
- `playwright.config.ts` – Playwright end-to-end test harness; it builds the frontend, launches the Go backend with a temporary SQLite database, and runs UI flows against `http://localhost:8080`.
- `tests/e2e.spec.ts` – High-level journey test that exercises critical user paths with Playwright.
- `public/` – Static assets copied verbatim into the built output.

## Development workflow
1. Install dependencies
   ```bash
   npm install
   ```
2. Start the Vite dev server with hot module reloading
   ```bash
   npm run dev
   ```
3. Build the production bundle
   ```bash
   npm run build
   ```
4. Lint and auto-fix with ESLint
   ```bash
   npm run lint
   ```

## Testing
### Unit tests (Vitest)
Run the headless unit tests with:
```bash
npm test
```
Vitest executes the suites colocated with the utilities under `src/utils/*.spec.js` using the jsdom environment. Add additional specs next to the components or helpers they cover to grow the suite.

### End-to-end tests (Playwright)
Run the full-stack UI regression test with:
```bash
npm run test:e2e
```
The Playwright configuration automatically builds the frontend, starts the Go backend (`go run backend/main.go`) against a temporary SQLite database, and drives the Chromium browser through the scenarios defined in `tests/e2e.spec.ts`. Ensure Go is available locally before running this command.
