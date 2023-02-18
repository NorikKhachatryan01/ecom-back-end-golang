package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Above is provided the main components strucutres 

//Customer Model
type Customer struct{
	ID 		 int    `json:"id"`
	Username string `json:"user_name"`
	Basket   		[]Order		   
}
var customers []Customer			//Customers list 
 	
// Product Model
type Product struct {
	ID 	  int 	  `json:"id"`
	Name  string  `json:"name"`
	Price int 	  `json:"price"`
}
var products []Product			  //Products list

// Category Model
type Category struct {
	ID 	 int 	`json:"id"`
	Name string `json:"name"`
}
var categories []Category		 //Categories list 	

// Order Model
type Order struct {
	ID 	       int    `json:"id"`
	ProductID  int    `json:"product_id"`
	CustomerID int    `json:"customer_id"`
}  
var orders []Order   			//Orders list

//Payment Model
type Payment struct {
    ID            int     `json:"id"`
    CustomerID    int     `json:"customer_id"`
    Amount        float64 `json:"amount"`
    PaymentMethod string  `json:"payment_method"`
}
var payments []Payment

func main() {
	r := mux.NewRouter()
	customers  = append(customers,Customer{ID: 99, Username: "Customer_1",Basket: []Order{} })
	products = append(products,Product{ID: 1, Name: "apple", Price: 100})
	products = append(products,Product{ID: 2, Name: "test1", Price: 200})
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
	r.HandleFunc("/payments", handlePayments).Methods("POST")
  (http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}


func handlePayments(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        var payment Payment
        json.NewDecoder(r.Body).Decode(&payment)

        // Validate the payment
        if payment.Amount <= 0 {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprint(w, "Payment amount must be greater than 0")
            return
        }

        // Process the payment
        if payment.PaymentMethod == "credit_card" {
            // process credit card payment
            fmt.Println("Processing credit card payment...")
        } else if payment.PaymentMethod == "paypal" {
            // process PayPal payment
            fmt.Println("Processing PayPal payment...")
        } else {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprint(w, "Invalid payment method")
            return
        }

        // Add the payment to the database
        payment.ID = len(payments) + 1
        payments = append(payments, payment)

        // Return success response
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "Payment processed successfully")
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
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