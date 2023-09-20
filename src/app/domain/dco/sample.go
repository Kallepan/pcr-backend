package dco

import "errors"

type AddSampleRequest struct {
	SampleId   string `json:"sample_id" binding:"required"`
	FullName   string `json:"full_name" binding:"required"`
	Sputalysed bool   `json:"sputalysed"`
	Comment    string `json:"comment,omitempty"`
	Birthdate  string `json:"birthdate" binding:"required"`
	Material   string `json:"material" binding:"required,uppercase,alpha"`
}

func (a *AddSampleRequest) Validate() error {
	a.Sputalysed = a.Sputalysed || false

	return nil
}

type UpdateSampleRequest struct {
	FullName   *string `json:"full_name,omitempty"`
	Sputalysed *bool   `json:"sputalysed,omitempty"`
	Comment    *string `json:"comment,omitempty"`
}

func (u *UpdateSampleRequest) Validate() error {
	if u.FullName == nil && u.Sputalysed == nil && u.Comment == nil {
		return errors.New("no fields to update")
	}

	return nil
}
