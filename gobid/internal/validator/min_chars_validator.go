package validator

import (
	"fmt"
	"unicode/utf8"
)

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func (e *Evaluator) CheckMinChars(key, value string, n int) *Evaluator {
	return e.CheckField(
		MinChars(value, n),
		key,
		fmt.Sprintf(MIN_CHARS_MESSAGE, n),
	)
}
