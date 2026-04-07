# Infrastruktur (Pulumi + AWS)

Modul **IaC** terpisah dari aplikasi Go. Saat ini hanya ada stack contoh **`dev`** (sandbox / pengembangan). **Staging dan production** belum dibuat — nanti cukup tambah stack baru, misalnya:

```bash
pulumi stack init staging
pulumi stack init prod
```

## Apa yang didefinisikan

- **AWS ECR** — repository Docker untuk image `findings-api` (nama bisa diubah lewat config).
- Scan on push diaktifkan; tag `Project` / `Managed` untuk pelacakan.

Tanpa akun AWS, kamu tetap bisa **mengompilasi** program Pulumi (`go build`) dan melewati CI; **`pulumi preview` / `up`** membutuhkan kredensial AWS.

## Prasyarat (untuk preview / deploy)

1. [Pulumi CLI](https://www.pulumi.com/docs/install/)
2. [AWS CLI](https://aws.amazon.com/cli/) atau variabel lingkungan `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`
3. Go 1.22+

## Langkah cepat

Dari folder `infra/`:

```bash
# Login Pulumi (Pulumi Cloud gratis / atau backend file — lihat dokumentasi Pulumi)
pulumi login

# Pakai stack dev (sudah ada file Pulumi.dev.yaml dengan region contoh)
pulumi stack select dev  # atau: pulumi stack init dev

# Opsional: ubah nama repo ECR
pulumi config set repositoryName findings-api

# Rencana perubahan (butuh kredensial AWS)
pulumi preview

# Terapkan (membuat resource nyata di AWS — hati-hati biaya)
pulumi up
```

Output penting setelah deploy: `ecrRepositoryUrl` (untuk `docker push` dari pipeline nanti).

## Tanpa environment cloud

- Cukup commit kode di repo sebagai bukti **IaC**; reviewer bisa baca `main.go`.
- Jalankan `go build` atau biarkan job **Pulumi IaC (compile)** di GitHub Actions memverifikasi kompilasi.

## Catatan

- Region default di `Pulumi.dev.yaml`: `ap-southeast-1` (sesuaikan kebutuhan).
- Jangan commit secret AWS; gunakan OIDC / GitHub Secrets bila otomasi CI perlu `pulumi up`.
