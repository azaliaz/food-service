package rest

import (
	"context"
	"github.com/azaliaz/food-service/internal/models"
	"github.com/azaliaz/food-service/internal/server/rest/handlers"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Port     string
	Database Storage
}

type Storage interface {
	AuthUser(username, password string) (int, error)
	InsertProduct(userID int, product models.Product, date string) error
	GetProducts(userID int, mealtype string, date string) ([]models.Product, error)
	GetProduct(userID int, productID int, date string) (models.Product, error)
	GetSumCalories(userID int, mealtype string, date string) (float64, float64, float64, float64, error)
	DeleteProduct(userID int, productID int, date string) error
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
