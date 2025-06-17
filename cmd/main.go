package main

import (
	"log"
	"os"
	"os/signal"

	"golang.org/x/net/context"
	"shop"
	"shop/internal/handlers"
	"shop/internal/repository"
	"shop/internal/services"
	"shop/storage/postgresql"
)

const dsn = "host=localhost user=postgres dbname=simple_shop password=postgres sslmode=disable port=8091"
const port = "8080"

func main() {
	storage := postgresql.New(dsn)

	defer storage.Close()

	//storage.DB.AutoMigrate(&models.Locations{},
	//	&models.MeasureUnits{},
	//	&models.Products{},
	//	&models.Orders{},
	//	&models.OrderItems{})

	repos := repository.NewRepository(storage.DB)
	service := services.NewService(repos)
	handler := handlers.NewHandler(service)

	srv := new(shop.Server)
	go func() {
		if err := srv.Run(port, handler.InitRoute()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Println("Server started on port " + port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured while shutting down server: %s", err.Error())

	}
}
