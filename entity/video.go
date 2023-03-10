package entity

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Link       string `json:"link" gorm:"type:varchar(255)"`
	Judul      string `json:"judul" gorm:"type:varchar(255)"`
	Durasi     string `json:"durasi" gorm:"type:varchar(255)"`
	PlaylistID uint
}
