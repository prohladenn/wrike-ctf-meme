package model

import "gorm.io/gorm"

type Meme struct {
	gorm.Model

	Name string

	DirPath  string
	FileName string

	UserID         uint
	AuthorUsername string `gorm:"-"`
	MemeTemplateID uint
}

type MemeDTO struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CreatedAt     string `json:"created_at"`
	OwnerID       uint   `json:"owner_id"`
	OwnerUsername string `json:"owner_username"`
	TemplateID    uint   `json:"template_id"`
}

func (m *Meme) ToDTO() *MemeDTO {
	return &MemeDTO{
		ID:            m.ID,
		Name:          m.Name,
		CreatedAt:     m.CreatedAt.String(),
		OwnerID:       m.UserID,
		OwnerUsername: m.AuthorUsername,
		TemplateID:    m.MemeTemplateID,
	}
}

func MemesToDTOs(memes []*Meme) []*MemeDTO {
	dtos := make([]*MemeDTO, 0, len(memes))
	for _, meme := range memes {
		dtos = append(dtos, meme.ToDTO())
	}
	return dtos
}
