package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCreateOrder(t *testing.T) {
	// Set up a mock request with some JSON data
	reqBody := `{

        "id": "12345002",

        "status": "PENDING_INVOICE",

        "items": [{

                "id": "123456",

                "description": "a product description",

                "price": 12.40,

                "quantity": 1

        }],

        "total": 12.40,

        "currencyUnit": "USD"

}`
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's output
	rr := httptest.NewRecorder()

	// Call the handler function, passing in the response recorder and request
	handleCreateOrder(rr, req)

	// Check the response body is what we expect
	expectedBody := `{"id":"12345002"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestHandleUpdateOrder(t *testing.T) {
	// Set up a mock request with some JSON data
	reqBody := `{
		"status": "PAID"
	}`
	req, err := http.NewRequest("PUT", "/orders/1", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's output
	rr := httptest.NewRecorder()

	// Call the handler function, passing in the response recorder and request
	handleUpdateOrder(rr, req)

	// Check the status code is what we expect
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	// Check the response body is what we expect
	expectedBody := `Status Successfully Updated`
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedBody)
	}
}

func TestHandleGetOrders(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetOrders)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}
