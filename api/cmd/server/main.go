package main

import (
	"context"
	"fmt"
	"github.com/azaliaz/food-service/internal/db/postgres"
	"github.com/azaliaz/food-service/internal/models"
	"github.com/azaliaz/food-service/internal/server/rest"
	"log"
	"os"
)

func main() {
	db, err := NewDB(os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"))
	if err != nil {
		log.Fatal(err)
	}

	server, err := NewServer("3000", db)
	if err != nil {
		log.Fatal(err)
	}
	server.Start(context.Background())

}

type Storage interface {
	AuthUser(username, password string) (int, error)
	InsertProduct(userID int, product models.Product, date string) error
	GetProducts(userID int, mealtype string, date string) ([]models.Product, error)
	GetProduct(userID int, productID int, date string) (models.Product, error)
	GetSumCalories(userID int, mealtype string, date string) (float64, float64, float64, float64, error)
	DeleteProduct(userID int, productID int, date string) error
}

func NewDB(host, port, username, password, dbname string) (Storage, error) {
	db, err := postgres.NewPostgres(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbname))
	if err != nil {
		return nil, err
	}
	return db, nil
}

type Server interface {
	Start(context.Context)
}

func NewServer(port string, database Storage) (Server, error) {
	return rest.Server{Port: port, Database: database}, nil
}
