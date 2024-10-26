package storage

import "meme-generator/internal/model"

type UserRepository interface {
	Create(user *model.User) (*model.User, error)
	All(limit, page int) ([]*model.User, error)
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	CheckExistsByUsername(username string) (bool, error)
}

type TemplateRepository interface {
	LastRecords() ([]*model.Template, error)
	Create(template *model.Template) (*model.Template, error)
	FindByID(id uint) (*model.Template, error)
	FindByUserID(userID uint) ([]*model.Template, error)
}

type MemeRepository interface {
	LastRecords() ([]*model.Meme, error)
	Create(meme *model.Meme) (*model.Meme, error)
	All(limit, page int) ([]model.Meme, error)
	FindByID(id uint) (*model.Meme, error)
	FindByUserID(userID uint) ([]*model.Meme, error)
}
