package required

import (
	"fmt"
	"strings"
	"testing"

	validator "github.com/mohamed-elalem/struct-validator"
)

func TestValidator(t *testing.T) {
	x := new(int)
	*x = 5

	tests := []struct {
		input  interface{}
		output validator.Errors
	}{
		{
			input: struct {
				a string `validate:"required"`
			}{},
			output: validator.Errors{"a": []error{fmt.Errorf("required")}},
		},
		{
			input: struct {
				a string `validate:"required"`
			}{"a"},
			output: validator.Errors{},
		},
		{
			input: struct {
				ptr *int `validate:"required"`
			}{x},
			output: validator.Errors{},
		},
		{
			input: struct {
				ptr *int `validate:"required"`
			}{},
			output: validator.Errors{"ptr": []error{fmt.Errorf("required")}},
		},
		{
			input: struct {
				interf interface{} `validate:"required"`
			}{},
			output: validator.Errors{"interf": []error{fmt.Errorf("required")}},
		},
		{
			input: struct {
				interf interface{} `validate:"required"`
			}{*x},
			output: validator.Errors{},
		},
	}

	for _, test := range tests {
		errors := validator.Validate(test.input)
		if len(errors) != len(test.output) {
			t.Errorf("validator.Validate(%+v) got %+v want %+v", test.input, errors, test.output)
		}
		for field, expectedErrors := range test.output {
			if errs, ok := errors[field]; !ok {
				t.Errorf("expected entry %s to exist in error", field)
			} else {
				for _, expectedError := range expectedErrors {
					found := false
					for _, err := range errs {
						if strings.Contains(err.Error(), expectedError.Error()) {
							found = true
							break
						}
					}
					if found == false {
						t.Errorf("expected error %v to exist on field %s", expectedError, field)
					}
				}
			}
		}
	}
}
