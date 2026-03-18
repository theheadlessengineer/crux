# crux-plugin-terraform-azure

Generates Azure Terraform infrastructure for `{{ service.name }}`: Azure Database, Redis Cache, and IAM modules.

## Generated Files

| File | Purpose |
|---|---|
| `infra/terraform/azure/main.tf` | Provider config and platform module calls |
| `infra/terraform/azure/variables.tf` | Input variables |
| `infra/terraform/azure/outputs.tf` | Output values |
| `infra/terraform/azure/backend.tf` | Azure Storage remote state backend |

## Configuration

| Question | Default | Description |
|---|---|---|
| `azure_subscription_id` | `00000000-...` | Azure subscription ID |
| `azure_location` | `westeurope` | Primary Azure location |
| `azure_remote_state_storage_account` | `companyterraformstate` | Storage account for state |
