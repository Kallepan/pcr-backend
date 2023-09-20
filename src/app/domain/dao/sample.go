package dao

type Sample struct {
	SampleId   string  `json:"sample_id"`
	FullName   string  `json:"full_name"`
	Sputalysed bool    `json:"sputalysed"`
	Comment    *string `json:"comment"`
	Birthdate  *string `json:"birthdate"`
	Material   *string `json:"material"`

	CreatedAt string  `json:"created_at"`
	CreatedBy *string `json:"created_by"`

	Panels string `json:"panels"`
}
