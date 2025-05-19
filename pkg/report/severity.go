package report

// Severity describes how bad a finding is once classified.
type Severity string

const (
	SevCritical Severity = "critical"
	SevHigh     Severity = "high"
	SevLow      Severity = "low"
	SevNone     Severity = "none"
)
