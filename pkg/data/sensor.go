package data

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	TemperatureSource = "temperature"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// SensorData stores sensor readings for different source types eg. temperature, light etc.
type SensorData struct {
	Id         string  `json:"id" validate:"required"`
	Value      float64 `json:"value" validate:"required,gte=-30,lte=100"`
	SourceType string  `json:"source_type" validate:"required"`
}

func validateData(data SensorData) error {
	if data.SourceType != TemperatureSource {
		return fmt.Errorf("%s source type is not supported", data.SourceType)
	}
	err := validate.Struct(data)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return fmt.Errorf("%s", ValidationError(validationErrors))
	}
	return nil
}

func ValidationError(errs validator.ValidationErrors) string {
	var s []string
	for _, e := range errs {
		s = append(s, e.Error())
	}
	return strings.Join(s, "\n")
}
