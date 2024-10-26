package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `gorm:"uniqueIndex;not null"`
	Password string
}

type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:       u.ID,
		Username: u.Username,
	}
}

func UsersToDTOs(users []*User) []*UserDTO {
	dtos := make([]*UserDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, user.ToDTO())
	}
	return dtos
}
