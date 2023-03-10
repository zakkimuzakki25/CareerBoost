package entity

import "gorm.io/gorm"

// import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Intro    string     `json:"intro" gorm:"type:VARCHAR(255)"`
	Playlist []Playlist `json:"playlist" gorm:"foreignKey:CourseID"`
	User     []User     `gorm:"many2many:users_course"`
	Star     uint       `json:"star"`
}
