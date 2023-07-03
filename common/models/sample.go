package models

type SampleAnalysis struct {
	Run      *string `json:"run"`
	Device   *string `json:"device"`
	Position *int    `json:"position"`
	RunDate  *string `json:"run_date"`
	Deleted  *bool   `json:"deleted"`

	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`

	Sample Sample `json:"sample"`
	Panel  Panel  `json:"panel"`
}

type Sample struct {
	SampleId   string  `json:"sample_id"`
	FullName   string  `json:"full_name"`
	Sputalysed bool    `json:"sputalysed"`
	Comment    *string `json:"comment"`
	Birthdate  string  `json:"birthdate"`

	CreatedAt string  `json:"created_at"`
	CreatedBy *string `json:"created_by"`
}

type Panel struct {
	PanelId     string `json:"panel_id"`
	DisplayName string `json:"display_name"`
	ReadyMix    *bool  `json:"ready_mix"`
}
