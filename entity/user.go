package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName string     `json:"full_name" gorm:"type:VARCHAR(100);NOT NULL"`
	Username string     `json:"username" gorm:"type:VARCHAR(20);UNIQUE"`
	Email    string     `json:"email" gorm:"type:VARCHAR(200);UNIQUE"`
	Password string     `json:"password" gorm:"type:VARCHAR(255);NOT NULL"`
	Skills   []string   `json:"skills" gorm:"type:VARCHAR(20)"`
	Interest []Interest `json:"interestID" gorm:"many2many:users_interest;"`
}

type UserRegister struct {
	FullName   string `json:"full_name" gorm:"type:VARCHAR(100);NOT NULL"`
	Username   string `json:"username" gorm:"type:VARCHAR(20)"`
	Email      string `json:"email" gorm:"type:VARCHAR(200);UNIQUE"`
	Password   string `json:"password" gorm:"type:VARCHAR(255);NOT NULL"`
	InterestID []uint `json:"interestID" gorm:";NOT NULL"`
}

type UserLogin struct {
	Email    string `json:"email" gorm:"type:VARCHAR(200);NOT NULL"`
	Password string `json:"password" gorm:"type:VARCHAR(255);NOT NULL"`
}

type UserClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewUserClaims(email string, exp time.Duration) UserClaims {
	return UserClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
