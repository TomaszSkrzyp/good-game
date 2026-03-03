package models

import (
	"time"
)

type Game struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ESPNID    string    `gorm:"column:espn_id;uniqueIndex" json:"espnId"`
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

	GameQuality GameQuality `gorm:"embedded" json:"gameQuality"`
	Status      string      `json:"status"`
}

type GameQuality struct {
	QualityScore uint `gorm:"column:quality_score;default:0" json:"qualityScore"`

	IsBigScoring bool `gorm:"column:is_big_scoring;default:false" json:"isBigScoring"`
	IsBigGame    bool `gorm:"column:is_big_game;default:false" json:"isBigGame"`
	IsClutch     bool `gorm:"column:is_clutch;default:false" json:"isClutch"`
	IsStarDuel   bool `gorm:"column:is_star_duel;default:false" json:"isStarDuel"`
	IsHugeSwing  bool `gorm:"column:is_huge_swing;default:false" json:"isHugeSwing"`
	IsShootout   bool `gorm:"column:is_shootout;default:false" json:"isShootout"`
	IsGritty     bool `gorm:"column:is_gritty;default:false" json:"isGritty"`
}
