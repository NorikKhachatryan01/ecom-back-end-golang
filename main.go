package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Customer struct{
	ID 		 int    `json:"id"`
	Username string `json:"user_name"`
	Basket   		[]Order		   
}

var customers []Customer
 	

// Product Model
type Product struct {
	ID 	  int 	  `json:"id"`
	Name  string  `json:"name"`
	Price int `json:"price"`
}

var products []Product


// Category Model
type Category struct {
	ID 	 int 	`json:"id"`
	Name string `json:"name"`
}

var categories []Category

// Order Model
type Order struct {
	ID 	       int    `json:"id"`
	ProductID  int    `json:"product_id"`
	CustomerID int    `json:"customer_id"`
}

var orders []Order

func main() {


	r := mux.NewRouter()



	customers  = append(customers,Customer{ID: 99, Username: "Customer_1",Basket: []Order{} })
	products = append(products,Product{ID: 1, Name: "apple", Price: 100})
	products = append(products,Product{ID: 2, Name: "test1", Price: 200})
	products = append(products,Product{ID: 3, Name: "test2", Price: 300})
	products = append(products,Product{ID: 4, Name: "test3", Price: 400})
	products = append(products,Product{ID: 5, Name: "test4", Price: 500})
	products = append(products,Product{ID: 6, Name: "test5", Price: 600})
	products = append(products,Product{ID: 7, Name: "test6", Price: 700})
	products = append(products,Product{ID: 8, Name: "test7", Price: 800})
	orders = append(orders, Order{})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

  

	for _, product := range products {
        fmt.Println(product)
    }

	
	r.HandleFunc("/products", handleProducts)
	r.HandleFunc("/products/{id}", handlePurchase).Methods("POST")
	r.HandleFunc("/categories", handleCategories)
	r.HandleFunc("/orders", handleOrders)
	r.HandleFunc("/customer", handleCustomers)
  (http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}


func handlePurchase(w http.ResponseWriter, r *http.Request){

	vars :=  mux.Vars(r)
	id, ok :=  vars["id"]

	if !ok {
		fmt.Println("id missing in paremters")
	}


	i, err := strconv.Atoi(id)
	if err !=nil{
		panic(err)
	}
	orderFactory(99,i)
}

func orderFactory(customer_id int,product_id int){
	var id int 
	for _,customer := range customers{
		if(customer.ID == customer_id){
			id = len(orders) + 1
			orders =  append(orders,Order{ID: id, ProductID: product_id, CustomerID: customer_id})
			
			break
		}
	}
}

func handleCustomers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCustomers(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts(w, r)
	case http.MethodPost:
		createProduct(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	product.ID = len(products) + 1
	products = append(products, product)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	for i, item := range products {
		if item.ID == product.ID {
			products[i] = product
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func handleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCategories(w, r)
	case http.MethodPost:
		createCategory(w, r)
	case http.MethodPut:
		updateCategory(w, r)
	case http.MethodDelete:
		deleteCategory(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	json.NewDecoder(r.Body).Decode(&category)

	category.ID = len(categories) + 1
	categories = append(categories, category)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	json.NewDecoder(r.Body).Decode(&category)

	for i, item := range categories {
		if item.ID == category.ID {
			categories[i] = category
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	json.NewDecoder(r.Body).Decode(&category)

	for i, item := range categories {
		if item.ID == category.ID {
			categories = append(categories[:i], categories[i+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOrders(w, r)
	case http.MethodPost:
		createOrder(w, r)
	case http.MethodPut:
		updateOrder(w, r)
	case http.MethodDelete:
		deleteOrder(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)

	order.ID = len(orders) + 1
	orders = append(orders, order)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)

	for i, item := range orders {
		if item.ID == order.ID {
			orders[i] = order
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)

	for i, item := range orders {
		if item.ID == order.ID {
			orders = append(orders[:i], orders[i+1:]...)
			break
		}
	}

	w.WriteHeader(http.StatusNoContent)
}