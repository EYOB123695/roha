package main

import (
	"log"

	"github.com/EYOB123695/roha/initializers"
	model "github.com/EYOB123695/roha/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	err := initializers.DB.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
		&model.Tag{},
		&model.Like{},
		&model.Bookmark{},
		&model.UserActivityLog{},
		&model.UserInterest{},
	)
	if err != nil {
		log.Fatal("failed to migrate database")
	}
}