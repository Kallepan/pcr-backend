package samplespanels

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/kaka/pcr-backend/packages/projectpath"
)

func TestGetFormattedBirthdate(t *testing.T) {
	// Test nil birthdate
	assert.Equal(t, "NA", getFormattedBirthdate(nil))

	// Test non-nil birthdate
	birthdate := "2020-01-01T00:00:00Z"
	assert.Equal(t, "2020-01-01", getFormattedBirthdate(&birthdate))
}

func TestGetFormattedSampleID(t *testing.T) {
	// Test non-numeric sample id
	sampleID := "ABC"
	assert.Equal(t, "ABC", getFormattedSampleID(sampleID))

	// Test numeric sample id
	sampleID = "123"
	assert.Equal(t, "123", getFormattedSampleID(sampleID))

	eightLetterSampleID := "12345678"
	assert.Equal(t, "1234 5678", getFormattedSampleID(eightLetterSampleID))

	twelveLetterSampleID := "123456789012"
	assert.Equal(t, "1234 567890 12", getFormattedSampleID(twelveLetterSampleID))

	thirteenLetterSampleID := "1234567890123"
	assert.Equal(t, "1234567890123", getFormattedSampleID(thirteenLetterSampleID))
}

func TestCreateCopy(t *testing.T) {
	root := projectpath.Root
	template := fmt.Sprintf("%s/templates/v1.xlsm", root)
	println(os.Getwd())
	outputPath, err := createCopy(template)

	assert.Nil(t, err)
	assert.NotEmpty(t, outputPath)

	// Delete the file
	err = os.Remove(*outputPath)

	assert.Nil(t, err)
}