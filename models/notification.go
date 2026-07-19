package models

type Notification struct {
	Type   string `gorm:"type:char(36);primaryKey" json:"type"`
	Value  string `json:"value"`
	UserId string `gorm:"type:char(36);primaryKey" json:"user_id"`
	IsRead int    `gorm:"default:0" json:"is_read"`
}
