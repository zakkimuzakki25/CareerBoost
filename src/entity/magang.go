package entity

import (
	"time"

	"gorm.io/gorm"
)

type Magang struct {
	gorm.Model
	Logo         string     `json:"logo" gorm:"default:null"`
	Perusahaan   string     `json:"perusahaan" gorm:"type:VARCHAR(255);NOT NULL"`
	Skill        []Skill    `json:"skillID" gorm:"many2many:magangs_skill"`
	Interest     []Interest `json:"interestID" gorm:"many2many:magangs_interest"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(255);NOT NULL"`
	Applican     []User     `json:"Applican" gorm:"foreignkey:MentorID"`
	StatusMagang string     `json:"status_magang" gorm:"type:VARCHAR(255);NOT NULL"`
	JangkaWaktu  string     `json:"jangka_waktu" gorm:"type:VARCHAR(255);NOT NULL"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:text"`
	Rate         float32    `json:"rate"`
	Fee          float32    `json:"fee"`
}

type MagangReqByID struct {
	ID uint `json:"id"`
}

type MagangRespData struct {
	CreatedAt    time.Time      `json:"createdAt"`
	Logo         string         `json:"logo"`
	Perusahaan   string         `json:"perusahaan"`
	Skill        []RespSkill    `json:"skill"`
	Interest     []RespInterest `json:"interest"`
	Lokasi       string         `json:"lokasi"`
	Applied      uint           `json:"applied"`
	Deskripsi    string         `json:"deskripsi"`
	Rate         float32        `json:"rate"`
	Fee          float32        `json:"fee"`
	JangkaWaktu  string         `json:"jangka_waktu"`
	StatusMagang string         `json:"status_magang"`
}

type MagangAdd struct {
	gorm.Model
	Logo         string  `json:"logo"`
	Perusahaan   string  `json:"perusahaan"`
	Skill        []uint  `json:"skillID"`
	Interest     []uint  `json:"interestID"`
	Lokasi       string  `json:"lokasi"`
	Deskripsi    string  `json:"deskripsi"`
	Rate         float32 `json:"rate"`
	Fee          float32 `json:"fee"`
	JangkaWaktu  string  `json:"jangka_waktu"`
	StatusMagang string  `json:"status_magang"`
}

type MagangParam struct {
	PostID int64 `uri:"magang_id" gorm:"column:id"`
	PaginationParam
}
