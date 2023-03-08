package entity

import "gorm.io/gorm"

// import "gorm.io/gorm"

type Mentor struct {
	gorm.Model
	ProfilePhoto string     `json:"profile_photo" gorm:"default:null"`
	FullName     string     `json:"full_name" gorm:"type:VARCHAR(255);NOT NULL"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(255);NOT NULL"`
	Skill        []Skill    `json:"skill" gorm:"type:NOT NULL;many2many:mentor_skill"`
	Interest     []Interest `json:"interest" gorm:"many2many:mentors_interest"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:VARCHAR(255);NOT NULL"`
	Rate         int32      `json:"rate"`
	Fee          float32    `json:"fee"`
	Mentee       []User     `json:"mentee"`
	Exp          []Exp      `json:"exp"`
}

type MentorReqByID struct {
	ID uint `json:"id"`
}

type MentorRespData struct {
	FullName  string     `json:"full_name"`
	Skill     []Skill    `json:"skill"`
	Lokasi    string     `json:"lokasi"`
	Interest  []Interest `json:"interest"`
	Deskripsi string     `json:"deskripsi"`
	Rate      int32      `json:"rate"`
	Fee       float32    `json:"fee"`
}
