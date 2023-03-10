package entity

import "gorm.io/gorm"

type Exp struct {
	gorm.Model
	Logo       string   `json:"logo"`
	Skill      string   `json:"skill"`
	Perusahaan string   `json:"perusahaan"`
	Mentor     []Mentor `json:"mentor" gorm:"foreignkey:MentorID"`
	MentorID   uint
}
