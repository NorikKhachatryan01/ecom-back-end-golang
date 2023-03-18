package main

import (
	"log"
	"net/http"

	"github.com/NorikKhachatryan01/go-ecommerce/pkg/routes"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
func main(){
	r := mux.NewRouter()
	routes.ProductStoreRoutes(r)
	http.Handle("/",r)
	log.Fatal(http.ListenAndServe("localhost:8081",r))
}




































// package main

// import(
// 	"fmt"
// 	"log"
// 	"encoding/json"
// 	"math/rand"
// 	"net/http"
// 	"strconv"
// 	"github.com/gorilla/mux"
// 	"github.com/gorilla/handlers"

// )


// var products []Product

// func getProducts(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(products)
// }

// func deleteProduct(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "application/json")
// 	params :=mux.Vars(r)
// 	for index, item := range products{

// 		if item.ID == params["id"]{
// 			products = append(products[:index], products[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(products)
// }

// func getProduct(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "application/json")
// 	params :=mux.Vars(r)

// 	for _, item := range products{
// 		if item.ID == params["id"]{
// 			json.NewEncoder(w).Encode(item)
// 			return
// 		}
// 	}
// }

// func createProduct(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "application/json")
// 	var product Product
// 	_ = json.NewDecoder(r.Body).Decode(&product)
// 	product.ID = strconv.Itoa(rand.Intn(10000000000000))
// 	products = append(products, product)
// 	json.NewEncoder(w).Encode(products)
// 	return
// }

// func updateProduct(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	for index, item := range products{
// 		if item.ID == params["id"] {
// 			products = append(products[:index], products[index+1:]... )	
// 			var product Product 	
// 			_ = json.NewDecoder(r.Body).Decode(&product)
// 			product.ID = params["id"]
// 			products = append(products, product)
// 			json.NewEncoder(w).Encode(products)
// 			return
// 		}
// 	}
// }

// func main(){
// 	r := mux.NewRouter()

// 	products = append(products, Product{ID: "1", Name: "apple", Description: "delicious", Price: 10, Category: &Category{ID: "1", Type: "food"}})
// 	products = append(products, Product{ID: "2", Name: "banna", Description: "delicious", Price: 20, Category: &Category{ID: "1", Type: "food"}})
	
// 	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
//     originsOk := handlers.AllowedOrigins([]string{"*"})
//     methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
// 	(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))

// 	fmt.Printf("Starting server at port 8000\n")
// 	log.Fatal(http.ListenAndServe("8000",r))
// }