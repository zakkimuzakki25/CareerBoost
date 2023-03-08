package entity

import "gorm.io/gorm"

const (
	GoLang_skill     = "Go-lang"
	HTML_skill       = "HTML"
	CSS_skill        = "CSS"
	JavaScript_skill = "Javascript"
	PHP_skill        = "PHP"
	Phyton_skill     = "Phyton"
	Ruby_skill       = "Ruby"
	React_skill      = "React"
	MySQL_skill      = "MySQL"
	Java_skill       = "Java"
)

type Skill struct {
	ID     uint     `gorm:"primaryKey"`
	Nama   string   `gorm:"type:VARCHAR(30)"`
	Mentor []Mentor `gorm:"many2many:mentor_skill;"`
	User   []User   `gorm:"many2many:user_skill;"`
}

func SeedSkills(sql *gorm.DB) error {
	var categories []Skill

	if err := sql.First(&categories).Error; err != gorm.ErrRecordNotFound {
		return err
	}
	categories = []Skill{
		{
			Nama: GoLang_skill,
		},
		{
			Nama: HTML_skill,
		},
		{
			Nama: CSS_skill,
		},
		{
			Nama: JavaScript_skill,
		},
		{
			Nama: Java_skill,
		},
		{
			Nama: MySQL_skill,
		},
		{
			Nama: PHP_skill,
		},
		{
			Nama: Phyton_skill,
		},
		{
			Nama: React_skill,
		},
		{
			Nama: Ruby_skill,
		},
	}

	if err := sql.Create(&categories).Error; err != nil {
		return err
	}
	return nil
}
