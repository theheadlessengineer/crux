# crux-plugin-terraform-aws

**Tier:** 1 (Official) | **Version:** 1.0.0 | **Phase:** Pilot

Terraform AWS infrastructure for crux-generated services — provider config, remote state backend, and variable definitions.

## Questions

| ID | Type | Prompt | Default |
|---|---|---|---|
| `aws_region` | select | Primary AWS region | `eu-west-1` |
| `aws_environments` | select | Environments to provision | `dev,staging,prod` |
| `aws_remote_state_bucket` | input | Terraform state S3 bucket | `company-terraform-state` |

## Generated Files

| File | Description |
|---|---|
| `infra/terraform/main.tf` | Provider config with cost allocation default tags |
| `infra/terraform/backend.tf` | S3 remote state with DynamoDB locking |
| `infra/terraform/variables.tf` | Region, environment, service name variables |
| `infra/terraform/outputs.tf` | Service name and region outputs |

## Cost Allocation Tags

All AWS resources receive these tags automatically via `default_tags`:

```
service     = <service-name>
team        = <team>
environment = <environment>
cost-centre = engineering
managed-by  = terraform
```
