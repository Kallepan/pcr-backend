package samples

import (
	"database/sql"

	"gitlab.com/kaka/pcr-backend/common/database"
	"gitlab.com/kaka/pcr-backend/common/models"
)

func FetchSampleInformationFromDatabase(sampleID string) (*models.Sample, error) {
	var sample models.Sample

	query :=
		`SELECT sample_id,samples.full_name,created_at,users.username,birthdate,sputalysed,comment
		FROM samples 
		LEFT JOIN users ON samples.created_by = users.user_id
		WHERE sample_id = $1;`

	row := database.Instance.QueryRow(query, sampleID)

	if err := row.Scan(&sample.SampleId, &sample.FullName, &sample.CreatedAt, &sample.CreatedBy, &sample.Birthdate, &sample.Sputalysed, &sample.Comment); err != nil {
		return nil, err
	}

	return &sample, nil
}

func SampleExists(sample_id string) bool {
	query := `
		SELECT EXISTS(
			SELECT sample_id 
			FROM samples 
			WHERE sample_id = $1
		)`

	var exists bool
	err := database.Instance.QueryRow(query, sample_id).Scan(&exists)

	if err != nil && err != sql.ErrNoRows {
		return false
	}

	return exists
}
