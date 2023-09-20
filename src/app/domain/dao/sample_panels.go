package dao

type SampleAnalysis struct {
	Run      *string `json:"run"`
	Device   *string `json:"device"`
	Position *int    `json:"position"`
	RunDate  *string `json:"run_date"`

	CreatedBy *string `json:"created_by"`
	CreatedAt string  `json:"created_at"`

	Sample Sample `json:"sample"`
	Panel  Panel  `json:"panel"`
}
