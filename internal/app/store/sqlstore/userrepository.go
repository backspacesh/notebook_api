package sqlstore

import (
	"database/sql"
	"rest_api/internal/app/model"
	"rest_api/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) Create(u *model.User) error {
	err := u.Validate()
	if err != nil {
		return err
	}

	err = u.BeforeCreate()
	if err != nil {
		return err
	}

	return ur.store.db.QueryRow(
		"INSERT INTO users (name, email, encrypted_password) VALUES ($1, $2, $3) RETURNING id",
		u.Name,
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := model.User{}

	if err := ur.store.db.QueryRow(
		"SELECT id, email, encrypted_password FROM users where email = $1",
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return &u, nil
}
