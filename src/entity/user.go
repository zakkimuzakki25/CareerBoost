package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName     string     `json:"full_name" gorm:"type:VARCHAR(255);NOT NULL"`
	Username     string     `json:"username" gorm:"type:VARCHAR(20);UNIQUE"`
	Email        string     `json:"email" gorm:"type:VARCHAR(255);UNIQUE"`
	Password     string     `json:"password" gorm:"type:VARCHAR(255);NOT NULL"`
	Skill        []Skill    `json:"skillID" gorm:"many2many:users_skill;default:null"`
	Interest     []Interest `json:"interestID" gorm:"many2many:users_interest"`
	TanggalLahir time.Time  `json:"tanggal_lahir" gorm:"default:null"`
	TempatLahir  string     `json:"tempat_lahir" gorm:"type:VARCHAR(50);default:null"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(50);default:null"`
	ProfilePhoto string     `json:"profile_photo" gorom:"default:null"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:VARCHAR(250);default:null"`
	Mentor       []Mentor   `json:"mentors" gorm:"many2many:user_mentors"`
	Magang       []Magang   `json:"magangs" gorm:"many2many:user_magangs"`
	Course       []Course   `json:"courses" gorm:"many2many:user_courses"`
}

type UserRegister struct {
	FullName   string `json:"full_name" gorm:"NOT NULL" binding:"required"`
	Username   string `json:"username" gorm:"NOT NULL" binding:"required,min=5,max=20"`
	Email      string `json:"email" gorm:"NOT NULL" binding:"required,email"`
	Password   string `json:"password" gorm:"NOT NULL" binding:"required,min=8"`
	InterestID []uint `json:"interestID" gorm:"NOT NULL"`
}

type UserLogin struct {
	Email    string `json:"email" gorm:"NOT NULL" binding:"required,email"`
	Password string `json:"password" gorm:"NOT NULL" binding:"required"`
}

type UserProfilePage struct {
	Email        string         `json:"email"`
	FullName     string         `json:"full_name"`
	Lokasi       string         `json:"lokasi"`
	ProfilePhoto string         `json:"profile_photo"`
	Deskripsi    string         `json:"deskripsi"`
	TanggalLahir time.Time      `json:"tanggal_lahir"`
	TempatLahir  string         `json:"tempat_lahir"`
	InterestID   []RespInterest `json:"interest"`
}

type UserHome struct {
	FullName     string `json:"full_name"`
	ProfilePhoto string `json:"profile_photo"`
}

type UserMentorHistory struct {
	Interest string `json:"interest"`
	FullName string `json:"full_name"`
}

type UserClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewUserClaims(id uint, exp time.Duration) UserClaims {
	return UserClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
