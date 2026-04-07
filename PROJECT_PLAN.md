# Rencana Proyek — Findings API (Technical Assessment)

Dokumen ringkas hasil perencanaan: layanan **Findings API** (Go/Gin + PostgreSQL), fokus REST + auth, container, Kubernetes, CI/CD, dan IaC. **Integrasi LLM tidak termasuk** dalam scope fase ini.

---

## 1. Ringkasan

| Item | Pilihan |
|------|---------|
| Bahasa & framework | Go 1.22+, Gin |
| Database | PostgreSQL + migrasi (`golang-migrate` / `goose`) |
| Domain contoh | *Security Asset Registry* atau *Vulnerability / Findings mini* — satu entitas utama + relasi user |

---

## 2. Yang Diuji (Assessment)

- **Go REST API** — desain endpoint, validasi, error handling, pola repository/service.
- **Autentikasi** — JWT untuk alur pengguna; **HTTP Basic Auth** sederhana dari **environment variable** untuk rute internal/operasional.
- **Docker** — image multi-stage, user non-root, healthcheck.
- **Kubernetes** — Deployment, Service, probes, Secret/ConfigMap, **HPA** untuk skenario traffic tinggi, connection pooling DB.
- **CI/CD** — pipeline (mis. GitHub Actions / GitLab CI): test, lint (opsional), build image, push registry, deploy.
- **IaC** — Terraform atau Pulumi minimal (satu lingkungan atau modul kecil: VPC/cluster/registry/VM).

---

## 3. Endpoint Contoh

### 3.1 Publik (tanpa auth)

| Metode | Path | Keterangan |
|--------|------|------------|
| `GET` | `/health` | Liveness |
| `GET` | `/ready` | Readiness (cek koneksi DB) |
| `GET` | `/api/v1/public/info` | Info service/versi (tanpa data sensitif) |

### 3.2 Auth pengguna (JWT)

| Metode | Path | Keterangan |
|--------|------|------------|
| `POST` | `/api/v1/auth/register` | Registrasi user (opsional sesuai kebutuhan demo) |
| `POST` | `/api/v1/auth/login` | Mengembalikan JWT |
| `GET` | `/api/v1/me` | Profil dari token |
| `GET` | `/api/v1/findings` | Contoh list resource |
| `POST` | `/api/v1/findings` | Contoh create resource |

*(Nama resource `findings` bisa diganti `assets` jika kebutuhan assessment berubah.)*

### 3.3 Basic Auth dari env (internal)

Kredensial hanya dari environment, **tidak** di-hardcode di repo.

| Variabel contoh | Fungsi |
|-----------------|--------|
| `BASIC_AUTH_USER` | Username Basic Auth |
| `BASIC_AUTH_PASSWORD` | Password (di produksi dari Secret K8s / CI) |

| Metode | Path | Keterangan |
|--------|------|------------|
| `GET` | `/internal/metrics` | Uji Basic Auth berhasil/gagal |
| `POST` | `/internal/cache/purge` | Simulasi aksi admin (body boleh kosong) |

**Perilaku dev:** dokumentasikan apakah env kosong menonaktifkan Basic Auth atau menolak semua akses ke `/internal/*`.

**Tes cepat**

```bash
curl http://localhost:8080/health
curl -u "$BASIC_AUTH_USER:$BASIC_AUTH_PASSWORD" http://localhost:8080/internal/metrics
```

---

## 4. Arsitektur Target

```
[Client] → [Ingress / LB] → [Pod Gin API] → [PostgreSQL]
```

**Lokal:** `docker-compose` — layanan API + Postgres (Redis opsional untuk rate limit/session).

---

## 5. Fase Kerja

| Fase | Isi | Estimasi kasar |
|------|-----|----------------|
| 0 | Struktur repo (`cmd/`, `internal/`), konfigurasi env | 0,5–1 hari |
| 1 | Skema DB, migrasi, CRUD + JWT + middleware Basic Auth untuk `/internal/*` | 2–4 hari |
| 2 | Dockerfile multi-stage, `docker-compose` dev | 0,5–1 hari |
| 3 | Manifest K8s (atau Helm minimal): probes, resources, HPA, Secret | 1–2 hari |
| 4 | CI/CD + deploy ke cloud / cluster | 1–2 hari |
| 5 | IaC (folder `terraform/` atau setara) | ~1 hari (bisa paralel dengan fase 4) |

**Total tanpa LLM:** sekitar **5–8 hari** kerja fokus (atau ~1–1,5 minggu part-time).

---

## 6. Deliverable Pengumpulan

1. README: cara jalan lokal, daftar env, diagram arsitektur (ASCII/cukup).
2. Source + migrasi PostgreSQL.
3. `Dockerfile` + `docker-compose.yml`.
4. Manifest Kubernetes (+ Helm jika dipakai).
5. Workflow CI/CD (YAML).
6. IaC sesuai pilihan stack.

---

## 7. Checklist Sebelum Submit

- [ ] Tidak ada secret atau password di git; hanya placeholder + dokumentasi.
- [ ] `/health` dan `/ready` berfungsi.
- [ ] Publik vs JWT vs Basic Auth dapat diuji dengan jelas.
- [ ] Migrasi DB dapat dijalankan (dokumentasi atau init job).
- [ ] Image container non-root dan sekecil mungkin (multi-stage).
- [ ] HPA + probes terdokumentasi untuk skenario traffic tinggi.
- [ ] Pipeline hijau minimal sekali end-to-end (build + test).

---

## 8. Catatan Scope (LLM)

Integrasi LLM lokal / GPU **sengaja tidak dimasukkan** dalam rencana implementasi saat ini. Bisa ditambahkan sebagai fase terpisah jika requirement berubah.

---

*Dokumen ini merangkum keputusan perencanaan proyek assessment; sesuaikan nama resource, cloud provider, dan tooling CI dengan akses tim Anda.*
