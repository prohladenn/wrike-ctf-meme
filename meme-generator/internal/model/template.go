package model

import "gorm.io/gorm"

type Template struct {
	gorm.Model

	Name     string
	DirPath  string
	FileName string

	UserID uint
}

type TemplateDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	OwnerID   uint   `json:"owner_id"`
}

func (m *Template) ToDTO() *TemplateDTO {
	return &TemplateDTO{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt.String(),
		OwnerID:   m.UserID,
	}
}

func TemplatesToDTOs(templates []*Template) []*TemplateDTO {
	dtos := make([]*TemplateDTO, 0, len(templates))
	for _, template := range templates {
		dtos = append(dtos, template.ToDTO())
	}
	return dtos
}

type TemplatePrivateDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	PrivateInfo string `json:"private_info"`
	OwnerID     uint   `json:"owner_id"`
}

func (m *Template) ToPrivateDTO(privateInfo string) *TemplatePrivateDTO {
	return &TemplatePrivateDTO{
		ID:          m.ID,
		Name:        m.Name,
		CreatedAt:   m.CreatedAt.String(),
		PrivateInfo: privateInfo,
		OwnerID:     m.UserID,
	}
}
