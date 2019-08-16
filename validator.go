package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

func Validate(model interface{}) Errors {
	return validate(reflect.ValueOf(model))
}

var mu sync.Mutex

type Errors map[string][]error

func (e Errors) add(name string, err error) {
	mu.Lock()
	e[name] = append(e[name], err)
	mu.Unlock()
}

func validate(model reflect.Value) Errors {
	errorsBag := make(Errors)
	wg := sync.WaitGroup{}

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
				wg.Add(1)
				go func(fieldType reflect.StructField, field reflect.Value) {
					if err := validateFunc(fieldType, field); err != nil {
						errorsBag.add(fieldType.Name, err)
					}
					wg.Done()
				}(fieldType, field)
			}
		}
	}
	wg.Wait()

	return errorsBag
}

func extractValidators(validators string) []string {
	return strings.Split(validators, ",")
}
