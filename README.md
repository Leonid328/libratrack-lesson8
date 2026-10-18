# LibraTrack — Lesson 8 (PostgreSQL + GORM)

Homework: database integration, migrations, N+1 fix.

## Run

```bash
docker compose up -d
export DATABASE_URL=postgres://libratrack:libratrack@localhost:5432/libratrack?sslmode=disable
export MIGRATIONS_PATH=migrations
go run ./cmd/api
```

CI workflow: `.github/workflows/theme8-check.yml`

## Lesson 8 Report

### Migration prompt

Generate golang-migrate SQL for authors, categories, books.category_id (ON DELETE RESTRICT),
book_authors (CASCADE on book_id, RESTRICT on author_id), copies, users. Indexes on FKs.
Up/down pairs. No AutoMigrate.

### N+1 before

List without Preload: 1 SELECT books + N category + N authors (~25ms for 8 books).

### N+1 after

`Preload("Authors").Preload("Category")` — constant 3–4 SELECTs with `WHERE id IN (...)` (~8ms).
