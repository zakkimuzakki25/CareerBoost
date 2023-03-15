package entity

import (
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model
	Judul      string     `json:"judul" gorm:"type:VARCHAR(255)"`
	Foto       string     `json:"foto" gorm:"type:VARCHAR(255)"`
	Deskripsi  string     `json:"deskripsi" gorm:"type:text"`
	Intro      string     `json:"intro" gorm:"type:text"`
	InterestID uint       `json:"interest_id"`
	Playlist   []Playlist `json:"playlist"`
	Rate       float32    `json:"rate"`
	Price      float32    `json:"price"`
	Vote       uint       `json:"vote"`
}

type CourseAdd struct {
	Foto       string     `json:"foto"`
	Judul      string     `json:"judul"`
	Deskripsi  string     `json:"deskripsi"`
	Intro      string     `json:"intro"`
	Playlist   []Playlist `json:"playlist"`
	Rate       float32    `json:"rate"`
	Price      float32    `json:"price"`
	InterestID uint       `json:"interest_id"`
	Vote       uint       `json:"vote"`
}

type CourseRespData struct {
	Judul     string         `json:"judul"`
	Deskripsi string         `json:"deskripsi"`
	Intro     string         `json:"intro"`
	Playlist  []RespPlaylist `json:"playlist"`
	Rate      float32        `json:"rate"`
	Price     float32        `json:"price"`
}

type CourseParam struct {
	ID int64 `uri:"course_id" gorm:"column:id"`
	PaginationParam
}

type CourseReqByID struct {
	ID uint `json:"id"`
}

type CourseSearch struct {
	Key string `form:"search"`
}

type Playlist struct {
	gorm.Model
	Nama     string  `json:"nama" gorm:"type:VARCHAR(55)"`
	Video    []Video `json:"video"`
	Durasi   string  `json:"durasi" gorm:"type:VARCHAR(55)"`
	Course   Course  `json:"course"`
	CourseID uint    `json:"course_id"`
}

type Video struct {
	gorm.Model
	Link       string `json:"link" gorm:"type:varchar(255)"`
	Judul      string `json:"judul" gorm:"type:varchar(255)"`
	Durasi     string `json:"durasi"`
	PlaylistID uint   `json:"playlist_id"`
}

type RespPlaylist struct {
	Nama     string
	Durasi   string
	Video    []RespVideo
	CourseID uint
}

type RespVideo struct {
	Link       string `json:"link" gorm:"type:varchar(255)"`
	Judul      string `json:"judul" gorm:"type:varchar(255)"`
	Durasi     string `json:"durasi"`
	PlaylistID uint   `json:"playlist_id"`
}
