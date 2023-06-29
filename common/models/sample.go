package models

type SampleAnalysis struct {
	Run      *string `json:"run"`
	Device   *string `json:"device"`
	Position *int    `json:"position"`
	Deleted  *bool   `json:"deleted"`

	CreatedBy *string `json:"created_by"`
	CreatedAt string  `json:"created_at"`

	Sample   Sample   `json:"sample"`
	Analysis Analysis `json:"analysis"`
}

type Sample struct {
	SampleID   string `json:"sample_id"`
	FullName   string `json:"full_name"`
	Sputalysed bool   `json:"sputalysed"`
	Comment    string `json:"comment"`
	Birthdate  string `json:"birthdate"`

	CreatedAt string  `json:"created_at"`
	CreatedBy *string `json:"created_by"`
}

type Analysis struct {
	AnalysisId  string `json:"analysis_id"`
	DisplayName string `json:"display_name"`

	ReadyMix bool `json:"ready_mix"`
	IsActive bool `json:"is_active"`
}
