package handlers

import (
	"encoding/json"
	"github.com/azaliaz/food-service/internal/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
	Storage Storage
}

type Storage interface {
	InsertProduct(username string, product models.Product) error
	GetProducts(username, mealtype string) ([]models.Product, error)
	GetProduct(username, id string) (models.Product, error)
	GetSumCalories(username, mealtype string) (float64, float64, float64, float64, error)
	DeleteProduct(username, id string) error
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

func (h *handler) CheckOptions(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Обработка запроса OPTIONS для CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Vary", "Origin")
	w.WriteHeader(http.StatusOK)
	log.Println("Options request handled successfully")
	return
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	type Request struct {
		Username string         `json:"username"`
		Product  models.Product `json:"product"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(req)

	err = h.Storage.InsertProduct(req.Username, req.Product)
	if err != nil {
		log.Println(err)
		return
	}
	breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories(req.Username, "breakfast")
	if err != nil {
		log.Println(err)
		return
	}
	lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories(req.Username, "lunch")
	if err != nil {
		log.Println(err)
		return
	}
	dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories(req.Username, "dinner")
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
		SavedProduct: req.Product,
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodDelete {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	username := r.URL.Query().Get("username")
	product, err := h.Storage.GetProduct(username, id)
	if err != nil {
		log.Println(err)
		return
	}
	err = h.Storage.DeleteProduct(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories(username, "breakfast")
	if err != nil {
		log.Println(err)
		return
	}
	lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories(username, "lunch")
	if err != nil {
		log.Println(err)
		return
	}
	dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories(username, "dinner")
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
		Message:      "successfully deleted product",
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

	body, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	mealtype := r.URL.Query().Get("mealtype")
	username := r.URL.Query().Get("username")
	if mealtype == "" {
		breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories(username, "breakfast")
		if err != nil {
			log.Println(err)
			return
		}
		lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories(username, "lunch")
		if err != nil {
			log.Println(err)
			return
		}
		dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories(username, "dinner")
		if err != nil {
			log.Println(err)
			return
		}

		type Response struct {
			Message       string
			TotalNutriens struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			} `json:"totalNutriens"`
		}

		resp := Response{
			Message: "successfully created product",
			TotalNutriens: struct {
				Calories     float64 `json:"calories"`
				Protein      float64 `json:"protein"`
				Fat          float64 `json:"fat"`
				Carbohydrate float64 `json:"carbohydrate"`
			}{
				Calories:     breakfastCalories + lunchCalories + dinnerCalories,
				Protein:      breakfastProtein + lunchProtein + dinnerProtein,
				Fat:          breakfastFat + lunchFat + dinnerFat,
				Carbohydrate: breakfastCarbo + lunchCarbo + dinnerCarbo,
			},
		}
		body, err := json.Marshal(resp)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(body)
		return
	}

	products, err := h.Storage.GetProducts(username, mealtype)
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
