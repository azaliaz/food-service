package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/salavad/food-service/internal/models"
)

type handler struct {
	Storage Storage
}

type Storage interface {
	InsertProduct(product models.Product) error
	GetProducts(mealtype string) ([]models.Product, error)
	GetSumCalories(mealtype string) (float64, float64, float64, float64, error)
	DeleteProduct(id string) error
}

type Handler interface {
	Register(router *httprouter.Router)
	CreateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	GetProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params)
	DeleteProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params)
}

func NewHandler(storage Storage) Handler {
	return &handler{storage}
}

func (h *handler) Register(router *httprouter.Router) {
	router.POST("/products", h.CreateProduct)
	router.GET("/products", h.GetProducts)
	router.DELETE("/products", h.DeleteProduct)
	router.OPTIONS("/products", h.CheckOptions)
}

func (h *handler) DeleteOptions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Обработка запроса OPTIONS для CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8888")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Vary", "Origin")
	w.WriteHeader(http.StatusOK)
	log.Println("Options request handled successfully")
	return
}

func (h *handler) CheckOptions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Обработка запроса OPTIONS для CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8888")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Vary", "Origin")
	w.WriteHeader(http.StatusOK)
	log.Println("Options request handled successfully")
	return
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8888")

	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(product)

	err = h.Storage.InsertProduct(product)
	if err != nil {
		log.Println(err)
		return
	}
	breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories("breakfast")
	if err != nil {
		log.Println(err)
		return
	}
	lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories("lunch")
	if err != nil {
		log.Println(err)
		return
	}
	dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories("dinner")
	if err != nil {
		log.Println(err)
		return
	}
	type Response struct {
		Message       string         `json:"message"`
		SavedProduct  models.Product `json:"savedProduct"`
		TotalCalories struct {
			Breakfast struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"breakfast"`
			Lunch struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"lunch"`
			Dinner struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"dinner"`
		} `json:"totalCalories"`
	}

	resp := Response{
		Message:      "successfully created product",
		SavedProduct: product,
		TotalCalories: struct {
			Breakfast struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"breakfast"`
			Lunch struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"lunch"`
			Dinner struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"dinner"`
		}{
			Breakfast: struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			}{
				Calories:     breakfastCalories,
				Protein:      breakfastProtein,
				Fat:          breakfastFat,
				Carbohydrate: breakfastCarbo,
			},
			Lunch: struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			}{
				Calories:     lunchCalories,
				Protein:      lunchProtein,
				Fat:          lunchFat,
				Carbohydrate: lunchCarbo,
			},
			Dinner: struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			}{
				Calories:     dinnerCalories,
				Protein:      dinnerProtein,
				Fat:          dinnerFat,
				Carbohydrate: dinnerCarbo,
			},
		},
	}

	log.Println(resp)

	body, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))
	w.Write(body)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8888")
	if r.Method != http.MethodDelete {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	err := h.Storage.DeleteProduct(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "product deleted"})
}

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:8888")

	mealtype := r.URL.Query().Get("mealtype")

	products, err := h.Storage.GetProducts(mealtype)
	if err != nil {
		log.Println(err)
		return
	}
	type Response struct {
		Message string
		Data    []models.Product
	}
	resp := &Response{Message: "OK", Data: products}
	body, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(body)
}
