# crux-plugin-terraform-gcp

Generates GCP Terraform infrastructure for `{{ service.name }}`: Cloud SQL, Memorystore, and IAM modules.

## Generated Files

| File | Purpose |
|---|---|
| `infra/terraform/gcp/main.tf` | Provider config and platform module calls |
| `infra/terraform/gcp/variables.tf` | Input variables |
| `infra/terraform/gcp/outputs.tf` | Output values |
| `infra/terraform/gcp/backend.tf` | GCS remote state backend |

## Configuration

| Question | Default | Description |
|---|---|---|
| `gcp_project_id` | `company-project` | GCP project ID |
| `gcp_region` | `europe-west1` | Primary GCP region |
| `gcp_remote_state_bucket` | `company-terraform-state` | GCS state bucket |
