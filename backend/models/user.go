package models

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserName    string    `json:"userName"`
	Password    string    `json:"-"`
	Email       string    `json:"email"`
	RoleID      uint      `json:"roleId"`
	Role        Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	LastLoginAt time.Time `json:"lastLoginAt"`
	HideScores  bool      `json:"hideScores" gorm:"column:hide_scores"`
}

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}
