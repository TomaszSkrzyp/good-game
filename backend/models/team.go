package models

type Team struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"uniqueIndex" json:"teamName"`
	Abbreviation string     `gorm:"uniqueIndex" json:"abbreviation"`
	ConferenceID uint       `json:"conferenceId"`
	Conference   Conference `gorm:"foreignKey:ConferenceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"conference,omitempty"`
}

type Conference struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"uniqueIndex" json:"name"`
}

type TeamStats struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	TeamID uint   `gorm:"index" json:"teamId"`
	Team   Team   `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Season string `gorm:"index" json:"season"`
}
