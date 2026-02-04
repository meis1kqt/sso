package sqlite

import (
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/meis1kqt/sso/internal/domain/models"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		slog.Error("invalid connect to db", "error", err)
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {

	stmt, err := s.db.Prepare("SELECT * from users WHERE email = ?")
	if err != nil {
		return models.User{}, err
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User

	err = row.Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		return models.User{}, err
	}

	return user, nil

}

func (s *Storage) IsAdmin(ctx context.Context, userID int64)(bool, error) {

	stmt, err := s.db.Prepare("SELECT is_admin from users where user_id = ?")

	if err != nil {
		return false, err
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool

	err = row.Scan(&isAdmin)

	if err != nil {
		return false, err
	}

	return isAdmin, nil



}