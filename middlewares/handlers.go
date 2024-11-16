package middlewares

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type respose struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func CreateNewStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("unable to decode the request body. %v", err)
	}
	insertID := insertStock(stock)
	res := respose{
		ID:      insertID,
		Message: "stock successfully created",
	}
	json.NewEncoder(w).Encode(res)
}

func GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()
	if err != nil {
		log.Fatalf("error while fetching all stocks: %v", err)
	}
	json.NewEncoder(w).Encode(stocks)

}
func UpdateStock(w http.ResponseWriter, r *http.Request) {}
func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockId, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("invalid stock id, %v", err)
	}
	stock, err := getStock(int64(stockId))
	if err != nil {
		log.Fatalf("error while getting the stock, %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}
func DeleteStock() {}

// db actions

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlstatement := `INSERT into stocks(name, price, company) values ($1, $2, $3) returning stockid`
	var id int64

	err := db.QueryRow(sqlstatement, stock.Name, stock.Company, stock.Price).Scan(id)
	if err != nil {
		log.Fatalf("error executing insert stock query. %v", err)
	}
	fmt.Printf("Inserted a stock record %v", id)
	return id
}

func getStock(stockId int64) (models.Stock, error) {
	db := createConnection()
	defer db.Close()

	var stock models.Stock
	sqlStatement := `select * from stocks where stockid=$1`
	row := db.QueryRow(sqlStatement, stockId)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("unable to scan the row. %v", err)

	}

	return stock, err

}

func getAllStocks() ([]models.Stock, error) {

	db := createConnection()
	defer db.Close()

	sqlStatment := `select * from stocks`
	var stocks []models.Stock
	rows, err := db.Query(sqlStatment)
	if err != nil {
		log.Fatalf("error while fetching stocks")
	}
	for rows.Next() {
		var stock models.Stock

		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			return nil, err
		}

		stocks = append(stocks, stock)
	}

	return stocks, err
}
