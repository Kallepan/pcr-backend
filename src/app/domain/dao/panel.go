package dao

type Panel struct {
	PanelID     string `json:"panel_id"`
	DisplayName string `json:"display_name"`
	ReadyMix    *bool  `json:"ready_mix"`
}
