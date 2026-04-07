# Kubernetes — Findings API

Manifest untuk namespace `findings`: Postgres in-cluster (demo), **ConfigMap**, **Secret** (contoh), **Deployment** + **Service** API, **HPA** (CPU & memory).

## Isi

| File | Keterangan |
|------|------------|
| `namespace.yaml` | Namespace `findings` |
| `postgres.yaml` | PVC + Deployment + Service Postgres |
| `configmap.yaml` | `HTTP_ADDR`, `SERVICE_NAME`, `SERVICE_VERSION` |
| `secret.example.yaml` | Secret `findings-db` + `findings-api` (nilai demo — ganti untuk prod) |
| `deployment-api.yaml` | API: 2 replika, probes `/health` & `/ready`, resource requests untuk HPA |
| `service-api.yaml` | ClusterIP `:8080` |
| `hpa-api.yaml` | `minReplicas: 2`, `maxReplicas: 10`, target CPU 70%, memory 80% |
| `kustomization.yaml` | Kustomize (tanpa secret — terapkan secret manual) |

## Urutan apply

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.example.yaml
kubectl apply -k k8s/
```

Atau tanpa kustomize:

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secret.example.yaml
kubectl apply -f k8s/postgres.yaml -f k8s/configmap.yaml -f k8s/deployment-api.yaml -f k8s/service-api.yaml -f k8s/hpa-api.yaml
```

**HPA** membutuhkan **metrics-server** terpasang di cluster (`kubectl top pods` harus jalan).

## Image API

`deployment-api.yaml` memakai `findings-api:latest`. Build dari root repo:

```bash
docker build -t findings-api:latest .
```

**kind:**

```bash
kind load docker-image findings-api:latest
```

**minikube:**

```bash
minikube image load findings-api:latest
```

Atau set `imagePullPolicy: Never` bila semua image hanya lokal (sesuaikan di manifest).

Untuk cluster cloud: push ke registry (mis. ECR dari Pulumi `infra/`), lalu ganti field `image` + `imagePullPolicy: Always`.

## Akses dari laptop

```bash
kubectl -n findings port-forward svc/findings-api 8080:8080
curl -s http://127.0.0.1:8080/health
```

## Production / managed DB

Hapus atau jangan apply `postgres.yaml`; set `DATABASE_URL` di Secret `findings-api` ke URL database terkelola. Sesuaikan `JWT_SECRET` dan Basic Auth.

## Secret sendiri (disarankan)

```bash
cp k8s/secret.example.yaml k8s/secret.yaml
# edit k8s/secret.yaml — file secret.yaml di-.gitignore
kubectl apply -f k8s/secret.yaml
```
