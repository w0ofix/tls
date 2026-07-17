package models

type Ips struct {
	RegisteredIp string   `json:"registered_ip"`
	LoginIps     []string `json:"login_ips"`
}

type User struct {
	BaseModel
	Email    string `gorm:"unique;not null" json:"email"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"default:user" json:"role"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
	Ips      Ips    `gorm:"serializer:json" json:"-"`
}
