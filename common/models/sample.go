package models

type SampleAnalysis struct {
	ID        string `json:"id"`
	Run       string `json:"run"`
	Device    string `json:"device"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`

	// These three fields define the Analysis
	// Guaranteed by unique relationship
	Analyt   string `json:"analyt"`
	Material string `json:"material"`
	Assay    string `json:"assay"`
	ReadyMix bool   `json:"ready_mix"`
}

type Sample struct {
	Tagesnummer string `json:"tagesnummer"`
	Name        string `json:"name"`

	AssociatedAnalyses []SampleAnalysis `json:"associated_analyses"`

	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type Analysis struct {
	// These three fields define the Analysis
	// Guaranteed by unique relationship
	Analyt   string `json:"analyt"`
	Material string `json:"material"`
	Assay    string `json:"assay"`

	ReadyMix bool `json:"ready_mix"`
}
