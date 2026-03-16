package main

import (
	"fmt"
	"log"
	"os"
	httphandlers "test_project/httpHandlers"
	"test_project/repository"
	"test_project/server"
	"test_project/service"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load("config.env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Соединение успешно установлено")
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(25)
	repo := repository.NewDbRepo(db)
	walletService := service.NewWalletService(repo)
	httpHandlers := httphandlers.NewHTTPHandlers(walletService)
	srv := server.NewHTTPServer(httpHandlers)
	log.Fatal(srv.StartServer())

}
