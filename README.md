# Heard - Go API Skeleton

Quick start

1. Set up Postgres and create the `heard` database (or update `DATABASE_URL`).
2. Copy the `.env.example` file to `.env` and update the environment variables as needed:
   ```bash
   cp .env.example .env
   ```
3. Apply `database_setup.sql` to create tables and seed data.
4. Run:

```bash
go mod tidy
go run ./cmd/api
```

Endpoints

- `GET /companies` - list companies
- `GET /companies?id=1` - get company
- `POST /companies` - create company (JSON body)
- `PUT /companies?id=1` - update company (JSON body)
- `DELETE /companies?id=1` - delete company

Same pattern for `/users`, `/posts`, `/comments`.