package db

import (
	"github.com/hyperremix/handball-analyzer/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dsn = "host=localhost user=handball_analyzer password=handball_analyzer dbname=handball_analyzer port=5432 sslmode=disable TimeZone=Europe/Berlin"
var db *gorm.DB

func Get() *gorm.DB {
	if db != nil {
		return db
	}

	newDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	db = newDb

	db.AutoMigrate(&model.Referee{}, &model.Season{}, &model.League{}, &model.Team{}, &model.TeamMember{}, &model.GameEventBlueCard{}, &model.GameEventGoal{}, &model.GameEventPenalty{}, &model.GameEventRedCard{}, &model.GameEventSevenMeters{}, &model.GameEventTimeout{}, &model.GameEventYellowCard{}, &model.Game{})

	return db
}
