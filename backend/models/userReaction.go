package models

import "time"

type UserReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GameID    uint      `gorm:"uniqueIndex:idx_user_game" json:"gameId"`
	Game      Game      `gorm:"foreignKey:GameID" json:"game,omitempty"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_game" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Liked     int       `gorm:"column:liked" json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
}
