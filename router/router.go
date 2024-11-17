package router

import (
	"encoding/json"
	"go-postgres/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/stock/{id}", middlewares.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/stocks", middlewares.GetAllStocks).Methods("GET", "OPTIONS")
	router.HandleFunc("/stock", middlewares.CreateNewStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/deletestock/{id}", middlewares.DeleteStock).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/stock/{id}", middlewares.UpdateStock).Methods("PUT", "OPTIONS")

	// handle anyother requests
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Route not found",
		})
	})

	return router
}
