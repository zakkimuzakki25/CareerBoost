package entity

const (
	FrontEndInterest              = "Front End Developer"
	BackEndInterest               = "Back End Developer"
	DataScienceInterest           = "Data Science"
	ArtificialInteligenceInterest = "Artificial Inteligence"
	CyberSecurityInterest         = "Cyber Security"
)

type Interest struct {
	ID     uint     `gorm:"primaryKey"`
	Nama   string   `gorm:"type:VARCHAR(30)"`
	User   []User   `gorm:"many2many:users_interest"`
	Mentor []Mentor `gorm:"many2many:mentors_interest"`
	Magang []Magang `gorm:"many2many:magangs_interest"`
	Course []Course
}

type RespInterest struct {
	Nama string `json:"nama"`
}

type RespInterestWithID struct {
	Nama string `json:"label"`
	ID   string `json:"value"`
}
