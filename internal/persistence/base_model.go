package persistence

import "time"

type BaseModel struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"                        json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"                        json:"updateAt"`
}
