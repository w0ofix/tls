package models

type ModReport struct {
	ModId      string `gorm:"type:char(36);primaryKey" json:"mod_id"`
	ReportedBy string `gorm:"type:char(36);primaryKey" json:"reported_by"`
	Category   string `gorm:"type:char(36);primaryKey" json:"category"`
	Message    string `json:"message"`
}
