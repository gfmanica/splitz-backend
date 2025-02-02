package user

import (
	"context"
	"fmt"

	"github.com/gfmanica/splitz-backend/types"
	"github.com/jackc/pgx/v5"
)

type Store struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM users WHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	u := &types.User{}

	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM users WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	u := &types.User{}

	for rows.Next() {
		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func (s *Store) CreateUser(u types.User) error {
	_, err := s.db.Exec(context.Background(), "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, u.Password)

	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoUser(rows pgx.Rows) (*types.User, error) {
	u := &types.User{}

	err := rows.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return u, nil
}
