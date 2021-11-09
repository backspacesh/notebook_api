package sqlstore_test

import (
	"github.com/stretchr/testify/assert"
	"rest_api/internal/app/model"
	"rest_api/internal/app/store/sqlstore"
	"testing"
)

func TestArticleRepository_CreateArticle(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseString)
	defer teardown("users", "articles")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	a := model.TestArticle(t, u.ID)
	err := s.Article().CreateArticle(a)
	assert.NoError(t, err)
	assert.NotNil(t, a)
}

func TestArticleRepository_FindByHeading(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseString)
	defer teardown("users", "articles")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	a := model.TestArticle(t, u.ID)
	s.Article().CreateArticle(a)

	a1, err := s.Article().FindByHeading(a.Heading)
	assert.NoError(t, err)
	assert.NotNil(t, a1)
}

func TestArticleRepository_DeleteArticle(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseString)
	defer teardown("users", "articles")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	a := model.TestArticle(t, u.ID)
	err := s.Article().CreateArticle(a)
	assert.NoError(t, err)

	header, err := s.Article().DeleteArticle(a.ID)
	assert.NoError(t, err)
	assert.Equal(t, a.Heading, header)
}

func TestArticleRepository_ChangeArticleById(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseString)
	defer teardown("users", "articles")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	s.User().Create(u)

	a := model.TestArticle(t, u.ID)
	s.Article().CreateArticle(a)

	another := & model.Article{
		Heading: "Another Header",
		Text: "Another text",
		ID: a.ID,
	}
	err := s.Article().ChangeArticleById(another)
	assert.NoError(t, err)

	a, err = s.Article().FindByHeading("Another Header")
	assert.Equal(t, "Another Header", a.Heading)
	assert.Equal(t, "Another text", a.Text)
}
