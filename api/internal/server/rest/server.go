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
	InsertProduct(username string, product models.Product) error
	GetProducts(username, mealtype string) ([]models.Product, error)
	GetProduct(username, id string) (models.Product, error)
	GetSumCalories(username, mealtype string) (float64, float64, float64, float64, error)
	DeleteProduct(username, id string) error
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
