package validator

import "context"

type Validator interface {
	Valid(context.Context) Evaluator
}

type Evaluator map[string]any

func (e *Evaluator) addFieldError(key, message string) {
	if *e == nil {
		*e = make(Evaluator)
	}

	if _, exists := (*e)[key]; !exists {
		(*e)[key] = message
	}
}

func (e *Evaluator) CheckField(ok bool, key, message string) *Evaluator {
	if !ok {
		e.addFieldError(key, message)
	}
	return e
}
