package sqlstore

import (
	_ "database/sql"
	"rest_api/internal/app/model"
)

type ArticleRepository struct {
	store *Store
}


func (a *ArticleRepository) CreateArticle(ar *model.Article) error {
	return a.store.db.QueryRow(
		"INSERT INTO articles(article_header, article_text, author_id, creating_date) values ($1, $2, $3, now()::DATE) RETURNING id, creating_date",
		&ar.Heading,
		&ar.Text,
		&ar.AuthorID,
	).Scan(
		&ar.ID,
		&ar.Date,
	)
}

func (a *ArticleRepository) FindByHeading(header string) (*model.Article, error) {
	ar := &model.Article{}

	if err := a.store.db.QueryRow(
		"select a.id, a.article_header, a.article_text, u.name, a.creating_date from articles a left join users u on u.id=a.author_id where a.article_header=$1",
		header,
	).Scan(
		&ar.ID,
		&ar.Heading,
		&ar.Text,
		&ar.AuthorName,
		&ar.Date,
	); err != nil {
		return nil, err
	}

	return ar, nil
}

func (a *ArticleRepository) ShowAllArticles() ([]*model.Article, error) {
	ars := make([]*model.Article, 0)

	rows, err := a.store.db.Query(
	"select a.id, a.article_header, a.article_text, u.name, a.creating_date from articles a left join users u on u.id=a.author_id")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		ar := &model.Article{}
		if err := rows.Scan(
			&ar.ID,
			&ar.Heading,
			&ar.Text,
			&ar.AuthorName,
			&ar.Date,
		); err != nil {
			return nil, err
		}

		ars = append(ars, ar)
	}
	return ars, nil
}

func (a *ArticleRepository) DeleteArticle(id int) (string, error) {
	var articleHeader string

	if err := a.store.db.QueryRow(
		"DELETE FROM articles where id=$1 returning article_header",
		id,
	).Scan(
		&articleHeader,
	); err != nil {
		return "", err
	}

	return articleHeader, nil
}
func (a *ArticleRepository) ChangeArticleById(ar *model.Article) error {
	return a.store.db.QueryRow(
		"Update articles set article_header=$1, article_text=$2 where id=$3 returning article_header, article_text",
		ar.Heading,
		ar.Text,
		ar.ID,
	).Scan(
		&ar.Heading,
		&ar.Text,
	)
}

//select a.article_header, a.article_text, u.name, a.creating_date from articles a left join users u on u.id=a.author_id;
//select a.article_header, a.article_text, u.name, a.creating_date from articles a left join users u on u.id=a.author_id where a.article_header='Hello';
