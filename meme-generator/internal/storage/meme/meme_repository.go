package meme

import (
	"fmt"

	"meme-generator/internal/model"
	"meme-generator/internal/utils"

	"gorm.io/gorm"
)

var ErrNotFound = fmt.Errorf("meme not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) LastRecords() ([]*model.Meme, error) {
	var memes []*model.Meme

	if err := r.db.Order("created_at desc").Limit(100).Find(&memes).Error; err != nil {
		return nil, err
	}

	for i := range memes {
		authorID := memes[i].UserID
		user := model.User{}

		if err := r.db.Where("id = ?", authorID).First(&user).Error; err != nil {
			return nil, err
		}

		memes[i].AuthorUsername = user.Username
	}

	return memes, nil
}

func (r *Repository) Create(meme *model.Meme) (*model.Meme, error) {
	meme.ID = 0

	if err := r.db.Create(meme).Error; err != nil {
		return nil, err
	}

	return meme, nil
}

func (r *Repository) All(limit, page int) ([]model.Meme, error) {
	var memes []model.Meme

	if err := r.db.Scopes(utils.NewPaginate(limit, page).PaginatedResult).Find(&memes).Error; err != nil {
		return nil, fmt.Errorf("failed to get memes: %w", err)
	}

	for i := range memes {
		authorID := memes[i].UserID
		user := model.User{}

		if err := r.db.Where("id = ?", authorID).First(&user).Error; err != nil {
			return nil, err
		}

		memes[i].AuthorUsername = user.Username
	}

	return memes, nil
}

func (r *Repository) FindByID(id uint) (*model.Meme, error) {
	var meme model.Meme

	if err := r.db.Where("id = ?", id).First(&meme).Error; err != nil {
		return nil, err
	}

	authorID := meme.UserID
	user := model.User{}

	if err := r.db.Where("id = ?", authorID).First(&user).Error; err != nil {
		return nil, err
	}

	meme.AuthorUsername = user.Username

	return &meme, nil
}

func (r *Repository) FindByUserID(userID uint) ([]*model.Meme, error) {
	var memes []*model.Meme

	if err := r.db.Where("user_id = ?", userID).Find(&memes).Error; err != nil {
		return nil, err
	}

	for i := range memes {
		authorID := memes[i].UserID
		user := model.User{}

		if err := r.db.Where("id = ?", authorID).First(&user).Error; err != nil {
			return nil, err
		}

		memes[i].AuthorUsername = user.Username
	}

	return memes, nil
}
