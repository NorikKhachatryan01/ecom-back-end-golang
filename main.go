package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
 
)

//Above is provided the main components strucutres 

//Customer Model
type Customer struct{
	ID 		 int      `json:"id"`
	Username string   `json:"user_name"`
	Session  *Session `json:"session"`
	Basket   []Order  `json:"basket"`
}
var customers []Customer			//Customers list 
 	
// Product Model
type Product struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name        string             `json:"name" bson:"name"`
    Description string             `json:"description" bson:"description"`
    Price       int                `json:"price" bson:"price"`
    Available   bool               `json:"available" bson:"available"`
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

// Session Model
type Session struct {
	LastLogin     time.Time   `json:"last_login"`
	PagesVisited  []string    `json:"pages_visited"`
	ShoppingCart  []Order     `json:"shopping_cart"`
}


func getMongoClient() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

var (
	customerCollection *mongo.Collection
	productCollection  *mongo.Collection
	categoryCollection *mongo.Collection
	orderCollection    *mongo.Collection
	paymentCollection  *mongo.Collection
)

func initMongo() {
	// Get MongoDB client
	client, err := getMongoClient()
	if err != nil {
		log.Fatal(err)
	}

	// Get database and collections
	db := client.Database("myapp")
	customerCollection = db.Collection("customers")
	productCollection = db.Collection("products")
	categoryCollection = db.Collection("categories")
	orderCollection = db.Collection("orders")
	paymentCollection = db.Collection("payments")
}


collection := client.Database("mydb").Collection("products")
func main() {
	initMongo()
	


	r := mux.NewRouter()

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
	r.HandleFunc("/order_history", handleOrderHistory)
	r.HandleFunc("/customer", handleCustomers)
	r.HandleFunc("/payments", handlePayments).Methods("POST")
  (http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}

func handleOrderHistory(w http.ResponseWriter, r *http.Request) {
	// Retrieve all orders from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := orderCollection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error retrieving orders: %v", err)
		return
	}
	defer cursor.Close(ctx)

	// Decode orders into a slice
	var orders []Order
	if err := cursor.All(ctx, &orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error decoding orders: %v", err)
		return
	}

	// Return orders as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding orders: %v", err)
		return
	}
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


func handlePurchase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Get the product with the given ID
	product, err := getProductByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Product with ID %d not found", id)
		return
	}

	// Get the customer from the session information
	session, err := getSession(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error getting session information: %v", err)
		return
	}

	// Check if the customer has already placed an order for this product
	if hasCustomerOrderedProduct(session.CustomerID, id) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Customer has already ordered this product")
		return
	}

	// Create a new order
	order := Order{
		ID:         getNextOrderID(),
		ProductID:  id,
		CustomerID: session.CustomerID,
	}

	// Add the order to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := orderCollection.InsertOne(ctx, order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error inserting order into database: %v", err)
		return
	}

	// Update the shopping cart in the session information
	session.Cart = append(session.Cart, id)
	if err := saveSession(session, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error saving session information: %v", err)
		return
	}

	// Return the ordered product as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error encoding product: %v", err)
		return
	}
}

// getSession gets the session information from the request cookies.
func getSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil, err
	}
	sessionID := cookie.Value

	// Get the session information from the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var session Session
	if err := sessionCollection.FindOne(ctx, bson.M{"_id": sessionID}).Decode(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

// saveSession saves the session information to the response cookies.
func saveSession(session *Session, w http.ResponseWriter) error {
	// Update the session information in the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := sessionCollection.ReplaceOne(ctx, bson.M{"_id": session.ID}, session); err != nil {
		return err
	}

	// Set the session ID as a cookie in the response
	cookie := http.Cookie{
		Name:     "session",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour),
	}
	http.SetCookie(w, &cookie)

	return nil
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
		// Retrieve all customers from the database
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cursor, err := customerCollection.Find(ctx, bson.M{})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error retrieving customers: %v", err)
			return
		}
		defer cursor.Close(ctx)

		// Decode customers into a slice
		var customers []Customer
		if err := cursor.All(ctx, &customers); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error decoding customers: %v", err)
			return
		}

		// Return customers as JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(customers); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error encoding customers: %v", err)
			return
		}

	case http.MethodPost:
		// Decode new customer from request body
		var customer Customer
		if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Invalid request payload")
			return
		}

		// Insert new customer into the database
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		result, err := customerCollection.InsertOne(ctx, customer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error inserting customer: %v", err)
			return
		}

		// Return the ID of the new customer as JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(struct{ ID string }{ID: result.InsertedID.(primitive.ObjectID).Hex()}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error encoding ID: %v", err)
			return
		}
	}
}

func handleProducts(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getProducts(w, r, collection)
    case http.MethodPost:
        createProduct(w, r, collection)
    case http.MethodPut:
        updateProduct(w, r, collection)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
	// Get all customers from database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := customerCollection.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	// Iterate over the cursor and encode the results
	var customers []Customer
	for cur.Next(ctx) {
		var customer Customer
		err := cur.Decode(&customer)
		if err != nil {
			log.Fatal(err)
		}
		customers = append(customers, customer)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Set the response header and encode the results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

func getProducts(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
    var products []Product

    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: %s", err.Error())
        return
    }
    defer cursor.Close(context.Background())

    for cursor.Next(context.Background()) {
        var product Product
        err := cursor.Decode(&product)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "Error: %s", err.Error())
            return
        }
        products = append(products, product)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request, collection *mongo.Collection) {
    var product Product
    json.NewDecoder(r.Body).Decode(&product)

    // Set a new ID for the product
    product.ID = primitive.NewObjectID()

    // Insert the product into the MongoDB collection
    _, err := collection.InsertOne(context.Background(), product)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Error: %s", err.Error())
        return
    }
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