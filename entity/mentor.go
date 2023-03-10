package entity

import "gorm.io/gorm"

type Mentor struct {
	gorm.Model
	ProfilePhoto string     `json:"profile_photo" gorm:"default:null"`
	FullName     string     `json:"full_name" gorm:"type:VARCHAR(255);NOT NULL"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(255);NOT NULL"`
	Skill        []Skill    `json:"skillID" gorm:"many2many:mentors_skill"`
	Interest     []Interest `json:"interestID" gorm:"many2many:mentors_interest"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:VARCHAR(255);NOT NULL"`
	Rate         int32      `json:"rate"`
	Fee          int32      `json:"fee"`
	Mentee       []User     `json:"mentee" gorm:"foreignkey:MentorID"`
	Exp          []Exp      `json:"exp"`
}

type MentorReqByID struct {
	ID uint `json:"id"`
}

type MentorRespData struct {
	ProfilePhoto string         `json:"profile_photo"`
	FullName     string         `json:"full_name"`
	Skill        []RespSkill    `json:"skill"`
	Interest     []RespInterest `json:"interest"`
	Lokasi       string         `json:"lokasi"`
	Deskripsi    string         `json:"deskripsi"`
	Rate         int32          `json:"rate"`
	Fee          int32          `json:"fee"`
}

type MentorAdd struct {
	gorm.Model
	ProfilePhoto string `json:"profile_photo"`
	FullName     string `json:"full_name"`
	Lokasi       string `json:"lokasi"`
	Skill        []uint `json:"skillID"`
	Interest     []uint `json:"interestID"`
	Deskripsi    string `json:"deskripsi"`
	Rate         int32  `json:"rate"`
	Fee          int32  `json:"fee"`
	Exp          []Exp  `json:"exp"`
}

type MentorParam struct {
	PostID int64 `uri:"mentor_id" gorm:"column:id"`
	PaginationParam
}

type MentorFilter struct {
	InterestID []uint `json:"Interestid"`
}

type MentorSearch struct {
	Key string `json:"key"`
}
