package main

import (
	"log"

	"github.com/EYOB123695/roha/initializers"
	"github.com/EYOB123695/roha/repository"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	err := initializers.DB.AutoMigrate(
		&repository.User{},
		&repository.Post{},
		&repository.Comment{},
		&repository.Tag{},
		&repository.Like{},
		&repository.Bookmark{},
		&repository.UserActivityLog{},
		&repository.UserInterest{},
	)
	if err != nil {
		log.Fatal("failed to migrate database: ", err)
	}
}