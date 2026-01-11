package models

import (
	"time"
)

type Game struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	GameTime  time.Time `gorm:"index" json:"gameTime"`

	HomeTeamID uint `gorm:"index" json:"homeTeamId"`
	HomeTeam   Team `gorm:"foreignKey:HomeTeamID" json:"homeTeam,omitempty"`

	AwayTeamID uint `gorm:"index" json:"awayTeamId"`
	AwayTeam   Team `gorm:"foreignKey:AwayTeamID" json:"awayTeam,omitempty"`

	HomeTeamPoints uint `json:"homeTeamPoints,omitempty"`
	AwayTeamPoints uint `json:"awayTeamPoints,omitempty"`

	AvgRating   float64 `gorm:"column:avg_rating;->" json:"avgRating"`
	RatingCount int64   `gorm:"column:rating_count;->" json:"ratingCount"`
	Rating      int     `gorm:"column:rating;->" json:"rating"`
}
