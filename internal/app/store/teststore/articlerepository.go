package teststore

import (
	"errors"
	"rest_api/internal/app/model"
	"rest_api/internal/app/store"
	"time"
)

type ArticleRepository struct {
	store *Store
}

func (ar *ArticleRepository) CreateArticle(article *model.Article) error {
	article.Date = time.Now().String()
	ar.store.articles = append(ar.store.articles, article)

	article.ID = -1

	for key, value := range ar.store.articles {
		if value.Heading == article.Heading {
			article.ID = key
		}
	}

	if article.ID == -1 {
		return store.ErrCreate
	}

	return nil
}

func (ar *ArticleRepository) FindByHeading(header string) (*model.Article, error) {
	for _, value := range ar.store.articles {
		if value.Heading == header {
			return value, nil
		}
	}

	return nil, errors.New("article not found")
}

func (ar *ArticleRepository) ShowAllArticles() ([]*model.Article, error) {
	ars := make([]*model.Article, 0)

	for _, value := range ar.store.articles {
		ars = append(ars, value)
	}

	return ars, nil
}

func (ar *ArticleRepository) DeleteArticle(id int) (string, error) {
	if len(ar.store.articles) < id + 1 {
		return "", errors.New("article not found")
	}

	header := ar.store.articles[id].Heading
	copy(ar.store.articles[id:], ar.store.articles[id + 1:])
	ar.store.articles[len(ar.store.articles) - 1] = nil
	ar.store.articles = ar.store.articles[:len(ar.store.articles) - 1]

	return header, nil
}

func (ar *ArticleRepository) ChangeArticleById(article *model.Article) error {
	if len(ar.store.articles) < article.ID + 1 {
		return errors.New("article not found")
	}

	ar.store.articles[article.ID].Heading = article.Heading
	ar.store.articles[article.ID].Text = article.Text

	if ar.store.articles[article.ID].Heading != article.Heading ||
		ar.store.articles[article.ID].Text != article.Text {
		return errors.New("change went with wrong")
	}

	return nil
}