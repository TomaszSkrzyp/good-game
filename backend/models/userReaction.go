package models

import "time"

type UserReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GameID    uint      `gorm:"uniqueIndex:idx_user_game" json:"gameId"`
	Game      Game      `gorm:"foreignKey:GameID" json:"game,omitempty"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_game" json:"userId"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Rating    int       `gorm:"column:rating" json:"rating"`
	CreatedAt time.Time `json:"createdAt"`
}
