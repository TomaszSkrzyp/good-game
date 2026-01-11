package models

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoleID      uint      `json:"roleId"`
	Role        Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	LastLoginAt time.Time `json:"lastLoginAt"`
	UserName    string    `json:"userName"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
}

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
