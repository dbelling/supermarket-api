package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowFoods(t *testing.T) {
	req, err := http.NewRequest("GET", "/food", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowFoods)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestShowFood(t *testing.T) {
	req, err := http.NewRequest("GET", "/food/123-ABC", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ShowFood)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestDeleteFood(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/food/123-ABC", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteFood)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestCreateFood(t *testing.T) {
	var foodPayload = []byte(`{"ProduceCode":"ABCD-1234-DEFG-5678","Name":"Cucumber","UnitPrice":1.99}`)

	req, err := http.NewRequest("POST", "/foods", bytes.NewBuffer(foodPayload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateFood)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	var m map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &m)
	if m["ProduceCode"] != "ABCD-1234-DEFG-5678" {
		t.Errorf("Expected ProduceCode to be 'ABCD-1234-DEFG-5678'. Got '%v'", m["ProduceCode"])
	}
	if m["Name"] != "Cucumber" {
		t.Errorf("Expected Name to be 'Cucumber'. Got '%v'", m["Name"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["UnitPrice"] != 1.99 {
		t.Errorf("Expected Price to be '1.99'. Got '%v'", m["UnitPrice"])
	}
}
