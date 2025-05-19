package report

// Baseline marks how a finding compares to a previous (baseline) scan.
type Baseline string

const (
	BaseNew       Baseline = "new"
	BaseUnchanged Baseline = "unchanged"
	BaseResolved  Baseline = "resolved"
)
