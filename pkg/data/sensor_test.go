package data

import (
	"testing"
)

func Test(t *testing.T) {
	testData := SensorData{
		Id:         "1234",
		Value:      -1,
		SourceType: "temperature",
	}

	testInvalidData := SensorData{
		Id:         "1234",
		Value:      -1,
		SourceType: "light",
	}

	err := validateData(testData)

	if err != nil {
		t.Error(err)
	}

	err = validateData(testInvalidData)

	if err == nil {
		t.Error("test data is invalid")
	}

	if err.Error() != "light source type is not supported" {
		t.Errorf("Got %s Want: %s", err.Error(), "light source type is not supported")
	}
}
