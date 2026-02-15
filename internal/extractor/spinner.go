package extractor

// SpinnerExtractor handles spinner definition files (spinners/*.yaml).
// ALL spinners are persona-specific -- each persona defines its own
// loading animation verbs. There is no core spinner content.
//
// Unlike other extractors, spinners are YAML files (not markdown).
// They are classified via fileType routing and tracked as persona assets
// without Document parsing. The orchestrator handles them directly in
// the Walk switch (like commands and hooks).
type SpinnerExtractor struct{}

// NewSpinnerExtractor creates a SpinnerExtractor.
func NewSpinnerExtractor() *SpinnerExtractor {
	return &SpinnerExtractor{}
}
