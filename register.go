package validator

import (
	"reflect"
)

type validatorHandler func(reflect.StructField, reflect.Value) error

var validators map[string]validatorHandler

func init() {
	validators = make(map[string]validatorHandler)
}

// RegisterValidator maps a validator by name to a corrosponding handler
func RegisterValidator(name string, handler validatorHandler) {
	validators[name] = handler
}

func unregisterValidator(name string) {
	delete(validators, name)
}
