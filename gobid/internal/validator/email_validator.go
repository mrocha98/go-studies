package validator

import "regexp"

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmail(value string) bool {
	return Matches(value, EmailRX)
}

func (e *Evaluator) CheckIsEmail(key, value string) *Evaluator {
	return e.CheckMatches(key, value, EmailRX, IS_EMAIL_MESSAGE)
}
