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

func ShowFoods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /foods: showProduce")
	json.NewEncoder(w).Encode(Produce)
}

func DeleteFood(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE /food: DeleteFood")

	vars := mux.Vars(r)
	code := vars["code"]

	var deletedFood bool = false
	for index, food := range Produce {
		if food.Code == code {
			Produce = append(Produce[:index], Produce[index+1:]...)
			deletedFood = true
		}
	}

	if deletedFood {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func requestHandler() {
	produceRouter := mux.NewRouter().StrictSlash(true)
	produceRouter.HandleFunc("/foods", ShowFoods)
	produceRouter.HandleFunc("/food/{code}", DeleteFood).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", produceRouter))
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
