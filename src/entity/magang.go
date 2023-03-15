package entity

import (
	"gorm.io/gorm"
)

type Magang struct {
	gorm.Model
	Logo         string     `json:"logo" gorm:"default:null"`
	Perusahaan   string     `json:"perusahaan" gorm:"type:VARCHAR(255);NOT NULL"`
	Skill        []Skill    `json:"skillID" gorm:"many2many:magangs_skill"`
	Interest     []Interest `json:"interestID" gorm:"many2many:magangs_interest"`
	JobDesc      []JobDesc  `json:"job_description"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(255);NOT NULL"`
	StatusMagang string     `json:"status_magang" gorm:"type:VARCHAR(255);NOT NULL"`
	JangkaWaktu  string     `json:"jangka_waktu" gorm:"type:VARCHAR(255);NOT NULL"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:text"`
	Rate         float32    `json:"rate"`
	Fee          float32    `json:"fee"`
}

type JobDesc struct {
	gorm.Model
	Tugas    string `json:"tugas"`
	MagangID uint   `json:"magang_id"`
}

type JobDescResp struct {
	Tugas string `json:"tugas"`
}

type MagangReqByID struct {
	ID uint `json:"id"`
}

type MagangRekomendasiData struct {
	ID           uint         `json:"id"`
	Release      string       `json:"release"`
	Logo         string       `json:"logo"`
	Interest     RespInterest `json:"interest"`
	Lokasi       string       `json:"lokasi"`
	Apllied      uint         `json:"applied"`
	StatusMagang string       `json:"status_magang"`
}

type MagangRespData struct {
	Rekomendasi  []MagangRekomendasiData `json:"rekomendasi"`
	Perusahaan   string                  `json:"perusahaan"`
	Skill        []RespSkill             `json:"skill"`
	Interest     []RespInterest          `json:"interest"`
	Deskripsi    string                  `json:"deskripsi"`
	JangkaWaktu  string                  `json:"jangka_waktu"`
	StatusMagang string                  `json:"status_magang"`
	JobDesc      []JobDescResp           `json:"job_description"`
}

type MagangRespHome struct {
	ID           uint           `json:"id"`
	Release      string         `json:"release"`
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

type MagangRespRekomendasi struct {
	ID           uint           `json:"id"`
	Release      string         `json:"release"`
	Logo         string         `json:"logo"`
	Perusahaan   string         `json:"perusahaan"`
	Skill        []RespSkill    `json:"skill"`
	Interest     []RespInterest `json:"interest"`
	Lokasi       string         `json:"lokasi"`
	Applied      uint           `json:"applied"`
	StatusMagang string         `json:"status_magang"`
}

type MagangAdd struct {
	gorm.Model
	Logo         string    `json:"logo"`
	Perusahaan   string    `json:"perusahaan"`
	Skill        []uint    `json:"skillID"`
	Interest     []uint    `json:"interestID"`
	JobDesc      []JobDesc `json:"job_description"`
	Lokasi       string    `json:"lokasi"`
	Deskripsi    string    `json:"deskripsi"`
	Rate         float32   `json:"rate"`
	Fee          float32   `json:"fee"`
	JangkaWaktu  string    `json:"jangka_waktu"`
	StatusMagang string    `json:"status_magang"`
}

type MagangParam struct {
	ID int64 `uri:"magang_id" gorm:"column:id"`
	PaginationParam
}
