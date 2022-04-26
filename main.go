package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Vehicle struct {
	Id int
	Make string
	Model string
	Price int
}

var vehicles = []Vehicle{
	{1, "Toyota", "Corolla", 1000},
	{2, "Toyota", "Camry", 2000},
	{3, "Honda", "Civic", 1500},
}

func returnAllCars(w http.ResponseWriter, r *http.Request) {
	// send statusOK to header. 200
	w.WriteHeader(http.StatusOK)
	// encode into json write to vehicles
	json.NewEncoder(w).Encode(vehicles)
}

func returnCarsByBrand(w http.ResponseWriter, r *http.Request) {
	// vars equal to the request entered into the router
	// vars is a map of keyvalue pairs
	vars := mux.Vars(r)
	//carsM (make of cars) = a slice of the requests with keyword "make"
	carM := vars["make"]
	cars := &[]Vehicle{}
	for _, car := range vehicles {
		if car.Make == carM {
			*cars = append(*cars, car)
		}
		// reaching through memory address of Vehicles to
		//capture the value and appened to the list of values, the new car value
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cars)
}

func returnCarsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("unable to convert to string")
	}
	for _, car := range vehicles {
		if car.Id == carId {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(car)
		}
	}
}

func updateCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("unable to convert to string")
	}
	var updateCar Vehicle
	json.NewDecoder(r.Body).Decode(&updateCar)
	for k, v := range vehicles {
		if v.Id == carId {
			vehicles = append(vehicles[:k], vehicles[k+1:]... )
			vehicles = append(vehicles, updateCar)
		}
	}
	json.NewEncoder(w).Encode(vehicles)
	w.WriteHeader(http.StatusOK)
}

func createCar(w http.ResponseWriter, r *http.Request) {
	var newCar Vehicle
	json.NewDecoder(r.Body).Decode(&newCar)
	vehicles = append(vehicles, newCar)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

func removeCarById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carId, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("unable to convert to string")
	}
	for k, v := range vehicles {
		if v.Id == carId {
			vehicles = append(vehicles[:k], vehicles[k+1:]...)
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vehicles)
}

func main() {
	// create new route. ScrictSlash true = strict about trailing slash on endpoint.
	router := mux.NewRouter().StrictSlash(true)
	// Create endpoint that runs function on router. Method GET
	router.HandleFunc("/cars", returnAllCars).Methods("GET")
	router.HandleFunc("/cars/make/{make}", returnCarsByBrand).Methods("GET")
	router.HandleFunc("/cars/{id}", returnCarsById).Methods("GET")
	router.HandleFunc("/cars/{id}", updateCar).Methods("PUT")
	router.HandleFunc("/cars", createCar).Methods("POST")
	router.HandleFunc("/cars/{id}", removeCarById).Methods("DELETE")
	// Not using default "nil" router. Using router we have created.
	log.Fatal(http.ListenAndServe(":8080", router))
}