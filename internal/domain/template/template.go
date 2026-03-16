// Package template defines the domain interface and data types for template rendering.
package template

// TemplateData is the variable namespace available inside every template.
// Field names use lowercase keys to match the dot-notation in .tmpl files
// (e.g. {{ .service.name }}, {{ .meta.cli_version }}).
type TemplateData struct {
	Service    ServiceData    `json:"service"`
	Company    CompanyData    `json:"company"`
	Resilience ResilienceData `json:"resilience"`
	SLO        SLOData        `json:"slo"`
	Cost       CostData       `json:"cost"`
	Infra      InfraData      `json:"infra"`
	Meta       MetaData       `json:"meta"`
	Plugins    []string       `json:"plugins_used"`
	Answers    map[string]any `json:"answers"`
}

// ToMap converts TemplateData to a map[string]any using the lowercase key convention
// expected by the .tmpl files (e.g. {{ .service.name }}).
func (d *TemplateData) ToMap() map[string]any {
	return map[string]any{
		"service": map[string]any{
			"name":               d.Service.Name,
			"module":             d.Service.Module,
			"language":           d.Service.Language,
			"framework":          d.Service.Framework,
			"service_type":       d.Service.ServiceType,
			"compliance_profile": d.Service.ComplianceProfile,
			"mesh_mode":          d.Service.MeshMode,
			"team":               d.Service.Team,
			"namespace":          d.Service.Namespace,
			"environment":        d.Service.Environment,
			"version":            d.Service.Version,
		},
		"company": map[string]any{
			"name":               d.Company.Name,
			"registry":           d.Company.Registry,
			"vault_addr":         d.Company.VaultAddr,
			"trust_domain":       d.Company.TrustDomain,
			"correlation_header": d.Company.CorrelationHeader,
		},
		"resilience": map[string]any{
			"circuit_breaker_threshold": d.Resilience.CircuitBreakerThreshold,
			"timeout_db_ms":             d.Resilience.TimeoutDBMs,
			"timeout_http_ms":           d.Resilience.TimeoutHTTPMs,
			"timeout_kafka_ms":          d.Resilience.TimeoutKafkaMs,
			"retry_max_attempts":        d.Resilience.RetryMaxAttempts,
			"retry_backoff_base_ms":     d.Resilience.RetryBackoffBaseMs,
		},
		"slo": map[string]any{
			"availability_target": d.SLO.AvailabilityTarget,
			"p99_latency_ms":      d.SLO.P99LatencyMs,
			"error_budget_policy": d.SLO.ErrorBudgetPolicy,
		},
		"cost": map[string]any{
			"centre":             d.Cost.Centre,
			"team":               d.Cost.Team,
			"monthly_budget_usd": d.Cost.MonthlyBudgetUSD,
		},
		"infra": map[string]any{
			"cloud":          d.Infra.Cloud,
			"region":         d.Infra.Region,
			"data_residency": d.Infra.DataResidency,
			"registry":       d.Infra.Registry,
		},
		"meta": map[string]any{
			"cli_version":  d.Meta.CLIVersion,
			"generated_at": d.Meta.GeneratedAt,
			"generated_by": d.Meta.GeneratedBy,
		},
		"plugins_used": d.Plugins,
		"answers":      d.Answers,
	}
}

// ServiceData holds service-level variables.
type ServiceData struct {
	Name              string `json:"name"`
	Module            string `json:"module"`
	Language          string `json:"language"`
	Framework         string `json:"framework"`
	ServiceType       string `json:"service_type"`
	ComplianceProfile string `json:"compliance_profile"`
	MeshMode          bool   `json:"mesh_mode"`
	Team              string `json:"team"`
	Namespace         string `json:"namespace"`
	Environment       string `json:"environment"`
	Version           string `json:"version"`
}

// CompanyData holds organisation-level variables.
type CompanyData struct {
	Name              string `json:"name"`
	Registry          string `json:"registry"`
	VaultAddr         string `json:"vault_addr"`
	TrustDomain       string `json:"trust_domain"`
	CorrelationHeader string `json:"correlation_header"`
}

// ResilienceData holds resilience configuration variables.
type ResilienceData struct {
	CircuitBreakerThreshold int `json:"circuit_breaker_threshold"`
	TimeoutDBMs             int `json:"timeout_db_ms"`
	TimeoutHTTPMs           int `json:"timeout_http_ms"`
	TimeoutKafkaMs          int `json:"timeout_kafka_ms"`
	RetryMaxAttempts        int `json:"retry_max_attempts"`
	RetryBackoffBaseMs      int `json:"retry_backoff_base_ms"`
}

// SLOData holds SLO target variables.
type SLOData struct {
	AvailabilityTarget string `json:"availability_target"`
	P99LatencyMs       int    `json:"p99_latency_ms"`
	ErrorBudgetPolicy  string `json:"error_budget_policy"`
}

// CostData holds FinOps variables.
type CostData struct {
	Centre           string `json:"centre"`
	Team             string `json:"team"`
	MonthlyBudgetUSD int    `json:"monthly_budget_usd"`
}

// InfraData holds infrastructure variables.
type InfraData struct {
	Cloud         string `json:"cloud"`
	Region        string `json:"region"`
	DataResidency string `json:"data_residency"`
	Registry      string `json:"registry"`
}

// MetaData holds generation metadata variables.
type MetaData struct {
	CLIVersion  string `json:"cli_version"`
	GeneratedAt string `json:"generated_at"`
	GeneratedBy string `json:"generated_by"`
}

// Engine renders named templates with TemplateData to the filesystem.
type Engine interface {
	// Render renders the named template and writes the result to outputPath.
	Render(templateName string, data *TemplateData, outputPath string) error
}
