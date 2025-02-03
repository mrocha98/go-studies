package validator

import "strings"

func IsNotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func (e *Evaluator) CheckIsNotBlank(key, value string) *Evaluator {
	return e.CheckField(IsNotBlank(value), key, IS_NOT_BLANK_MESSAGE)
}
