package user

import (
	"fmt"

	"meme-generator/internal/model"
	"meme-generator/internal/utils"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *model.User) (*model.User, error) {
	user.ID = 0

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) All(limit, page int) ([]*model.User, error) {
	var users []*model.User

	if err := r.db.Scopes(utils.NewPaginate(limit, page).PaginatedResult).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (r *Repository) FindByID(id uint) (*model.User, error) {
	var user model.User

	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) CheckExistsByUsername(username string) (bool, error) {
	var user model.User

	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
