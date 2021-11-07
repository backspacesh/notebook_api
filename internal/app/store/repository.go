package store

import "rest_api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
}

type ArticleRepository interface {
	CreateArticle(*model.Article) error
	FindByHeading(string) (*model.Article, error)
	ShowAllArticles() ([]*model.Article, error)
	DeleteArticle(int) (string, error)
	ChangeArticleById(*model.Article) error
}
