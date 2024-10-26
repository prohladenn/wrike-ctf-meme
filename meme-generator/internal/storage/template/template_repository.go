package template

import (
	"fmt"

	"meme-generator/internal/model"

	"gorm.io/gorm"
)

var ErrNotFound = fmt.Errorf("template not found")

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	// NewRepository returns a new Repository, used to interact with the Template database table.
	return &Repository{
		db: db,
	}
}

func (r *Repository) LastRecords() ([]*model.Template, error) {
	var templates []*model.Template

	if err := r.db.Order("created_at desc").Limit(100).Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

func (r *Repository) Create(template *model.Template) (*model.Template, error) {
	template.ID = 0

	if err := r.db.Create(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (r *Repository) FindByID(id uint) (*model.Template, error) {
	var template model.Template

	if err := r.db.Where("id = ?", id).First(&template).Error; err != nil {
		return nil, err
	}

	return &template, nil
}

func (r *Repository) FindByUserID(userID uint) ([]*model.Template, error) {
	var templates []*model.Template

	if err := r.db.Where("user_id = ?", userID).Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}
