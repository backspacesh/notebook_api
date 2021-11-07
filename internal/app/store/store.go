package store

type Store interface {
	User() UserRepository
	Article() ArticleRepository
}
