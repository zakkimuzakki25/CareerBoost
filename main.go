package main

import (
	"CareerBoost/entity"
	"CareerBoost/handler"
	"CareerBoost/sdk/config"
	"CareerBoost/sdk/database"
	"log"

	"gorm.io/gorm"
)

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

	sql, err := database.InitDB(databaseConfig)
	if err != nil {
		log.Fatal("init db failed,", err)
	}

	db := sql.GetInstance()
	db.AutoMigrate(entity.User{}, entity.Interest{}, entity.Course{})

	if err := seedInterest(db); err != nil {
		panic("GAGAL SEED")
	}

	handler := handler.Init(cnfg, db)
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
