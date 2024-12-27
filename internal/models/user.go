package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email     string    `gorm:"type:varchar(255);" json:"email"`
	Name      string    `gorm:"type:varchar(255);" json:"name"`
	Forms     []Form    `gorm:"foreignKey:UserID" json:"forms"`
	PubKey    string    `gorm:"type:varchar(1024);not null;unique" json:"pub_key"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type UserRequest struct {
	Email  string `json:"email"`
	Name   string `json:"name"`
	PubKey string `json:"pub_key"`
}

type UserResponse struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Name   string    `json:"name"`
	Forms  []Form    `json:"forms"`
	PubKey string    `json:"pub_key"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:     u.ID,
		Email:  u.Email,
		Name:   u.Name,
		Forms:  u.Forms,
		PubKey: u.PubKey,
	}
}

func (u *UserRequest) ToUser() User {
	return User{
		Email:  u.Email,
		Name:   u.Name,
		PubKey: u.PubKey,
	}
}
