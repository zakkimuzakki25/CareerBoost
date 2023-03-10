package entity

import "gorm.io/gorm"

type Playlist struct {
	gorm.Model
	Nama  string  `gorm:"type:VARCHAR(55)"`
	Video []Video `json:"video" gorm:"foreignKey:PlaylistID"`
	// Course   Course  `json:"course" gorm:"foreignKey:CourseID"`
	CourseID uint
}
