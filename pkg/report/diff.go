package report

import "time"

// makeKey builds a stable identifier for a finding.
// We ignore timestamp & baseline because they vary run-to-run.
func makeKey(f Finding) string {
	return f.Service + "|" + f.Method + "|" + f.Payload
}

// ApplyBaseline tags findings in curr as BaseNew / BaseUnchanged,
// and (optionally) appends any findings that disappeared as BaseResolved.
func ApplyBaseline(curr, prev []Finding) {
	prevMap := make(map[string]Finding, len(prev))
	for _, f := range prev {
		prevMap[makeKey(f)] = f
	}

	// Track which previous keys we’ve matched.
	matched := make(map[string]struct{}, len(prevMap))

	// 1) Mark current slice.
	for i := range curr {
		key := makeKey(curr[i])
		if _, ok := prevMap[key]; ok {
			curr[i].Baseline = BaseUnchanged
			matched[key] = struct{}{}
		} else {
			curr[i].Baseline = BaseNew
		}
	}

	// Add “resolved” findings (present before, gone now).
	for key, old := range prevMap {
		if _, ok := matched[key]; ok {
			continue
		}
		curr = append(curr, Finding{
			Service:   old.Service,
			Method:    old.Method,
			Payload:   old.Payload,
			Error:     old.Error,
			Severity:  old.Severity,
			Baseline:  BaseResolved,
			Timestamp: time.Now(),
		})
	}
}
