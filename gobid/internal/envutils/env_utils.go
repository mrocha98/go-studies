package envutils

import "os"

type Env interface {
	APIHost() string
	APIPort() string
	DBHost() string
	DBPort() string
	DBName() string
	DBUser() string
	DBPassword() string
	PasswordPepper() string
}

type OSEnv struct{}

func NewOSEnv() Env {
	return OSEnv{}
}

func (env OSEnv) APIHost() string {
	return os.Getenv("GOBID_API_HOST")
}

func (env OSEnv) APIPort() string {
	return os.Getenv("GOBID_API_PORT")
}

func (env OSEnv) DBHost() string {
	return os.Getenv("GOBID_DATABASE_HOST")
}

func (env OSEnv) DBPort() string {
	return os.Getenv("GOBID_DATABASE_PORT")
}

func (env OSEnv) DBName() string {
	return os.Getenv("GOBID_DATABASE_NAME")
}

func (env OSEnv) DBUser() string {
	return os.Getenv("GOBID_DATABASE_USER")
}

func (env OSEnv) DBPassword() string {
	return os.Getenv("GOBID_DATABASE_PASSWORD")
}

func (env OSEnv) PasswordPepper() string {
	return os.Getenv("GOBID_PASSWORD_PEPPER")
}
