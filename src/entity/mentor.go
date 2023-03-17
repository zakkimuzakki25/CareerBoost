package entity

import "gorm.io/gorm"

type Mentor struct {
	gorm.Model
	ProfilePhoto string     `json:"profile_photo" gorm:"default:null"`
	FullName     string     `json:"full_name" gorm:"type:VARCHAR(255);NOT NULL"`
	Work         string     `json:"work" gorm:"type:VARCHAR(255);NOT NULL"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(255);NOT NULL"`
	Skill        []Skill    `json:"skillID" gorm:"many2many:mentors_skill"`
	Interest     []Interest `json:"interestID" gorm:"many2many:mentors_interest"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:text;NOT NULL"`
	Rate         float32    `json:"rate"`
	Fee          float32    `json:"fee"`
	Exp          []Exp      `json:"exp"`
	Perusahaan   string     `json:"perusahaan" gorm:"type:text"`
	WA           string     `json:"wa"`
	IG           string     `json:"ig"`
	Email        string     `json:"email"`
	User         []User     `json:"users" gorm:"many2many:user_mentors;ForeignKey:ID"`
}

type MentorReqByID struct {
	ID uint `json:"id"`
}

type MentorRespData struct {
	ID           uint           `json:"id"`
	ProfilePhoto string         `json:"profile_photo"`
	FullName     string         `json:"full_name"`
	Skill        []RespSkill    `json:"skill"`
	Interest     []RespInterest `json:"interest"`
	Lokasi       string         `json:"lokasi"`
	Deskripsi    string         `json:"deskripsi"`
	Rate         float32        `json:"rate"`
	Fee          float32        `json:"fee"`
	Exp          []ExpResp      `json:"exp"`
	Perusahaan   string         `json:"perusahaan"`
}

type MentorRespDataLangganan struct {
	ID           uint           `json:"id"`
	ProfilePhoto string         `json:"profile_photo"`
	FullName     string         `json:"full_name"`
	Skill        []RespSkill    `json:"skill"`
	Interest     []RespInterest `json:"interest"`
	Lokasi       string         `json:"lokasi"`
	Deskripsi    string         `json:"deskripsi"`
	Rate         float32        `json:"rate"`
	Fee          float32        `json:"fee"`
	WA           string         `json:"wa"`
	IG           string         `json:"ig"`
	Email        string         `json:"email"`
}

type MentorAdd struct {
	gorm.Model
	ProfilePhoto string  `json:"profile_photo"`
	FullName     string  `json:"full_name"`
	Work         string  `json:"work"`
	Perusahaan   string  `json:"perusahaan"`
	Lokasi       string  `json:"lokasi"`
	Skill        []uint  `json:"skillID"`
	Interest     []uint  `json:"interestID"`
	Deskripsi    string  `json:"deskripsi"`
	Rate         float32 `json:"rate"`
	Fee          float32 `json:"fee"`
	Exp          []Exp   `json:"exp"`
	WA           string  `json:"wa"`
	IG           string  `json:"ig"`
	Email        string  `json:"email"`
}

type MentorParam struct {
	ID uint `uri:"mentor_id" gorm:"column:id"`
	PaginationParam
}

type Exp struct {
	gorm.Model
	Logo       string `json:"logo"`
	Skill      string `json:"skill"`
	Perusahaan string `json:"perusahaan"`
	MentorID   uint
}

type ExpResp struct {
	Logo       string `json:"logo"`
	Skill      string `json:"skill"`
	Perusahaan string `json:"perusahaan"`
}

type MentorRespHome struct {
	ID           uint           `json:"id"`
	ProfilePhoto string         `json:"profile_photo"`
	Nama         string         `json:"nama"`
	Work         string         `json:"work"`
	Lokasi       string         `json:"lokasi"`
	Rate         float32        `json:"rate"`
	Deskripsi    string         `json:"deskripsi"`
	Bidang       []RespInterest `json:"bidang"`
	Skill        []RespSkill    `json:"skill"`
	Exp          []ExpResp      `json:"exp"`
}
