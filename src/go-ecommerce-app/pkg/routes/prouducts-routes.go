package routes

import (
	"github.com/NorikKhachatryan01/go-ecommerce/pkg/controllers"
	"github.com/gorilla/mux"
)

var ProductStoreRoutes = func (router *mux.Router)  {

	router.HandleFunc("/products", controllers.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", controllers.GetProductById).Methods("GET")
	router.HandleFunc("/products", controllers.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", controllers.DeleteProduct).Methods("DELETE")
}