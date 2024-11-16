package router

import (
	"github.com/gorilla/mux"
	"go-postgres/middlewares"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/stock/{id}", middlewares.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/stocks", middlewares.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/stock", middlewares.CreateNewStock).Methods("POST", "OPTIONS")
	// router.HandleFunc("/deletestock/{id}", middlewares.DeleteStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/stock/{id}", middlewares.UpdateStock).Methods("DELETE", "OPTIONS")

	return router
}
