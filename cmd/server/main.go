package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/salavad/food-service/internal/db/postgres"
	"github.com/salavad/food-service/internal/models"
	"github.com/salavad/food-service/internal/server/rest"
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
	InsertProduct(product models.Product) error
	GetProducts(mealtype string) ([]models.Product, error)
	GetSumCalories(mealtype string) (float64, float64, float64, float64, error)
	CreateUser(user *models.User) error
	Find(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	DeleteProduct(id string) error
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
