package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Food struct {
	ProduceCode string
	Name        string
	UnitPrice   float64
}

var Produce []Food

func ShowFoods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /food: showFoods")
	json.NewEncoder(w).Encode(Produce)
}

func ShowFood(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /food/{code}: showFood")
	handlerChannel := make(chan Food)

	go func() {
		vars := mux.Vars(r)
		code := vars["code"]
		var foundFood Food
		for index, food := range Produce {
			if food.ProduceCode == code {
				foundFood = Produce[index]
				break
			}
		}
		handlerChannel <- foundFood
	}()

	foundFood := <-handlerChannel
	if foundFood.ProduceCode != "" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(foundFood)
}

func CreateFood(w http.ResponseWriter, r *http.Request) {
	fmt.Println("POST /food: CreateFood")
	handlerChannel := make(chan Food)

	go func() {
		var createdFood Food
		reqBody, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(reqBody, &createdFood)
		// update our global Articles array to include
		// our new Article
		Produce = append(Produce, createdFood)
		handlerChannel <- createdFood
	}()

	foodCreated := <-handlerChannel

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(foodCreated)
}

func DeleteFood(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DELETE /food: DeleteFood")
	handlerChannel := make(chan bool)

	go func() {
		vars := mux.Vars(r)
		code := vars["code"]
		var deletedFood bool = false
		for index, food := range Produce {
			if food.ProduceCode == code {
				Produce = append(Produce[:index], Produce[index+1:]...)
				deletedFood = true
				break
			}
		}
		handlerChannel <- deletedFood
	}()

	foodDeleted := <-handlerChannel
	if foodDeleted {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func requestHandler() {
	produceRouter := mux.NewRouter().StrictSlash(true)
	produceRouter.HandleFunc("/food/{code}", ShowFood)
	produceRouter.HandleFunc("/food", ShowFoods)
	produceRouter.HandleFunc("/foods", CreateFood).Methods("POST")
	produceRouter.HandleFunc("/food/{code}", DeleteFood).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", produceRouter))
}

func main() {
	// Instantitate in-memory database array of food.
	Produce = []Food{
		Food{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: 3.64},
		Food{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: 2.99},
		Food{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: 0.79},
		Food{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: 3.59},
	}
	requestHandler()
}
