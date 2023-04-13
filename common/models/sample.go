package models

type SampleAnalysis struct {
	AnalysisId string `json:"analysis_id"`
	CreatedBy  string `json:"created_by"`
	CreatedAt  string `json:"created_at"`
}

type Sample struct {
	Tagesnummer string `json:"tagesnummer"`
	Name        string `json:"name"`

	AssociatedAnalyses []SampleAnalysis `json:"associated_analyses"`

	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type Analysis struct {
	AnalysisId string `json:"analysis_id"`

	// These three fields define the Analysis
	// Guaranteed by unique relationship
	Analyt   string `json:"analyt"`
	Material string `json:"material"`
	Assay    string `json:"assay"`

	ReadyMix bool `json:"ready_mix"`
}
