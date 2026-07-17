package mod

type Mod struct {
	ID          string `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	Author      string `gorm:"type:char(36);primaryKey" json:"author"`
	Game 	  	string `json:"game"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	Image       string `json:"image"`
	Path        string `json:"path"`
	Status      string `json:"status"`
}
