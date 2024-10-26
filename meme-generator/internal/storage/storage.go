package storage

type Store interface {
	Users() UserRepository
	Templates() TemplateRepository
	Memes() MemeRepository
}
