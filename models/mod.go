package models

type Mod struct {
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	Author      string `gorm:"type:char(36);primaryKey" json:"author"`
	Game        string `json:"game"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	Path        string `json:"path"`
	Status      int    `gorm:"default:0" json:"status"`
}
