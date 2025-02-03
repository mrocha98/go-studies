package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

var (
	ErrDuplicatedUserNameOrEmail = errors.New("userName or email already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
)

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

func (us *UserService) AuthenticateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := us.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, ErrInvalidCredentials
		}
		return uuid.UUID{}, err
	}

	if err := us.passwordHasher.Compare(
		user.PasswordHash, []byte(password), user.PasswordSalt); err != nil {
		if errors.Is(err, cryptoutils.ErrMismatchedHashAndPassword) {
			return uuid.UUID{}, ErrInvalidCredentials
		}
		return uuid.UUID{}, err
	}

	return user.ID, nil
}
