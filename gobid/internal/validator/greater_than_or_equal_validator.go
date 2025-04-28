package validator

import "strings"

func GreaterThanOrEqual[T number](value T, toCompare T) bool {
	return float64(value) >= float64(toCompare)
}

func (e *Evaluator) CheckIsGreaterThanOrEqualInt(key string, value, toCompare int64) *Evaluator {
	return e.CheckField(GreaterThanOrEqual(value, toCompare), key, GREATER_OR_EQUAL)
}

func (e *Evaluator) CheckIsGreaterThanOrEqualFloat(key string, value, toCompare float64) *Evaluator {
	return e.CheckField(GreaterThanOrEqual(value, toCompare), key, strings.ReplaceAll(GREATER_OR_EQUAL, "%d", "%f"))
}
