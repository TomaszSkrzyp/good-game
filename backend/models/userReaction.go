package models

import "time"

type UserReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GameID    uint      `json:"gameId"`
	Game      Game      `gorm:"foreignKey:GameID" json:"game,omitempty"`
	UserID    uint      `json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Liked     int       `json:"liked"`
	CreatedAt time.Time `json:"createdAt"`
}
