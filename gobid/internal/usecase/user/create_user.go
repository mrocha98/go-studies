package user

import (
	"context"

	"github.com/mrocha98/go-studies/gobid/internal/validator"
)

type CreateUserReq struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator

	eval.CheckIsNotBlank("userName", req.UserName).
		CheckMaxChars("userName", req.UserName, 50).
		CheckIsNotBlank("email", req.Email).
		CheckIsEmail("email", req.Email).
		CheckIsNotBlank("bio", req.Bio).
		CheckMinChars("bio", req.Bio, 10).
		CheckMaxChars("bio", req.Bio, 255).
		CheckIsNotBlank("password", req.Password).
		CheckMinChars("password", req.Password, 8).
		CheckMaxChars("password", req.Password, 64)

	return eval
}
