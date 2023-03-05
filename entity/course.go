package entity

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	// video string `json:"course_video"`
}
