package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/go-sql-driver/mysql"
)

// Order represents the structure of an order.
type Order struct {
	ID           string      `json:"id"`
	Status       string      `json:"status"`
	Items        []OrderItem `json:"items"`
	Total        float64     `json:"total"`
	CurrencyUnit string      `json:"currencyUnit"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}

// OrderItem represents the structure of an order item.
type OrderItem struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Handle requests for creating an order.
	r.Post("/orders", handleCreateOrder)

	// Handle requests for updating an order.
	r.Put("/orders/{id}", handleUpdateOrder)

	// Handle requests for fetching orders based on all the fields of the order in a sorted and filtered way.
	r.Get("/orders", handleGetOrders)

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Connects to the MySQL database.
func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/orders")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// Parse the order data from the request body.
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Validate the order data.
	if order.ID == "" || len(order.Items) == 0 || order.Total <= 0 || order.CurrencyUnit == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Insert the order into the database.
	db, err := connectToDB()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	insertOrderQuery := "INSERT INTO orders (id, status, total, currency_unit, items) VALUES (?, ?, ?, ?, ?)"
	orderItems, _ := json.Marshal(order.Items)
	_, err = db.Exec(insertOrderQuery, order.ID, order.Status, order.Total, order.CurrencyUnit, orderItems)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Return the created order ID as the response.
	type CreateOrderResponse struct {
		ID string `json:"id"`
	}
	response := CreateOrderResponse{ID: order.ID}
	json.NewEncoder(w).Encode(response)
}

func handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the order to update from the URL path.
	id := chi.URLParam(r, "id")

	// Parse the updated order data from the request body.
	var updateData struct {
		Status string `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Update the order in the database.
	db, err := connectToDB()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	updateOrderQuery := "UPDATE orders SET status = ? WHERE id = ?"
	_, err = db.Exec(updateOrderQuery, updateData.Status, id)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Return a success response.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Status Successfully Updated"))

}

// Handle requests for fetching orders based on all the fields of the order in a sorted and filtered way.
func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	// Get database connection
	db, err := connectToDB()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute SELECT statement
	rows, err := db.Query("SELECT id, status, items, total, currency_unit, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s'), DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') FROM orders")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create slice of orders
	orders := make([]Order, 0)

	// Iterate over rows and scan into orders slice
	for rows.Next() {
		var order Order
		var createdAtStr, updatedAtStr string
		var items []byte

		if err := rows.Scan(&order.ID, &order.Status, &items, &order.Total, &order.CurrencyUnit, &createdAtStr, &updatedAtStr); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Parse createdAt and updatedAt values into time.Time
		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		order.CreatedAt = createdAt
		order.UpdatedAt = updatedAt

		// Unmarshal items field into slice of Item structs
		if err := json.Unmarshal(items, &order.Items); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		orders = append(orders, order)
	}

	// Marshal orders slice into JSON
	orderData, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Set Content-Type header and write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(orderData)
}
