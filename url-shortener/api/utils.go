package api

import "math/rand/v2"

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func GenerateRandomString() string {
	const n = 8
	charsetLen := len(charset)
	bytes := make([]byte, n)
	for i := range n {
		bytes[i] = charset[rand.IntN(charsetLen)]
	}
	return string(bytes)
}
