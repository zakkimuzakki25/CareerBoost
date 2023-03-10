package entity

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
	User   []User   `gorm:"many2many:users_skill"`
	Mentor []Mentor `gorm:"many2many:mentors_skill"`
}

type RespSkill struct {
	Nama string `json:"nama"`
}
