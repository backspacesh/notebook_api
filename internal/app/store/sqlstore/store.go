package sqlstore

import (
	"database/sql"
	"rest_api/internal/app/store"
)

type Store struct {
	db *sql.DB
	userRepository *UserRepository
	articleRepository *ArticleRepository
}

func New(db *sql.DB) *Store {
	return &Store {
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{
			s,
		}
	}

	return s.userRepository
}

func (s *Store) Article() store.ArticleRepository {
	if s.articleRepository == nil {
		s.articleRepository = &ArticleRepository{
			s,
		}
	}

	return s.articleRepository
}
