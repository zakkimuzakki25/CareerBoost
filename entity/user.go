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
	Skills       []string   `json:"skills" gorm:"many2many:users_skill"`
	Interest     []Interest `json:"interestID" gorm:"many2many:users_interest"`
	TanggalLahir *time.Time `json:"tanggal_lahir"`
	TempatLahir  string     `json:"tempat_lahir" gorm:"type:VARCHAR(50)"`
	Lokasi       string     `json:"lokasi" gorm:"type:VARCHAR(50)"`
	Courses      []Course   `json:"courses" gorm:"many2many:users_courses;"`
	ProfilePhoto string     `json:"profile_photo"`
	Deskripsi    string     `json:"deskripsi" gorm:"type:VARCHAR(250)"`
}

type UserRegister struct {
	FullName   string `json:"full_name" gorm:"NOT NULL"`
	Username   string `json:"username" gorm:"NOT NULL"`
	Email      string `json:"email" gorm:"NOT NULL"`
	Password   string `json:"password" gorm:"NOT NULL"`
	InterestID []uint `json:"interestID" gorm:"NOT NULL"`
}

type UserLogin struct {
	Email    string `json:"email" gorm:"NOT NULL"`
	Password string `json:"password" gorm:"NOT NULL"`
}

type UserProfilePage struct {
	Email        string     `json:"email"`
	FullName     string     `json:"full_name"`
	Lokasi       string     `json:"lokasi"`
	ProfilePhoto string     `json:"profile_photo"`
	Deskripsi    string     `json:"deskripsi"`
	TanggalLahir time.Time  `json:"tanggal_lahir"`
	TempatLahir  string     `json:"tempat_lahir"`
	InterestID   []Interest `json:"interestID"`
}

type UserHome struct {
	FullName     string `json:"full_name"`
	ProfilePhoto string `json:"profile_photo"`
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
