package main

import (
	"CareerBoost/sdk/config"
	"CareerBoost/sdk/database"
	"CareerBoost/src/entity"
	"CareerBoost/src/handler"
	"fmt"
	"log"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
	"gorm.io/gorm"
)

// test
func main() {
	cnfg := config.Init()
	if err := cnfg.Load(".env"); err != nil {
		log.Fatalln("failed to load env file")
	}

	databaseConfig := database.Config{
		Username: cnfg.Get("DB_USERNAME"),
		Password: cnfg.Get("DB_PASSWORD"),
		Host:     cnfg.Get("DB_HOST"),
		Port:     cnfg.Get("DB_PORT"),
		Database: cnfg.Get("DB_DATABASE"),
	}

	supClient := supabasestorageuploader.NewSupabaseClient(
		cnfg.Get("PROJECT_URL"),
		cnfg.Get("PROJECT_API_KEYS"),
		cnfg.Get("STORAGE_NAME"),
		cnfg.Get("STORAGE_PATH"),
	)

	sql, err := database.InitDB(databaseConfig)
	if err != nil {
		log.Fatal("init db failed,", err)
	}

	db := sql.GetInstance()
	db.AutoMigrate(entity.User{})
	db.AutoMigrate(entity.Interest{})
	db.AutoMigrate(entity.Skill{})
	db.AutoMigrate(entity.Mentor{})
	db.AutoMigrate(entity.Exp{})
	db.AutoMigrate(entity.Video{})
	db.AutoMigrate(entity.Playlist{})
	db.AutoMigrate(entity.Course{})
	db.AutoMigrate(entity.Magang{})

	if err := seedInterest(db); err != nil {
		fmt.Println(err)
		panic("GAGAL SEED INTEREST")
	}
	if err := seedSkills(db); err != nil {
		fmt.Println(err)
		panic("GAGAL SEED INTEREST")
	}

	handler := handler.Init(cnfg, db, supClient)
	handler.Run()

}

func seedInterest(sql *gorm.DB) error {
	var categories []entity.Interest

	if err := sql.First(&categories).Error; err != gorm.ErrRecordNotFound {
		return err
	}
	categories = []entity.Interest{
		{
			Nama: entity.FrontEndInterest,
		},
		{
			Nama: entity.BackEndInterest,
		},
		{
			Nama: entity.DataScienceInterest,
		},
		{
			Nama: entity.ArtificialInteligenceInterest,
		},
		{
			Nama: entity.CyberSecurityInterest,
		},
	}

	if err := sql.Create(&categories).Error; err != nil {
		return err
	}
	return nil
}

func seedSkills(sql *gorm.DB) error {
	var categories []entity.Skill

	if err := sql.First(&categories).Error; err != gorm.ErrRecordNotFound {
		return err
	}
	categories = []entity.Skill{
		{
			Nama: entity.GoLang_skill,
		},
		{
			Nama: entity.HTML_skill,
		},
		{
			Nama: entity.CSS_skill,
		},
		{
			Nama: entity.JavaScript_skill,
		},
		{
			Nama: entity.Java_skill,
		},
		{
			Nama: entity.MySQL_skill,
		},
		{
			Nama: entity.PHP_skill,
		},
		{
			Nama: entity.Phyton_skill,
		},
		{
			Nama: entity.React_skill,
		},
		{
			Nama: entity.Ruby_skill,
		},
	}

	if err := sql.Create(&categories).Error; err != nil {
		return err
	}
	return nil
}
