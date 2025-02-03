package validator

import (
	"fmt"
	"unicode/utf8"
)

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func (e *Evaluator) CheckMaxChars(key, value string, n int) *Evaluator {
	return e.CheckField(
		MaxChars(value, n),
		key,
		fmt.Sprintf(MAX_CHARS_MESSAGE, n),
	)
}
