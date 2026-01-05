package models

type Team struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `json:"teamName"`
	ConferenceID uint       `json:"conferenceId"`
	Conference   Conference `gorm:"foreignKey:ConferenceID" json:"conference,omitempty"`
}

type Conference struct {
	ID   uint   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name string `json:"name"`
}

type TeamStats struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	TeamID uint   `json:"teamId"`
	Team   Team   `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Season string `json:"season"`
}
