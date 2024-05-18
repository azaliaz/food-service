package postgres

import (
	"database/sql"
	"github.com/azaliaz/food-service/internal/models"
	"log"

	_ "github.com/lib/pq"
)

type Postgres struct {
	Conn *sql.DB
}

func NewPostgres(dbURL string) (Postgres, error) {
	db := Postgres{}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return db, err
	}
	db.Conn = conn

	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	log.Println("database connection established")
	return db, nil
}

func (p Postgres) InsertProduct(username string, product models.Product) error {
	query := `INSERT INTO products (name, mealtype, calories, protein, fat, carbohydrates, grams, username) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := p.Conn.Exec(query, product.Name, product.Mealtype, product.Calories, product.Protein, product.Fat, product.Carbohydrates, product.Grams, username)
	if err != nil {
		return err
	}

	return nil
}
func (p Postgres) DeleteProduct(username, id string) error {
	query := "DELETE FROM products WHERE id = $1 and username=$2"
	_, err := p.Conn.Exec(query, id, username)
	if err != nil {
		return err
	}
	return nil
}
func (p Postgres) GetProducts(username, mealtype string) ([]models.Product, error) {
	query := `SELECT name, mealtype, fat, grams, protein, carbohydrates, calories, id FROM products WHERE mealtype=$1 and username=$2`
	rows, err := p.Conn.Query(query, mealtype, username)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	products := make([]models.Product, 0)
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.Name, &product.Mealtype, &product.Fat, &product.Grams, &product.Protein, &product.Carbohydrates, &product.Calories, &product.ID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (p Postgres) GetProduct(username, id string) (models.Product, error) {
	query := `SELECT name, mealtype, fat, grams, protein, carbohydrates, calories, id from products WHERE id = $1 and username = $2`
	row := p.Conn.QueryRow(query, id, username)
	if row.Err() != nil {
		return models.Product{}, row.Err()
	}

	var product models.Product
	err := row.Scan(&product.Name, &product.Mealtype, &product.Fat, &product.Grams, &product.Protein, &product.Carbohydrates, &product.Calories, &product.ID)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (p Postgres) GetSumCalories(username, mealtype string) (float64, float64, float64, float64, error) {
	query := `SELECT calories, protein, fat, carbohydrates FROM products WHERE mealtype = $1 and username = $2`
	rows, err := p.Conn.Query(query, mealtype, username)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	var totalcalories, totalProtein, totalFat, totalCarbohydrates float64
	for rows.Next() {
		var calories, protein, fat, carbohydrates float64
		err = rows.Scan(&calories, &protein, &fat, &carbohydrates)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		totalcalories += calories
		totalProtein += protein
		totalFat += fat
		totalCarbohydrates += carbohydrates
	}
	return totalcalories, totalProtein, totalFat, totalCarbohydrates, nil
}