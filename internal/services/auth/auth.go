package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/meis1kqt/sso/internal/domain/models"
	"github.com/meis1kqt/sso/internal/lib/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log      *slog.Logger
	storage  Storage
	tokenTTL time.Duration
}

type Storage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (userID int64, err error)
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(log *slog.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int64) (string, error) {

	user, err := a.storage.User(ctx, email)

	if err != nil {
		slog.Error("failed to get User", "error", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		slog.Error("invalid password", "error", err)
		return "", err
	}

	app, err := a.storage.App(ctx, appID)
	if err != nil {
		slog.Error("invalid app", "error", err)
		return "", err
	}

	slog.Info("user is logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)

	if err != nil {
		slog.Error("token error", "error", err)
		return "", err
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	slog.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		slog.Error("password hash error", "error", err)
		return 0, err
	}

	id, err := a.storage.SaveUser(ctx, email, passHash)

	if err != nil {
		slog.Error("saveuser", "error", err)
		return 0, err
	}

	slog.Info("user registered")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	isAdmin, err := a.storage.IsAdmin(ctx, userID)

	if err != nil {
		slog.Error("bag", "error", err)
		return false, err
	}

	return isAdmin, nil
}
