package model

import "testing"

func TestArticle(t *testing.T, id int) *Article {
	t.Helper()

	return &Article {
		Heading: "TestArticle",
		Text: "TestArticle",
		AuthorID: id,
	}
}
