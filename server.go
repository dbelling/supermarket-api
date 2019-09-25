package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Food struct {
	Code  string
	Name  string
	Price string
}

var Produce []Food

func showProduce(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /produce: showProduce")
	json.NewEncoder(w).Encode(Produce)
}

func requestHandler() {
	http.HandleFunc("/produce", showProduce)
	log.Fatal(http.ListenAndServe(":9000", nil))
}

func main() {
	// Instantitate in-memory database array of food.
	Produce = []Food{
		Food{Code: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", Price: "$3.46"},
		Food{Code: "E5T6-9UI3-TH15-QR88", Name: "Peach", Price: "$2.99"},
		Food{Code: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", Price: "$0.79"},
		Food{Code: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", Price: "$3.59"},
	}
	requestHandler()
}
