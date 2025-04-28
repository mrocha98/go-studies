package validator

import "strings"

func GreaterThan[T number](value T, toCompare T) bool {
	return float64(value) > float64(toCompare)
}

func (e *Evaluator) CheckIsGreaterThanInt(key string, value, toCompare int64) *Evaluator {
	return e.CheckField(GreaterThanOrEqual(value, toCompare), key, GREATER)
}

func (e *Evaluator) CheckIsGreaterThanFloat(key string, value, toCompare float64) *Evaluator {
	return e.CheckField(GreaterThanOrEqual(value, toCompare), key, strings.ReplaceAll(GREATER, "%d", "%f"))
}
