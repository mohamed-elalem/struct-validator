package required

import (
	"fmt"
	"reflect"

	validator "github.com/mohamed-elalem/struct-validator"
)

var requiredValidator validator.ValidatorHandler = func(t reflect.StructField, v reflect.Value) error {
	err := fmt.Errorf("field %s is required", t.Name)
	switch v.Kind() {
	case reflect.String:
		if v.String() == "" {
			return err
		}
	case reflect.Ptr:
		if v.IsNil() {
			return err
		}
	default:
		if v.IsNil() || !v.IsValid() {
			return err
		}
	}

	return nil
}

func init() {
	validator.RegisterValidator("required", requiredValidator)
}
