package model

type Article struct {
	ID int `json:"id"`
	Heading string `json:"article_heading"`
	Text string `json:"article_text"`
	Date string `json:"creating_date"`
	AuthorID int `json:"author_id,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
}