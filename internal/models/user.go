package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email     string    `gorm:"type:varchar(255);"`
	Name      string    `gorm:"type:varchar(255);"`
	Forms     []Form    `gorm:"foreignKey:UserID"`
	PubKey    string    `gorm:"type:varchar(1024);not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type UserRequest struct {
	Email  string
	Name   string
	PubKey string
}

func (u *UserRequest) ToUser() User {
	return User{
		Email:  u.Email,
		Name:   u.Name,
		PubKey: u.PubKey,
	}
}
