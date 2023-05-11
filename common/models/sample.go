package models

import (
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullInt64
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// MarshalJSON for NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type SampleAnalysis struct {
	Run      NullString `json:"run"`
	Device   NullString `json:"device"`
	Position NullInt64  `json:"position,omitempty"`

	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`

	Sample   Sample   `json:"sample"`
	Analysis Analysis `json:"analysis"`
}

type Sample struct {
	SampleID   string     `json:"sample_id"`
	FullName   string     `json:"full_name"`
	Sputalysed bool       `json:"sputalysed"`
	Comment    NullString `json:"comment"`

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
