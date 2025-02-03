package user

import (
	"context"

	"github.com/mrocha98/go-studies/gobid/internal/validator"
)

type LoginUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	return *eval.CheckIsNotBlank("email", req.Email).
		CheckIsEmail("email", req.Email).
		CheckIsNotBlank("password", req.Password).
		CheckMaxChars("password", req.Password, 120)

}
