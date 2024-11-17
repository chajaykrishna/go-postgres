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
	_ "github.com/lib/pq"
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
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stockid, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("error converting stockid to int. %v", err)
	}
	var stock models.Stock
	json.NewDecoder(r.Body).Decode(&stock)

	updatedRows := updateStock(stockid, stock)

	if err != nil {
		log.Fatalf("error while updating the stock. %v", err)
	}
	msg := fmt.Sprintf("stock updated. total rows affected: %v", updatedRows)
	res := respose{
		ID:      int64(stockid),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

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
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	stockid, err := strconv.Atoi(param["id"])
	if err != nil {
		log.Fatalf("invalid stock id, %v", err)
	}

	id := deleteStock(int64(stockid))
	msg := fmt.Sprintf("stock successfully deleted. id: %v", id)
	res := respose{
		ID:      id,
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

// db actions

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlstatement := `INSERT into stocks(name, price, company) values ($1, $2, $3) returning stockid`
	var id int64

	err := db.QueryRow(sqlstatement, stock.Name, stock.Price, stock.Company).Scan(&id)
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

func updateStock(stockid int, updatedstock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `update stocks set name=$2, price=$3, company=$4 where stockid=$1`

	res, err := db.Exec(sqlStatement, stockid, updatedstock.Name, updatedstock.Price, updatedstock.Company)

	if err != nil {
		log.Fatalf("unable to execute the query")
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while updating the stock. %v", err)
	}
	return rowsAffected
}

func deleteStock(stockid int64) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `delete from stocks where stockid=$1`
	res, err := db.Exec(sqlStatement, stockid)
	if err != nil {
		log.Fatalf("Error while executing the delete query. %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	return rowsAffected
}
