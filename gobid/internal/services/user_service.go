package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrocha98/go-studies/gobid/internal/cryptoutils"
	"github.com/mrocha98/go-studies/gobid/internal/store/pgstore"
)

type UserService struct {
	pool           *pgxpool.Pool
	queries        *pgstore.Queries
	passwordHasher cryptoutils.PasswordHasher
}

var ErrDuplicatedUserNameOrEmail = errors.New("userName or email already exists")

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:           pool,
		queries:        pgstore.New(pool),
		passwordHasher: cryptoutils.NewArgon2PasswordHasher(),
	}
}

func (us *UserService) CreateUser(
	ctx context.Context,
	userName, email, password, bio string,
) (uuid.UUID, error) {
	salt, hash, err := us.passwordHasher.Hash(password)
	if err != nil {
		return uuid.UUID{}, nil
	}

	args := pgstore.CreateUserParams{
		UserName:     userName,
		Email:        email,
		PasswordHash: hash,
		PasswordSalt: salt,
		Bio:          bio,
	}
	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == _UNIQUE_VIOLATION {
			return uuid.UUID{}, ErrDuplicatedUserNameOrEmail
		}
		return uuid.UUID{}, err
	}

	return id, nil
}
