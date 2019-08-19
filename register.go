package validator

import (
	"reflect"
)

type ValidatorHandler func(reflect.StructField, reflect.Value) error

var validators map[string]ValidatorHandler

func init() {
	validators = make(map[string]ValidatorHandler)
}

// RegisterValidator maps a validator by name to a corrosponding handler.
// These handlers should not mutate any global data as these handlers run in parallel.
func RegisterValidator(name string, handler ValidatorHandler) {
	validators[name] = handler
}

func unregisterValidator(name string) {
	delete(validators, name)
}
