package handlers

import (
	"encoding/json"
	"github.com/azaliaz/food-service/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type handler struct {
	Storage Storage
}

type Storage interface {
	AuthUser(username, password string) (int, error)
	InsertProduct(userID int, product models.Product, date string) error
	GetProducts(userID int, mealtype string, date string) ([]models.Product, error)
	GetProduct(userID int, productID int, date string) (models.Product, error)
	GetSumCalories(userID int, mealtype string, date string) (float64, float64, float64, float64, error)
	DeleteProduct(userID int, productID int, date string) error
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
	router.POST("/users", h.Auth)
	router.OPTIONS("/users", h.CheckOptions)
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

func (h *handler) Auth(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(req)

	userID, err := h.Storage.AuthUser(req.Username, req.Password)
	if err != nil {
		log.Println(err)
		return
	}
	type Response struct {
		UserID int `json:"user_id"`
	}
	resp := Response{UserID: userID}

	body, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(body))
	w.Write(body)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Set("Access-Control-Allow-Origin", "*")

	type Request struct {
		UserID  string         `json:"userID"`
		Date    string         `json:"date"`
		Product models.Product `json:"product"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(req)

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	err = h.Storage.InsertProduct(userID, req.Product, req.Date)
	if err != nil {
		log.Println(err)
		return
	}
	breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories(userID, "breakfast", req.Date)
	if err != nil {
		log.Println(err)
		return
	}
	lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories(userID, "lunch", req.Date)
	if err != nil {
		log.Println(err)
		return
	}
	dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories(userID, "dinner", req.Date)
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
	productID, err := strconv.Atoi(r.URL.Query().Get("id"))
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	date := r.URL.Query().Get("date")
	product, err := h.Storage.GetProduct(userID, productID, date)
	if err != nil {
		log.Println(err)
		return
	}
	err = h.Storage.DeleteProduct(userID, productID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories(userID, "breakfast", date)
	if err != nil {
		log.Println(err)
		return
	}
	lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories(userID, "lunch", date)
	if err != nil {
		log.Println(err)
		return
	}
	dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories(userID, "dinner", date)
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
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		log.Println(err)
	}
	date := r.URL.Query().Get("date")
	log.Println(mealtype, userID, date)
	if mealtype == "" {
		breakfastCalories, breakfastProtein, breakfastFat, breakfastCarbo, err := h.Storage.GetSumCalories(userID, "breakfast", date)
		if err != nil {
			log.Println(err)
			return
		}
		lunchCalories, lunchProtein, lunchFat, lunchCarbo, err := h.Storage.GetSumCalories(userID, "lunch", date)
		if err != nil {
			log.Println(err)
			return
		}
		dinnerCalories, dinnerProtein, dinnerFat, dinnerCarbo, err := h.Storage.GetSumCalories(userID, "dinner", date)
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

	products, err := h.Storage.GetProducts(userID, mealtype, date)
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
