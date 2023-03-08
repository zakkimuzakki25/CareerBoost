package entity

const (
	FrontEndInterest              = "Front End"
	BackEndInterest               = "Back End"
	DataScienceInterest           = "Data Science"
	ArtificialInteligenceInterest = "Artificial Inteligence"
	CyberSecurityInterest         = "Cyber Security"
)

type Interest struct {
	ID   uint   `gorm:"primaryKey"`
	Nama string `gorm:"type:VARCHAR(30)"`
	User []User `gorm:"many2many:users_interest"`
}
