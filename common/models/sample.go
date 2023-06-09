package models

type SampleAnalysis struct {
	Run       *string `json:"run"`
	Device    *string `json:"device"`
	Position  *int    `json:"position"`
	Completed *bool   `json:"completed"`

	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`

	Sample   Sample   `json:"sample"`
	Analysis Analysis `json:"analysis"`
}

type Sample struct {
	SampleID   string `json:"sample_id"`
	FullName   string `json:"full_name"`
	Sputalysed bool   `json:"sputalysed"`
	Comment    string `json:"comment"`

	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type Analysis struct {
	AnalysisID string `json:"analysis_id"`

	Analyt   string `json:"analyt"`
	Material string `json:"material"`
	Assay    string `json:"assay"`

	ReadyMix bool `json:"ready_mix"`
}
