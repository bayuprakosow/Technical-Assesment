# Panduan Presentasi — Findings API (Technical Assessment)

Dokumen ini membantu kamu **memaparkan proyek secara jujur** (keterbatasan lokal), **menautkan ke rubrik penilaian**, dan **menggambarkan arah ke depan** tanpa mengada-ada infrastruktur yang belum ada.

---

## 1. Pembuka (30–45 detik)

- Saya membangun **Findings API**: REST API dengan **Go + Gin + PostgreSQL**, fokus pola yang dipakai di lingkungan produksi (auth, container, orkestrasi, otomasi, IaC).
- Scope mengikuti brief assessment: **REST + auth**, **Docker**, **Kubernetes**, **CI/CD**, **IaC**; **LLM** sengaja tidak dimasukkan sesuai rencana.
- Saya akan jelaskan **apa yang sudah jalan di mesin lokal / CI**, **apa yang direpresentasikan sebagai artefak** (manifest, Pulumi), dan **rencana penguatan** bila ada akses cloud/cluster sungguhan.

---

## 2. Keterbatasan lokal (transparansi ~1 menit)

Ucapkan dengan tenang; asesor biasanya menghargai kejujuran.

| Keterbatasan | Dampak | Yang saya lakukan sebagai gantinya |
|--------------|--------|-------------------------------------|
| Tidak ada **staging/production** tetap | Tidak ada deploy berkelanjutan ke cloud | Pipeline CI memverifikasi **build, test, image**; manifest & IaC siap untuk environment nyata |
| Cluster K8s mungkin **tidak selalu jalan** di laptop | `kubectl apply` / HPA tidak selalu didemokan live | Menyediakan manifest **lengkap + README**; menjelaskan alur apply dan kebutuhan **metrics-server** untuk HPA |
| **Registry & deploy otomatis** belum di-wire | CD “push + deploy” belum di GitHub Actions | **Pulumi** mendefinisikan **ECR**; workflow bisa ditambah job push/deploy bila ada **secrets** registry & kubeconfig |
| Secret **demo** di compose / `secret.example` | Bukan pola production | Dokumentasi: ganti via **Secret K8s** / **GitHub Secrets**; file sensitif di-`.gitignore` |

**Kalimat siap pakai:**  
*“Secara lokal saya prioritaskan **repeatability** dan **artefak yang bisa direview**: kode, tes, Docker, manifest Kubernetes, dan Pulumi. Integrasi penuh ke cloud saya arahkan sebagai langkah berikutnya begitu ada environment dan kredensial resmi.”*

---

## 3. Peta ke poin penilaian (inti presentasi ~3–5 menit)

### Go REST API

- **Sudah:** Endpoint sesuai rencana (publik, JWT, findings, internal), validasi JSON, respons error konsisten, pola **repository**, migrasi **golang-migrate**.
- **Batasan / next:** Lapisan **service** eksplisit bisa dipisah jika tim ingin domain logic lebih tebal dan handler lebih tipis.

### Autentikasi

- **Sudah:** **JWT** untuk user flow; **HTTP Basic Auth** untuk `/internal/*` dari **environment / Secret**, perilaku **503** jika Basic tidak dikonfigurasi (terdokumentasi).

### Docker

- **Sudah:** **Multi-stage**, user **non-root**, **HEALTHCHECK**; `docker-compose` untuk dev (API + Postgres).

### Kubernetes

- **Sudah (artefak):** Namespace, **Deployment** (2 replika, resource requests), **Service**, **ConfigMap** + **Secret** contoh, **liveness/readiness** ke `/health` dan `/ready`, **HPA** (CPU & memory), Postgres in-cluster untuk demo cluster.
- **Batasan lokal:** Ingress/LB di diagram arsitektur **belum** diwujudkan sebagai manifest (bisa disebut sebagai **langkah berikutnya** untuk exposure publik).

### CI/CD

- **Sudah:** **GitHub Actions** — `go mod verify`, `vet`, `test`, build binary, **build image**, compile program **Pulumi**.
- **Vision:** Tambah **golangci-lint**, lalu job **push image** (mis. ke ECR dari Pulumi) dan **deploy** (Helm/Kustomize + kubeconfig atau GitOps) saat akses tersedia.

### IaC

- **Sudah:** **Pulumi (Go)** — modul minimal **AWS ECR**, stack **dev**, README cara `preview`/`up`.
- **Vision:** Stack **staging/prod**, OIDC GitHub → AWS, atau provider lain sesuai kebijakan perusahaan.

### Dokumentasi & uji manual

- **Sudah:** Halaman **`/docs`** (panduan setup + referensi API), **Postman collection** JSON untuk import, **README** + **`k8s/README.md`** + **`infra/README.md`**.

---

## 4. Vision ke depan (1–2 menit)

1. **CD:** Setelah ada registry & cluster — extend workflow: build → push tag → deploy ke namespace terpisah (staging lalu prod).
2. **Keamanan:** Secret hanya dari **vault/Secret Manager**; hilangkan nilai demo dari repo; opsional **SOPS/Sealed Secrets**.
3. **Observabilitas:** Metrics (Prometheus), tracing, log terstruktur — menunjang HPA dan troubleshooting.
4. **Ingress & TLS:** Satu manifest Ingress + sertifikat (cert-manager) agar selaras dengan gambar arsitektur target.
5. **Arsitektur kode:** Layer **service** jika domain bertambah (approval workflow, notifikasi, dll.).

---

## 5. Alur demo singkat (jika ada waktu ~2 menit)

Pilih yang paling stabil di laptop kamu:

1. `docker compose up` **atau** `go run` + Postgres.
2. Browser: `http://localhost:8080/docs` — tunjuk pemisahan **panduan** vs **API**.
3. Postman / curl: **register → login → findings**; lalu **internal** dengan Basic Auth.
4. Tunjukkan di IDE: folder **`k8s/`** (Deployment, HPA), **`.github/workflows`**, **`infra/`** (Pulumi).

Jika cluster tidak jalan: buka file **`k8s/hpa-api.yaml`** dan **`deployment-api`** — jelaskan **probes** dan **mengapa** ada **requests** CPU/memory untuk HPA.

---

## 6. Penutup + undang diskusi

- **Ringkas:** “Yang dinilai dari brief sudah saya **wujudkan sebagai kode dan artefak**; bagian yang butuh **akun dan environment organisasi** saya jadikan **jalur evolusi** yang jelas.”
- **Tanya balik ke asesor (opsional):** preferensi registry (ECR/GCR/Artifactory), apakah Helm wajib, dan standar secret di perusahaan mereka.

---

*Sesuaikan durasi dengan format wawancara (10 vs 30 menit) dengan memperdalam demo atau deep-dive salah satu topik (mis. HPA vs connection pool).*
