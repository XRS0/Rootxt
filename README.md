# System Journal Portfolio

This repo contains a content-first portfolio for a backend engineer. It is a technical journal, a system interface, and a manifesto for operational clarity.

## 1. Overall architecture explanation
- The project is split into a Go backend API and a Vite + React frontend.
- Articles are the primary asset. They exist as Markdown drafts and are rendered to HTML only when published.
- The backend follows DDD: domain is pure, application orchestrates, infrastructure implements, HTTP is thin.
- The frontend renders the public journal and a minimal admin panel for writing and publishing.

## 2. Backend structure and responsibilities
```
cmd/server
  main.go                 # wiring, config, server startup

internal
  domain
    article               # entity + rules + repository interface
    user                  # entity + rules + repository interface

  application
    article               # use cases (create/update/publish/list/get/delete)
    user                  # login + bootstrap admin

  infrastructure
    auth                  # JWT + bcrypt hasher
    db                    # bun models + repositories + migrations
    http                  # router assembly
    markdown              # Markdown -> HTML rendering

  interfaces
    http
      handlers            # thin handlers, DTOs
      middleware          # JSON + CORS + auth
```

## 3. Database schema and migrations
Migrations live in `internal/infrastructure/db/migrations` using bun's migrator.

Tables:
- `users`
  - `id` BIGSERIAL PK
  - `email` TEXT UNIQUE NOT NULL
  - `password_hash` TEXT NOT NULL
  - `created_at` TIMESTAMPTZ
  - `updated_at` TIMESTAMPTZ

- `articles`
  - `id` BIGSERIAL PK
  - `title` TEXT NOT NULL
  - `slug` TEXT UNIQUE NOT NULL
  - `markdown_content` TEXT NOT NULL
  - `html_content` TEXT
  - `status` TEXT NOT NULL (`draft` or `published`)
  - `created_at` TIMESTAMPTZ
  - `updated_at` TIMESTAMPTZ

Slug uniqueness is enforced at the database level. The application layer also ensures uniqueness on create/update.

## 4. Domain and application layer code
- Domain `article.Article` is the source of truth for invariants. It owns slug creation, state transitions, and validation.
- Application use cases orchestrate repository calls and timestamps.
- Markdown rendering is injected via an interface and only invoked in the publish use case.
- Updating a published article moves it back to draft and clears HTML to force an explicit republish.

## 5. HTTP API
Base path: `/api`

Public:
- `GET /api/articles`
  - Returns published articles (title, slug, excerpt, timestamps)
- `GET /api/articles/{slug}`
  - Returns a published article (title, slug, html, timestamps)

Admin:
- `POST /api/admin/login`
  - Body: `{ "email": "", "password": "" }`
  - Returns: `{ "token": "..." }`

Authenticated (Bearer token):
- `GET /api/admin/articles`
- `POST /api/admin/articles` (title + markdown)
- `PUT /api/admin/articles/{id}`
- `DELETE /api/admin/articles/{id}`
- `POST /api/admin/articles/{id}/publish`
- `POST /api/admin/articles/{id}/unpublish`

## 6. Frontend routing and layout
Routes:
- `/` homepage (hero + beliefs + recent articles)
- `/articles` list
- `/articles/:slug` article page
- `/about`
- `/contact`
- `/admin` admin panel

Layout favors typography and obvious structure: oversized headings, strong dividers, monospaced metadata, no decoration.

## 7. Brutalist CSS example
From `web/src/styles/base.css`:
```
:root {
  --bg: #ffffff;
  --fg: #111111;
  --accent: #d10000;
  --border: 2px solid #111111;
}

.site-header {
  border-bottom: var(--border);
  text-transform: uppercase;
}

.article-list li {
  padding: 1rem 0;
  border-bottom: 1px solid #111111;
}
```

## 8. Admin panel structure
- Table-based article list (title + status)
- Plain text buttons for actions
- Markdown textarea editor with live preview
- Publish / unpublish workflow with explicit actions

## 9. Minimal deployment instructions
Backend:
1. Set environment variables:
   - `DATABASE_URL` (or `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`)
   - `JWT_SECRET`
   - `ADMIN_EMAIL` + `ADMIN_PASSWORD` (optional bootstrap)
   - `CORS_ORIGIN` (optional for frontend dev)
2. Run the server:
   - `go run ./cmd/server`

Frontend:
1. `cd web`
2. `npm install`
3. `npm run dev`
4. Set `VITE_API_BASE` if the API is not on the same origin.

Note: `go mod tidy` failed in this sandbox due to Go build cache permissions. Run it locally to generate `go.sum` and lock versions.

Docker:
1. `docker compose up --build`
2. Frontend: http://localhost:5173
3. Backend API: http://localhost:8080
4. Postgres: localhost:5432 (user/pass: postgres/postgres)
