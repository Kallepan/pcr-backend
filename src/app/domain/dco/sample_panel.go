package dco

type ResetSamplePanelRequest struct {
	PanelID  string `json:"panel_id" binding:"required"`
	SampleID string `json:"sample_id" binding:"required"`
}

type UpdateSamplePanelRequest struct {
	Deleted *bool `json:"deleted" binding:"required"`
}

type AddSamplePanelRequest struct {
	SampleID string `json:"sample_id" binding:"required"`
	PanelID  string `json:"panel_id" binding:"required"`
}
