# Technical Assessment — Findings API (Go + Gin + PostgreSQL)

API REST tahap awal: health/ready, info publik, JWT (register/login/findings), dan rute internal dengan HTTP Basic Auth dari environment.

## Prasyarat

- Go **1.22+**
- Docker & Docker Compose (opsional, untuk Postgres + API sekaligus)

## Struktur folder

```
cmd/server/          # entrypoint
internal/config/     # env & validasi konfigurasi
internal/db/         # koneksi PostgreSQL + migrasi
internal/handlers/   # HTTP handlers
internal/middleware/ # JWT & Basic Auth
internal/repository/ # akses data
internal/router/     # registrasi rute Gin
migrations/          # SQL golang-migrate
infra/               # IaC Pulumi (AWS ECR contoh, stack dev)
k8s/                 # Manifest Kubernetes (Deploy, SVC, CM, Secret contoh, HPA, Postgres)
```

## Menjalankan lokal (Postgres di Docker, API dengan `go run`)

1. Salin environment:

   ```bash
   cp .env.example .env
   ```

   Sesuaikan `JWT_SECRET` (minimal **32 karakter**) dan kredensial Basic Auth bila perlu.

2. Jalankan hanya database:

   ```bash
   docker compose up -d postgres
   ```

3. Dari root repo:

   ```bash
   go run ./cmd/server
   ```

   Server mendengarkan `HTTP_ADDR` (default `:8080`). Migrasi dijalankan otomatis saat startup. Dokumentasi web: `http://localhost:8080/` → `/docs` (beranda), `/docs/guide` (panduan setup), `/docs/api` (referensi API), `/docs/api-overview.xml` (XML), `/docs/postman-collection.json` (koleksi Postman v2.1 untuk import).

## Menjalankan penuh dengan Docker Compose

```bash
docker compose up --build
```

API: `http://localhost:8080` — `JWT_SECRET` di `docker-compose.yml` hanya untuk demo; ganti di lingkungan nyata.

## Variabel environment

| Variabel | Wajib | Keterangan |
|----------|--------|------------|
| `DATABASE_URL` | ya | DSN PostgreSQL |
| `JWT_SECRET` | ya | Minimal 32 karakter |
| `HTTP_ADDR` | tidak | Default `:8080` |
| `BASIC_AUTH_USER` | untuk `/internal/*` | Jika kosong → `503` pada rute internal |
| `BASIC_AUTH_PASSWORD` | untuk `/internal/*` | Sama seperti di atas |
| `SERVICE_NAME` | tidak | Default `findings-api` |
| `SERVICE_VERSION` | tidak | Default `0.1.0` |

File `.env` dibaca otomatis jika ada (`godotenv`).

## Troubleshooting

### `pq: role "postgres" does not exist`

Artinya **user di `DATABASE_URL`** (bagian sebelum `@`) tidak ada di server PostgreSQL yang sedang dipakai.

- **Postgres dari Docker / Compose** di repo ini memang memakai user `postgres`. Pastikan container jalan dan `DATABASE_URL` memakai `postgres:postgres@...` seperti di `.env.example` baris terakhir.
- **Postgres dari Homebrew di macOS** sering **tidak** punya role `postgres`; superuser default biasanya **nama user macOS-mu**. Ganti `DATABASE_URL`, misalnya:

  `postgres://namaloginmac@localhost:5432/findings?sslmode=disable`

  Buat database bila belum ada:

  ```bash
  createdb findings
  ```

  Atau buat role `postgres` sekali saja (opsional):

  ```bash
  createuser -s postgres
  ```

### `docker compose` / flag `--build` tidak dikenali

Lihat penjelasan di diskusi: pasang Docker Desktop / plugin Compose, atau pakai perintah `docker-compose` (dengan tanda hubung).

## Uji cepat

```bash
curl -s http://localhost:8080/health
curl -s http://localhost:8080/ready
curl -s http://localhost:8080/api/v1/public/info

curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"password123"}'

TOKEN="$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"password123"}' | jq -r .token)"

curl -s http://localhost:8080/api/v1/me -H "Authorization: Bearer $TOKEN"

curl -s -u "admin:admin" http://localhost:8080/internal/metrics
```

## GitHub Actions (CI tanpa production)

Workflow [`.github/workflows/ci.yml`](.github/workflows/ci.yml) hanya memverifikasi kode di runner GitHub:

- **Go:** `go mod verify`, `go vet`, `build` server, `go test -race ./...`
- **Docker:** `docker build` lokal di runner (**tanpa** push ke registry dan **tanpa** deploy)

Tidak perlu secret atau environment production. Setelah repo di-push ke GitHub, buka tab **Actions** untuk melihat status hijau. Job **Pulumi IaC (compile)** memastikan program di `infra/` dapat dikompilasi tanpa menjalankan `pulumi up`.

Detail penggunaan Pulumi (preview/deploy, stack `dev` vs nanti staging/prod): [`infra/README.md`](infra/README.md).

## Kubernetes

Manifest siap pakai ada di [`k8s/`](k8s/): namespace, Postgres + PVC, ConfigMap, Secret contoh, Deployment API (probes, 2 replika), Service, HPA. Ikuti [`k8s/README.md`](k8s/README.md) untuk `kubectl apply` dan pemuatan image lokal (kind/minikube).

## Rencana detail

Lihat [`PROJECT_PLAN.md`](./PROJECT_PLAN.md).
