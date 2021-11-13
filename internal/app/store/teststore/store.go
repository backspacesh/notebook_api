package teststore

import (
	"rest_api/internal/app/model"
	"rest_api/internal/app/store"
)

type Store struct {
	users []*model.User
	articles []*model.Article
	userRepository *UserRepository
	articleRepository *ArticleRepository
}

func New() *Store {
	return &Store{
		users: make([]*model.User, 0),
		articles: make([]*model.Article, 0),
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository == nil {
		s.userRepository = &UserRepository{s}
	}

	return s.userRepository
}

func (s *Store) Article() store.ArticleRepository {
	if s.articleRepository == nil {
		s.articleRepository = &ArticleRepository{s}
	}

	return s.articleRepository
}