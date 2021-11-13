package teststore

import (
	"rest_api/internal/app/model"
	"rest_api/internal/app/store"
)

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) Create(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeCreate(); err != nil {
		return err
	}

	ur.store.users = append(ur.store.users, user)

	user.ID = -1

	for key, value := range ur.store.users {
		if value.Email == user.Email {
			user.ID = key
		}
	}

	if user.ID == -1 {
		return store.ErrCreate
	}

	return nil
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, value := range ur.store.users {
		if value.Email == email {
			return value, nil
		}
	}

	return nil, store.ErrRecordNotFound
}