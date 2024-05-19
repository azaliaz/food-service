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

func (p Postgres) AuthUser(username, password string) (int, error) {

	query := `SELECT user_id FROM users WHERE username = $1 AND password = $2`
	row := p.Conn.QueryRow(query, username, password)
	if row.Err() != nil {
		return -1, row.Err()
	}

	var user_id int
	err := row.Scan(&user_id)
	if err == sql.ErrNoRows {
		query := `INSERT INTO users (username, password) VALUES ($1, $2)`
		_, err := p.Conn.Exec(query, username, password)
		if err != nil {
			return -1, err
		}
		query = `SELECT user_id FROM users WHERE username = $1 AND password = $2`
		row := p.Conn.QueryRow(query, username, password)
		if row.Err() != nil {
			return -1, row.Err()
		}

		err = row.Scan(&user_id)
		if err != nil {
			return -1, err
		}
	}
	return user_id, nil

}

func (p Postgres) InsertProduct(userID int, product models.Product, date string) error {
	query := `INSERT INTO products (name, mealtype, calories, protein, fat, carbohydrates, grams, user_id, eating_date) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := p.Conn.Exec(query, product.Name, product.Mealtype, product.Calories, product.Protein, product.Fat, product.Carbohydrates, product.Grams, userID, date)
	if err != nil {
		return err
	}

	return nil
}
func (p Postgres) DeleteProduct(userID int, productID int, date string) error {
	query := "DELETE FROM products WHERE product_id = $1 and user_id=$2 and eating_date = $3"
	_, err := p.Conn.Exec(query, productID, userID, date)
	if err != nil {
		return err
	}
	return nil
}
func (p Postgres) GetProducts(userID int, mealtype string, date string) ([]models.Product, error) {
	query := `SELECT name, mealtype, fat, grams, protein, carbohydrates, calories, product_id FROM products WHERE mealtype=$1 and user_id=$2 and eating_date = $3`
	rows, err := p.Conn.Query(query, mealtype, userID, date)
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

func (p Postgres) GetProduct(userID int, productID int, date string) (models.Product, error) {
	query := `SELECT name, mealtype, fat, grams, protein, carbohydrates, calories, product_id from products WHERE product_id = $1 and user_id = $2 and eating_date = $3`
	row := p.Conn.QueryRow(query, productID, userID, date)
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

func (p Postgres) GetSumCalories(userID int, mealtype string, date string) (float64, float64, float64, float64, error) {
	query := `SELECT calories, protein, fat, carbohydrates FROM products WHERE mealtype = $1 and user_id = $2 and eating_date = $3`
	rows, err := p.Conn.Query(query, mealtype, userID, date)
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
