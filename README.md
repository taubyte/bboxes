# bboxes

Build, test, and publish Taubyte WASI container images. Uses a directory build context and `docker build` (no tarballs, no `go generate`).

## Requirements

- **Go 1.25** (to build the tool)
- **Docker** (for building and running images)

## Usage

```text
bboxes <build|test|publish> <language> <version>
```

| Subcommand | Description |
|------------|-------------|
| `build`    | Build the image only; tag with the given version. |
| `test`     | Build the image, then run a smoke test (`docker run --rm <image> true`). |
| `publish`  | Build the image and push it to the registry. |

| Argument  | Values | Notes |
|-----------|--------|--------|
| `language` | `go`, `rs`, `as` | Target runtime. |
| `version`  | e.g. `v0.1.0`    | Image tag. |

### Image mapping

| language | Image (org/repo)           |
|----------|----------------------------|
| go       | taubyte/go-wasi            |
| rs       | taubyte/rust-wasi          |
| as       | taubyte/assembly-script-wasi |

Go builds **reactor-style** WASM modules: no `_start`/main, only exported functions the host calls (`//export`).

### Examples

```bash
# Build Go image and tag as v0.1.0
bboxes build go v0.1.0

# Publish Go image
bboxes publish go v0.2.0

# Build Rust image and run smoke test
bboxes test rs v0.1.0
```

### Publishing (registry login)

For `publish`, the tool runs `docker login`. Either:

- Run `docker login` yourself first, or
- Set **DOCKER_USER** and **DOCKER_TOKEN** (or **DOCKER_USERNAME** / **DOCKER_PASSWORD**).

## Build the tool

From the repo root:

```bash
go build -o bboxes .
```

Run directly:

```bash
go run . build go v0.1.0
```

## How it works

1. **Context:** For the chosen language, the tool copies `containers/<lang-dir>` and `containers/common` into a temporary directory (lang first, then common), and writes the language’s Dockerfile as `Dockerfile`.
2. **Build:** It runs `docker build` with that directory as the build context (no tarballs).
3. **Publish:** If you use `publish`, it runs `docker login` (if needed) and then `docker push` for the versioned image.
