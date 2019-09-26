package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestShowFoods(t *testing.T) {
	req, err := http.NewRequest("GET", "/foods", nil)
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
