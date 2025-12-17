# darulabror-api

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](#)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](#)
[![Swagger](https://img.shields.io/badge/Swagger-API%20Docs-85EA2D?logo=swagger&logoColor=000)](https://darulabror-717070183986.asia-southeast2.run.app/swagger/index.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Darul Abror](https://img.shields.io/badge/Darul%20Abror-API-111827)](#)

Darul Abror API service (public + admin area) built with Echo, GORM, PostgreSQL, and optional Google Cloud Storage integration for public article media.

## API Documentation (Swagger)

- Swagger UI: https://darulabror-717070183986.asia-southeast2.run.app/swagger/index.html
- OpenAPI JSON: https://darulabror-717070183986.asia-southeast2.run.app/swagger/doc.json
- OpenAPI YAML: https://darulabror-717070183986.asia-southeast2.run.app/swagger/doc.yaml

## Response format (convention)

Most JSON responses follow this envelope (see `internal/utils/response.go`):

Success:
```json
{
  "status": "success",
  "message": "OK",
  "data": {}
}
```

Error:
```json
{
  "status": "error",
  "message": "something went wrong"
}
```

Notes:
- Some endpoints intentionally return **No Content** (`201/204` with empty body) because handlers use `c.NoContent(...)`.

## Authentication (Admin)

1) Login to get JWT:
- `POST /admin/login`

2) Use token for admin endpoints:
- Header: `Authorization: Bearer <token>`

Role rules (see `api/routes/routes.go`):
- `/admin/*` requires role: `admin` or `superadmin`
- `/admin/admins*` requires role: `superadmin`

## Pagination

List endpoints support:
- `page` (default 1)
- `limit` (default 10, max 100)

Response uses:
```json
{
  "items": [],
  "meta": { "page": 1, "limit": 10, "total": 123 }
}
```

---

## Articles: `content` is flexible JSON (frontend-defined)

In the database/model, `Article.content` is stored as **JSONB** (`gorm.io/datatypes.JSON`).
That means the backend does **not** enforce a fixed schema for article body content.

**Frontend decides the shape**, e.g. a block editor style:
```json
{
  "blocks": [
    { "type": "heading", "level": 2, "text": "Judul" },
    { "type": "paragraph", "text": "Teks panjang..." },
    { "type": "image", "url": "https://example.com/image.jpg", "caption": "..." },
    { "type": "video", "url": "https://example.com/video.mp4" }
  ]
}
```

Important:
- Binary file (image/video) **is not stored inside `content`**.
- `content` should store **URLs** (or any reference you decide) that the frontend can render.
- This project currently supports **uploading only the article header image** via the create/update endpoints (see below).

### `photo_header`
`photo_header` is intended for the article card/cover image (thumbnail/banner).
It can be:
- a URL string you provide (`photo_header` form field), OR
- automatically set by uploading a file (`photo_header_file`)

---

## Admin Articles (multipart/form-data)

### POST /admin/articles
Create article using `multipart/form-data`.

Fields:
- `title` (string, required)
- `author` (string, required)
- `status` (string, optional: `draft|published`)
- `content` (string, required) → **must be valid JSON string**
- `photo_header` (string, optional) → header URL
- `photo_header_file` (file, optional) → upload header file; overrides `photo_header`

Response:
- `201 Created` (no body)

Requirements for `photo_header_file` upload:
- Set `PUBLIC_BUCKET` env var in server runtime (GCS). If not configured, upload will fail.

### PUT /admin/articles/:id
Same fields/behavior as create (multipart). Response:
- `200 OK` (no body)

---

## Health & Swagger

### GET /healthz
Response:
- `200 OK` body: `ok`

### GET /swagger/index.html
Swagger UI

---

## Public Endpoints

### GET /articles
Query:
- `page` (optional)
- `limit` (optional)

Response `200`:
```json
{
  "status": "success",
  "message": "articles fetched",
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Example",
        "photo_header": "https://example.com/header.jpg",
        "content": {},
        "author": "Admin",
        "status": "published",
        "created_at": 1734567890,
        "updated_at": 1734567890
      }
    ],
    "meta": { "page": 1, "limit": 10, "total": 1 }
  }
}
```

### GET /articles/:id
Response:
- `200` success envelope with `data` = article
- `400` if id invalid
- `404` if not found / not published

### POST /registrations
Request body: `dto.RegistrationDTO` (see `internal/dto/registration_dto.go`)
Response:
- `201 Created` (no body)

### POST /contacts
Request body:
```json
{
  "email": "user@example.com",
  "subject": "Question",
  "message": "Hello..."
}
```
Response:
- `201 Created` (no body)

---

## License
MIT. See [LICENSE](LICENSE).