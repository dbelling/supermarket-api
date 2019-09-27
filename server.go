package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidProduceCode = errors.New("invalid produce code")
	ErrInvalidName        = errors.New("invalid name")
	ErrInvalidUnitPrice   = errors.New("invalid unit price")
)

type Food struct {
	ProduceCode string
	Name        string
	UnitPrice   string
}

type FoodValidation interface {
	Validate(r *http.Request) error
}

var Produce []Food

func (f Food) Validate(r *http.Request) error {

	if len(f.ProduceCode) != 19 {
		return ErrInvalidProduceCode
	}

	chars := []rune(f.ProduceCode)

	if !govalidator.IsAlphanumeric(string(chars[0:3])) {
		return ErrInvalidProduceCode
	}

	if !govalidator.IsAlphanumeric(string(chars[5:8])) {
		return ErrInvalidProduceCode
	}

	if !govalidator.IsAlphanumeric(string(chars[10:13])) {
		return ErrInvalidProduceCode
	}

	if !govalidator.IsAlphanumeric(string(chars[15:18])) {
		return ErrInvalidProduceCode
	}

	if govalidator.IsNull(f.Name) {
		return ErrInvalidName
	}

	if !govalidator.IsFloat(f.UnitPrice) {
		return ErrInvalidUnitPrice
	}

	return nil
}

func Validate(r *http.Request, v FoodValidation) error {
	return v.Validate(r)
}

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

	reqBody, _ := ioutil.ReadAll(r.Body)
	x := bytes.TrimLeft(reqBody, " \t\r\n")
	isArray := len(x) > 0 && x[0] == '['

	foodsToCreate := make([]Food, 0)
	if isArray {
		decoder := json.NewDecoder(bytes.NewBufferString(string(reqBody)))
		err := decoder.Decode(&foodsToCreate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, food := range foodsToCreate {
			err := Validate(r, food)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			foodsToCreate = append(foodsToCreate, food)
		}
	} else {
		var createdFood Food
		json.Unmarshal(reqBody, &createdFood)
		err := Validate(r, createdFood)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		foodsToCreate = append(foodsToCreate, createdFood)
	}


	handlerChannel := make(chan []Food)
	go func(foods []Food) {
		for _, food := range foods {
			Produce = append(Produce, food)
		}
		handlerChannel <- foods
	}(foodsToCreate)

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
		Food{ProduceCode: "A12T-4GH7-QPL9-3N4M", Name: "Lettuce", UnitPrice: "3.64"},
		Food{ProduceCode: "E5T6-9UI3-TH15-QR88", Name: "Peach", UnitPrice: "2.99"},
		Food{ProduceCode: "YRT6-72AS-K736-L4AR", Name: "Green Pepper", UnitPrice: "0.79"},
		Food{ProduceCode: "TQ4C-VV6T-75ZX-1RMR", Name: "Gala Apple", UnitPrice: "3.59"},
	}
	requestHandler()
}
