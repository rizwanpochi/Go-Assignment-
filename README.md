# Documentation for the code

## Introduction

This code represents an HTTP server written in Go, which has 3 endpoints for creating, updating and retrieving orders from a MySQL database.

## Dependencies

The code depends on the following packages:

- `database/sql` for working with databases and SQL statements
- `encoding/json` for encoding and decoding JSON data
- `log` for logging messages
- `net/http` for handling HTTP requests and responses
- `time` for working with date and time values
- `github.com/go-chi/chi` for working with the HTTP router
- `github.com/go-chi/chi/middleware` for adding middleware to the router
- `github.com/go-sql-driver/mysql` for working with the MySQL database

## Types

### Order

The `Order` type is a struct that represents an order. It has the following fields:

- `ID` (string): the ID of the order
- `Status` (string): the status of the order
- `Items` ([]OrderItem): a slice of `OrderItem` structs that represent the items in the order
- `Total` (float64): the total cost of the order
- `CurrencyUnit` (string): the currency unit used for the total cost
- `CreatedAt` (time.Time): the date and time the order was created
- `UpdatedAt` (time.Time): the date and time the order was last updated

### OrderItem

The `OrderItem` type is a struct that represents an item in an order. It has the following fields:

- `ID` (string): the ID of the item
- `Description` (string): a description of the item
- `Price` (float64): the price of the item
- `Quantity` (int): the quantity of the item

## Functions

### `main()`

The `main()` function sets up the HTTP server and defines the routes for handling requests. The following routes are defined:

- `POST /orders`: creates a new order in the database
- `PUT /orders/{id}`: updates the status of an existing order in the database
- `GET /orders`: retrieves a sorted and filtered list of orders from the database

### `connectToDB()`

The `connectToDB()` function returns a connection to the MySQL database.

### `handleCreateOrder(w http.ResponseWriter, r *http.Request)`

The `handleCreateOrder()` function handles the `POST /orders` route. It parses the order data from the request body, validates it, and inserts the order into the database. If the order data is invalid or there is an error inserting the order into the database, an error response is sent. If the order is created successfully, the ID of the created order is returned in the response.

### `handleUpdateOrder(w http.ResponseWriter, r *http.Request)`

The `handleUpdateOrder()` function handles the `PUT /orders/{id}` route. It gets the ID of the order to update from the URL path, parses the updated order data from the request body, and updates the order in the database. If there is an error updating the order in the database, an error response is sent. If the order is updated successfully, a success response is sent.

### `handleGetOrders(w http.ResponseWriter, r *http.Request)`

This API endpoint allows you to fetch orders from the database based on all the fields of the order in a sorted and filtered way.
Endpoint URL
GET /orders

Query Parameters
sort: Specifies the field to sort_by. 

Valid values are:
id: Order ID
status: Order status	
items: Order items
total: Order total
currency_unit: Currency unit of the order
created_at: Time when the order was created
updated_at: Time when the order was last updated
order: Specifies the sort order. Valid values are:

asc: Sort in ascending order (default)
desc: Sort in descending order
id: Filter orders by order ID. The value should be an integer.

status: Filter orders by order status. The value should be a string.

item_name: Filter orders by the name of the items in the order. The value should be a string.

item_price: Filter orders by the price of the items in the order. The value should be a float.

item_quantity: Filter orders by the quantity of the items in the order. The value should be an integer.

total: Filter orders by the total amount of the order. The value should be a float.

currency_unit: Filter orders by the currency unit of the order. The value should be a string.

created_at: Filter orders by the time when the order was created. The value should be a string in the format YYYY-MM-DD HH:MM:SS.

updated_at: Filter orders by the time when the order was last updated. The value should be a string in the format YYYY-MM-DD HH:MM:SS.

Example Requests
Fetch all orders, sorted by created_at in descending order:
GET /orders?sort=created_at&order=desc

Fetch orders with status "shipped" and total greater than or equal to 100, sorted by created_at in ascending order:
GET /orders?status=shipped&total_gte=100&sort=created_at&order=asc

Response Format
The API returns a JSON object containing an array of order objects that match the specified filters and sorted according to the specified sorting criteria.

Each order object has the following properties:

id: Order ID (integer)
status: Order status (string)
items: Array of order items, where each item has the following properties:
name: Item name (string)
price: Item price (float)
quantity: Item quantity (integer)
total: Total amount of the order (float)
currency_unit: Currency unit of the order (string)
created_at: Time when the order was created (string in the format YYYY-MM-DD HH:MM:SS)
updated_at: Time when the order was last updated (string in the format YYYY-MM-DD HH:MM:SS)

## Query Parameters

| Parameter | Type | Description |
| --- | --- | --- |
| `id` | integer | Filters orders by ID. |
| `status` | string | Filters orders by status. |
| `item_name` | string | Filters orders by item name. |
| `item_quantity` | integer | Filters orders by item quantity. |
| `item_price` | number | Filters orders by item price. |
| `sort_by` | string | Sorts orders by the given field. Valid values are `id`, `status`, `total`, `created_at`, and `updated_at`. |
| `sort_order` | string | Sets the sort order. Valid values are `asc` and `desc`. |
| `page` | integer | Specifies the page number. |
| `page_size` | integer | Specifies the number of orders per page. |


To filter orders by status, we can use the status parameter

To filter orders by date range, we can use the start_date and end_date parameters:
http://localhost:8080/orders?start_date=2022-01-01&end_date=2023-01-31
This will retrieve all orders that were created between January 1st, 2022 and January 31st, 2023.

To sort orders by a particular field, we can use the sort_by parameter
http://localhost:8080/orders?sort_by=total

To sort orders in descending order, we can use the sort_order parameter:
http://localhost:8080/orders?sort_by=total&sort_order=desc	

We can also combine multiple parameters to create more complex queries
http://localhost:8080/orders?status=shipped&sort_by=total&sort_order=desc&page=2&page_size=10

Setting up Database in MySql

Install MySQL: You can download the MySQL community server from the official website of MySQL.

Open the MySQL client and connect to the MySQL server with the following command:
css
mysql -u root -p

Navigate to the directory where you saved the orders.sql file using the following command:
cd /path/to/sql/file/

Import the orders.sql file into the MySQL server with the following command:
source orders.sql;

This will create the necessary tables for the orders database.

Now that the database is set up, you can update the database connection string in the Go code to connect to the 'orders' database:
db, err := sql.Open("mysql", "orders_user:password@tcp(localhost:3306)/orders")
