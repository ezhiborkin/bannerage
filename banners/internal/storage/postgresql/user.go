package postgresql

import (
	"banners/domain/models"
	"banners/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) GetUserStorage(email string) (*models.User, error) {
	const op = "storage.postgresql.GetUserStorage"

	query, args, err := sq.Select("id", "email", "role", "password_hash").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := s.db.QueryRow(query, args...)

	user := &models.User{}
	err = row.Scan(&user.ID, &user.Email, &user.Role, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) CreateUserStorage(email, role string, passHash []byte) error {
	const op = "storage.postgresql.CreateUserStorage"

	query, args, err := sq.Insert("users").
		Columns("email", "role", "password_hash").
		Values(email, role, passHash).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
