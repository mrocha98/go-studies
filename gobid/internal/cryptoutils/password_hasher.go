package cryptoutils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/mrocha98/go-studies/gobid/internal/envutils"
	"golang.org/x/crypto/argon2"
)

type PasswordHasher interface {
	// Generates a secure hash of a password.
	//
	// The function takes a password string and returns:
	//   - A randomly generated salt
	//   - The generated password hash
	//   - Any error that occurred during the process
	Hash(password string) ([]byte, []byte, error)
	Compare(hashed, password, salt []byte) error
}

var ErrMismatchedHashAndPassword = errors.New(
	"PasswordHasher: hashedPassword is not the hash of the given password and salt",
)

type Argon2PasswordHasher struct {
	timeCost   uint32
	memoryCost uint32
	threads    uint8
	keyLen     uint32
	env        envutils.Env
}

func NewArgon2PasswordHasher() PasswordHasher {
	return Argon2PasswordHasher{
		timeCost:   uint32(1),         // Number of iterations
		memoryCost: uint32(64 * 1024), // 64 MB memory usage
		threads:    uint8(4),          // Degree of parallelism
		keyLen:     uint32(32),        // Desired length of the hash
		env:        envutils.NewOSEnv(),
	}
}

func (ph Argon2PasswordHasher) Hash(password string) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, nil, err
	}

	pepper := ph.env.PasswordPepper()
	passwordWithPepper := []byte(password + pepper)

	hash := argon2.IDKey(
		passwordWithPepper, salt,
		ph.timeCost, ph.memoryCost, ph.threads, ph.keyLen,
	)

	return salt, hash, nil
}

func (ph Argon2PasswordHasher) Compare(hashed, password, salt []byte) error {
	pepper := ph.env.PasswordPepper()
	passwordWithPepper := []byte(hex.EncodeToString(password) + pepper)

	hash := argon2.IDKey(
		passwordWithPepper, salt,
		ph.timeCost, ph.memoryCost, ph.threads, ph.keyLen,
	)

	if hex.EncodeToString(hash) != hex.EncodeToString(hashed) {
		return ErrMismatchedHashAndPassword
	}
	return nil
}
