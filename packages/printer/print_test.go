package printer

import (
	"testing"
)

// Test Label Creation
type labelTest struct {
	printData PrintData
	expected  string
}

var printData_one = PrintData{
	Position: "N",
	SampleID: "1234567890",
	Name:     "John Doe",
	PanelID:  "Panel 1",
	Device:   "Device 1",
	Run:      "Run 1",
	Date:     "2020-01-01",
}
var expected_one = `q256
N
A150,40,0,5,1,1,N,"N"
A30,0,0,2,1,1,N,"1234567890"
A30,20,0,2,1,1,N,"John Doe"
A30,50,0,2,1,1,N,"Panel 1"
A30,70,0,2,1,1,N,"Device 1Run 1"
A30,100,0,1,1,1,N,"2020-01-01"
P1
`

var printData_two = PrintData{
	Position: "N",
	SampleID: "1234567890",
	Name:     "John Doeä",
	PanelID:  "Panel 1",
	Device:   "Device 1",
	Run:      "Run 1",
	Date:     "2020-01-01",
}
var expected_two = `q256
N
A150,40,0,5,1,1,N,"N"
A30,0,0,2,1,1,N,"1234567890"
A30,20,0,2,1,1,N,"John Doe?"
A30,50,0,2,1,1,N,"Panel 1"
A30,70,0,2,1,1,N,"Device 1Run 1"
A30,100,0,1,1,1,N,"2020-01-01"
P1
`

var labelTests = []labelTest{
	{printData_one, expected_one},
	{printData_two, expected_two},
}

// TestPrint tests the Print function
func TestLabelCreation(t *testing.T) {
	for _, lt := range labelTests {
		actual, err := lt.printData.createLabel(globalTemplate)
		if err != nil {
			t.Errorf("error creating label: %s", err)
		}

		if actual != lt.expected {
			t.Errorf("expected %s, actual %s", actual, lt.expected)
		}
	}
}
