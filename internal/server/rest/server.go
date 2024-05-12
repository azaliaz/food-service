package rest

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/salavad/food-service/internal/models"
	"github.com/salavad/food-service/internal/server/rest/handlers"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Port     string
	Database Storage
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

func (s Server) Start(ctx context.Context) {
	log.Println("creating router")
	router := httprouter.New()

	log.Println("register user handler")
	handler := handlers.NewHandler(s.Database)
	handler.Register(router)

	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		panic(err)
	}
	log.Println("server is listening at " + s.Port)

	server := http.Server{
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		log.Println(server.Serve(listener))
	}()

	go func() {
		defer wg.Done()

		<-ctx.Done()
		log.Println(server.Shutdown(ctx))

	}()

	wg.Wait()
}
