package validator

import (
	"fmt"
	"reflect"
	"strings"
)

func Validate(model interface{}) errors {
	return validate(reflect.ValueOf(model))
}

type errors map[string][]error

func (e errors) add(name string, err error) {
	e[name] = append(e[name], err)
}

func validate(model reflect.Value) errors {
	errorsBag := make(errors)

	if model.Kind() != reflect.Struct {
		errorsBag.add("fatal", fmt.Errorf("models must be of type struct"))
		return errorsBag
	}

	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)
		fieldType := model.Type().Field(i)
		fieldTag := fieldType.Tag
		for _, validator := range extractValidators(fieldTag.Get("validate")) {
			if validateFunc, ok := validators[validator]; ok {
				if err := validateFunc(fieldType, field); err != nil {
					errorsBag.add(fieldType.Name, err)
				}
			}
		}
	}

	return errorsBag
}

func extractValidators(validators string) []string {
	return strings.Split(validators, ",")
}
