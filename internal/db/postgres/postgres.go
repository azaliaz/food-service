package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/salavad/food-service/internal/models"
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

func (p Postgres) CreateUser(user *models.User) error {
	err := user.BeforeCreating()
	if err != nil {
		return err
	}

	return p.Conn.QueryRow(`INSERT INTO users (email, encrypted_password) VALUES($1, $2) RETURNING id`,
		user.Email,
		user.EncryptedPassword,
	).Scan(&user.ID)
}

func (p Postgres) Find(id int) (*models.User, error) {
	user := &models.User{}
	err := p.Conn.QueryRow(`SELECT id, email, encrypted_password FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.EncryptedPassword)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p Postgres) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := p.Conn.QueryRow(`SELECT id, email, encrypted_password FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.EncryptedPassword)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p Postgres) InsertProduct(product models.Product) error {
	query := `INSERT INTO products (name, mealtype, calories, protein, fat, carbohydrates, grams) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := p.Conn.Exec(query, product.Name, product.Mealtype, product.Calories, product.Protein, product.Fat, product.Carbohydrates, product.Grams)
	if err != nil {
		return err
	}

	return nil
}
func (p Postgres) DeleteProduct(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := p.Conn.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (p Postgres) GetProducts(mealtype string) ([]models.Product, error) {
	query := `SELECT name, mealtype, fat, grams, protein, carbohydrates, calories, id FROM products WHERE mealtype=$1`
	rows, err := p.Conn.Query(query, mealtype)
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

func (p Postgres) GetSumCalories(mealtype string) (float64, float64, float64, float64, error) {
	query := `SELECT calories, protein, fat, carbohydrates FROM products WHERE mealtype = $1`
	rows, err := p.Conn.Query(query, mealtype)
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
