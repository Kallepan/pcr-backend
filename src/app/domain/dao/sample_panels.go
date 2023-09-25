package dao

type SamplePanel struct {
	Run      *string `json:"run"`
	Device   *string `json:"device"`
	Position *int    `json:"position"`
	RunDate  *string `json:"run_date"`

	CreatedBy *string `json:"created_by"`
	CreatedAt string  `json:"created_at"`

	Sample Sample `json:"sample"`
	Panel  Panel  `json:"panel"`
}

type Statistic struct {
	PanelID string `json:"panel_id"`
	Count   int    `json:"count"`
}
