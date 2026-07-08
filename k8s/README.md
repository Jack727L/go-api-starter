# Kubernetes manifests (skeleton)

Minimal manifests to run the API plus an in-cluster Redis. These are **cloud-agnostic**: nothing here is AWS-specific except what you put in the image URL and the secrets. Postgres is assumed to be **outside the cluster** (a managed database), wired in via the `postgres-secrets` Secret.

## Files

| File | What it is |
| :--- | :--- |
| `namespace.yaml` | `development` namespace everything else lives in |
| `backend-deployment.yaml` | API Deployment — set `image:` to your registry before applying |
| `backend-service.yaml` | ClusterIP service, port 80 → container 3000 |
| `backend-ingress.yaml` | Traefik ingress for `api.local` — change host/class for your cluster |
| `redis.yaml` | Single-replica Redis + Service (jobs queue; no persistence) |
| `backend-secret.example.yaml` | Template for `postgres-secrets` — **do not commit a real copy** |

## Deploy

```bash
kubectl apply -f k8s/namespace.yaml

# Create secrets (or copy backend-secret.example.yaml — real copies are gitignored)
kubectl create secret generic postgres-secrets -n development \
  --from-literal=POSTGRES_HOST=... \
  --from-literal=POSTGRES_PORT=5432 \
  --from-literal=POSTGRES_DATABASE=... \
  --from-literal=POSTGRES_USER=... \
  --from-literal=POSTGRES_PASSWORD=...

# Private registry only:
# kubectl create secret docker-registry registry-pull-secret -n development \
#   --docker-server=... --docker-username=... --docker-password=...

kubectl apply -f k8s/
```

Migrations are not applied by the cluster — run `sqitch deploy` against your database from CI or a one-off job (see the AWS section of the root README; the same applies to any cloud).

## Cloud provider mapping

The manifests do not care which cloud you use. What changes:

| Concern | AWS | GCP | Notes |
| :--- | :--- | :--- | :--- |
| Cluster | EKS | GKE | Manifests identical |
| Image registry | ECR (`<account>.dkr.ecr.<region>.amazonaws.com/...`) | Artifact Registry (`<region>-docker.pkg.dev/<project>/...`) | Set in `backend-deployment.yaml` |
| Registry auth | `imagePullSecrets` or IRSA/pod identity | Workload Identity (GKE pulls Artifact Registry natively — often no pull secret needed) | Remove `imagePullSecrets` if unused |
| Postgres | RDS | Cloud SQL | Both are just host/port/credentials in `postgres-secrets`. Cloud SQL often uses the [Auth Proxy sidecar](https://cloud.google.com/sql/docs/postgres/connect-kubernetes-engine) instead of a direct host |
| Redis | ElastiCache (or in-cluster) | Memorystore (or in-cluster) | For managed Redis, delete `redis.yaml` and point `REDIS_HOST`/`REDIS_PORT` at the managed endpoint |
| Ingress / LB | ALB (aws-load-balancer-controller) | GCLB (GKE Ingress) or any controller | Change `ingressClassName` and host; `traefik` here is a generic default |
| Secrets backing | Secrets Manager + External Secrets Operator (optional) | Secret Manager + External Secrets Operator (optional) | Plain k8s Secrets work everywhere; ESO is an upgrade path, not a requirement |
