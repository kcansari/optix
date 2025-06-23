package validator

import (
	"fmt"
	"os"
)

type FileValidator interface {
	Validate(filename string) error
}

type BasicFileValidator struct{}

func (v *BasicFileValidator) Validate(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file '%s' does not exist", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file '%s' is not readable: %v", filename, err)
	}
	defer file.Close()

	return nil
}

type ValidatorStrategy struct {
	validator FileValidator
}

func NewValidatorStrategy(validator FileValidator) *ValidatorStrategy {
	return &ValidatorStrategy{validator: validator}
}

func (vs *ValidatorStrategy) ValidateFile(filename string) error {
	return vs.validator.Validate(filename)
}

func NewBasicFileValidator() *BasicFileValidator {
	return &BasicFileValidator{}
}
