package models

import (
	"time"
)

type Game struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	GameTime       time.Time `json:"gameTime"`
	HomeTeamID     uint      `json:"homeTeamId"`
	HomeTeam       Team      `gorm:"foreignKey:HomeTeamID" json:"homeTeam,omitempty"`
	AwayTeamID     uint      `json:"awayTeamId"`
	AwayTeam       Team      `gorm:"foreignKey:AwayTeamID" json:"awayTeam,omitempty"`
	HomeTeamPoints uint      `json:"homeTeamPoints,omitempty"`
	AwayTeamPoints uint      `json:"awayTeamPoints,omitempty"`
	Rating         int       `json:"rating,omitempty"`
}
