package validator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRegisterValidator(t *testing.T) {
	fun := func(t reflect.StructField, v reflect.Value) error { return nil }
	handlerName := "test"
	RegisterValidator(handlerName, fun)
	defer unregisterValidator(handlerName)
	if _, ok := validators["test"]; !ok {
		t.Errorf("Expected handler %s to be registered", handlerName)
	}
}

func TestValidateWithCorrectValidators(t *testing.T) {
	type TestStruct struct {
		a string `validate:"required,email"`
		b string `validate:"email"`
	}

	testStruct := TestStruct{}

	emailValidatorCalls, requiredValidatorCalls := 0, 0
	emailValidator := func(_ reflect.StructField, _ reflect.Value) error {
		emailValidatorCalls++
		return nil
	}
	requiredValidator := func(_ reflect.StructField, _ reflect.Value) error {
		requiredValidatorCalls++
		return nil
	}

	RegisterValidator("required", requiredValidator)
	RegisterValidator("email", emailValidator)

	defer unregisterValidator("required")
	defer unregisterValidator("email")

	Validate(testStruct)

	if emailValidatorCalls != 2 {
		t.Errorf("email validator should been called 2 times got %d", emailValidatorCalls)
	}
	if requiredValidatorCalls != 1 {
		t.Errorf("required validator should have been called 1 time got %d", requiredValidatorCalls)
	}
}

func TestUnregisterValidator(t *testing.T) {
	fun := func(_ reflect.StructField, _ reflect.Value) error { return nil }
	handlerName := "test"
	RegisterValidator(handlerName, fun)
	if _, ok := validators[handlerName]; !ok {
		t.Errorf("%s should be a registered validator", handlerName)
	}
	unregisterValidator(handlerName)
	if _, ok := validators[handlerName]; ok {
		t.Errorf("%s should be an unregistered validator", handlerName)
	}
}

func TestExtractValidators(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"", []string{}},
		{"a", []string{"a"}},
	}

	for _, test := range tests {
		validators := extractValidators(test.input)
		for i := 0; i < len(test.output); i++ {
			if test.output[i] != validators[i] {
				t.Errorf("extractValidators(%q) got %+v want %+v", test.input, validators, test.output)
				break
			}
		}
	}
}

func TestValidationFailure(t *testing.T) {
	type TestStruct struct {
		a string `validate:"required,email"`
		b string `validate:"email"`
	}

	emailValidatorFunc := func(typ reflect.StructField, val reflect.Value) error {
		return fmt.Errorf("email")
	}

	requiredValidatorFunc := func(typ reflect.StructField, val reflect.Value) error {
		return fmt.Errorf("required")
	}

	RegisterValidator("required", requiredValidatorFunc)
	RegisterValidator("email", emailValidatorFunc)

	defer unregisterValidator("required")
	defer unregisterValidator("email")

	errors := Validate(TestStruct{})

	if len(errors) == 0 {
		t.Errorf("expected errors to be populated got %+v", errors)
	}

	if _, ok := errors["a"]; !ok {
		t.Errorf(`expected entry "a" to exist in %+v got nil`, errors)
	} else if len(errors["a"]) != 2 {
		t.Errorf(`expected entry "a" to contain 2 errors got %d, %+v`, len(errors["a"]), errors["a"])
		errorsFound := 0
		for _, err := range errors["a"] {
			switch err.Error() {
			case "email":
				fallthrough
			case "required":
				errorsFound++
			default:
				t.Errorf("error %q is not expected", err.Error())
			}
		}

		if errorsFound != 2 {
			t.Errorf("expected %d errors got %d", 2, errorsFound)
		}
	}
}

func TestValidateWrongModelType(t *testing.T) {
	errors := Validate(3)
	if len(errors) == 0 {
		t.Errorf("expected errors %+v not to be empty", errors)
	}
	if _, ok := errors["fatal"]; !ok {
		t.Errorf("expected fatal error to exist on %+v", errors)
	}
}

func BenchmarkValidate(b *testing.B) {
	type TestStruct struct {
		a string `validate:"required,email"`
		b string `validate:"email"`
	}

	emailValidatorFunc := func(typ reflect.StructField, val reflect.Value) error {
		return fmt.Errorf("email")
	}

	requiredValidatorFunc := func(typ reflect.StructField, val reflect.Value) error {
		return fmt.Errorf("required")
	}

	RegisterValidator("required", requiredValidatorFunc)
	RegisterValidator("email", emailValidatorFunc)

	defer unregisterValidator("required")
	defer unregisterValidator("email")

	testStruct := TestStruct{}

	for i := 0; i < b.N; i++ {
		Validate(testStruct)
	}
}
