package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Strucutures and collections
type Customer struct{
	ID             int
    Name           string
    Email          string
    LastLoginTime  time.Time
    PagesVisited   []string
    ItemsInCart    []string		   
}

type Category struct {
	ID 	 int 	`json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID 	  int 	  `json:"id"`
	Name  string  `json:"name"`
	Price int 	  `json:"price"`
}

type Order struct {
    ID           int
	CustomerID   int
    OrderDate    string
    TotalAmount  float64
	CreatedAt   time.Time
    SessionID   string
}

type Payment struct {
    ID            int     `json:"id"`
    CustomerID    int     `json:"customer_id"`
    Amount        float64 `json:"amount"`
    PaymentMethod string  `json:"payment_method"`
}

var customers  []Customer		
var products   []Product			  
var categories []Category		 
var orders     []Order 
var payments   []Payment
///////////////////////////////////////////////////////////////////////////

var db *sql.DB    //DB instance 

func main() {

	var err error
	db, err = sql.Open("mysql", "root:systemDesign2023$@tcp(host:3306)/ecommerce")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	customerID := 1 
    orders, err := getOrdersByCustomerID(db, customerID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Orders for customer %d:\n", customerID)
    for _, order := range orders {
        fmt.Printf("%d - Total amount: %f, Created at: %s, Session ID: %s\n", order.ID, order.TotalAmount, order.CreatedAt, order.SessionID)
    }

    // Get customer information
    customer, err := getCustomerByID(db, customerID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Customer information:\nName: %s\nEmail: %s\nLast login time: %s\nPages visited: %v\nItems in cart: %v\n",
        customer.Name, customer.Email, customer.LastLoginTime, customer.PagesVisited, customer.ItemsInCart)


	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	customers := getCustomers()
	products := getProducts()

	print(customers)
	print(orders)


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

func getOrdersByCustomerID(db *sql.DB, customerID int) ([]Order, error) {
    rows, err := db.Query("SELECT id, customer_id, total_amount, created_at, session_id FROM orders WHERE customer_id = ?", customerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    orders := []Order{}
    for rows.Next() {
        var order Order
        err := rows.Scan(&order.ID, &order.CustomerID, &order.TotalAmount, &order.CreatedAt, &order.SessionID)
        if err != nil {
            return nil, err
        }
        orders = append(orders, order)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return orders, nil
}

func getCustomerByID(db *sql.DB, customerID int) (Customer, error) {
    var customer Customer
    err := db.QueryRow("SELECT id, name, email, last_login_time, pages_visited, items_in_cart FROM customers WHERE id = ?", customerID).
        Scan(&customer.ID, &customer.Name, &customer.Email, &customer.LastLoginTime, &customer.PagesVisited, &customer.ItemsInCart)
    if err != nil {
        return Customer{}, err
    }
    return customer, nil
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

	_, err := strconv.Atoi(id)
	if err !=nil{
		panic(err)
	}
	orderFactory(99,"Vazgen","22/02/2023",123.4)
}

func orderFactory(customer_id int, customer_name string, order_date string, total_amount float64){
	var id int 
	for _,customer := range customers{
		if(customer.ID == customer_id){
			id = len(orders) + 1
			orders =  append(orders,Order{ID: id, CustomerID: customer_id,OrderDate: order_date,TotalAmount: total_amount})
			break
		}
	}
}

func handleCustomers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getCustomers()
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getCustomers() []Customer {
	rows, err := db.Query("SELECT * FROM customers")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID)
		if err != nil {
			panic(err.Error())
		}
		customers = append(customers, customer)
	}

	return customers
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

func handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getProducts()
	case http.MethodPost:
		createProduct(w, r)
	case http.MethodPut:
		updateProduct(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getProducts() []Product {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, product)
	}

	return products
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

func handleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getOrders()
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

func getOrders() ([]Order) {
    orders := []Order{}

    rows, err := db.Query("SELECT id, customer_id,  order_date, total_amount FROM orders")
    if err != nil {
        return nil
    }

    defer rows.Close()

    for rows.Next() {
        var order Order

        err := rows.Scan(&order.ID, &order.CustomerID,  &order.OrderDate, &order.TotalAmount)
        if err != nil {
            return nil
        }

        orders = append(orders, order)
    }

    return orders
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