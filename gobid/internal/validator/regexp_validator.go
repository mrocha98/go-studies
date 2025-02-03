package validator

import "regexp"

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (e *Evaluator) CheckMatches(key, value string, rx *regexp.Regexp, customMessage string) *Evaluator {
	message := customMessage
	if message == "" {
		message = MATCHES_MESSAGE
	}
	return e.CheckField(Matches(value, rx), key, message)
}
