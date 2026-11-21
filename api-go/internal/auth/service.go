package auth

import (
	"context"
	"errors"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   authRepository
	tokens *TokenManager
}

func NewService(repo authRepository, tokens *TokenManager) *Service {
	return &Service{
		repo:   repo,
		tokens: tokens,
	}
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func VerifyPassword(password, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}

func (s *Service) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	return s.repo.IsTokenRevoked(ctx, jti)
}

var ErrInvalidCredentials = errors.New("invalid email or password")

func (_service *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := _service.repo.FindUserByEmail(ctx, strings.TrimSpace(email))
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	if err := VerifyPassword(password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	return _service.tokens.Issue(strconv.FormatInt(user.ID, 10))
}

var ErrInvalidEmail = errors.New("invalid email format")

func (s *Service) Register(ctx context.Context, email, password string) error {
	email = strings.TrimSpace(email)

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return ErrInvalidEmail
	}

	passwordHash, err := HashPassword(password)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, email, passwordHash)
}

func (_service *Service) RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error {
	return _service.repo.RevokeToken(ctx, jti, expiresAt)
}

var ErrPasswordChangedConcurrently = errors.New("password changed during request")

func (_service *Service) ChangePassword(
	ctx context.Context,
	userID int64,
	oldPassword, newPassword string,
) error {
	currentHash, err := _service.repo.FindPasswordHash(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidCredentials
	}
	if err != nil {
		return err
	}

	if err := VerifyPassword(oldPassword, currentHash); err != nil {
		return ErrInvalidCredentials
	}

	newHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	updated, err := _service.repo.UpdatePasswordHash(ctx, userID, currentHash, newHash)
	if err != nil {
		return err
	}
	if !updated {
		return ErrPasswordChangedConcurrently
	}

	return nil
}
