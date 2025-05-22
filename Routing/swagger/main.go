package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-openapi/runtime/middleware"
)

// Pet represents a pet object
// swagger:model
type Pet struct {
	// The ID of the pet
	// example: 123
	ID int64 `json:"id"`

	// The name of the pet
	// required: true
	// example: Fido
	Name string `json:"name"`

	// The status of the pet
	// enum: [available, pending, sold]
	Status string `json:"status,omitempty"`
}

// NewPet represents the data needed to create a new pet
// swagger:model
type NewPet struct {
	// The name of the pet
	// required: true
	// example: Buddy
	Name string `json:"name"`

	// The status of the pet
	// enum: [available, pending, sold]
	Status string `json:"status,omitempty"`
}

// @Summary Add a new pet to the store
// @Description Add a new pet
// @Accept  json
// @Produce json
// @Param body body NewPet true "Pet object that needs to be added to the store"
// @Success 200 {object} Pet "Successful operation"
// @Failure 405 {string} string "Invalid input"
// @Router /pets [post]
func addPetHandler(w http.ResponseWriter, r *http.Request) {
	var newPet NewPet
	if err := json.NewDecoder(r.Body).Decode(&newPet); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// In a real application, you would save the pet to a database
	pet := Pet{ID: 1, Name: newPet.Name, Status: newPet.Status}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

// @Summary Get pet by ID
// @Description Get pet by ID
// @Produce json
// @Param petId path int64 true "ID of pet to return"
// @Success 200 {object} Pet "Successful operation"
// @Failure 400 {string} string "Invalid ID supplied"
// @Failure 404 {string} string "Pet not found"
// @Router /pets/{petId} [get]
func getPetHandler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	petIDStr := vars.Get("petId")
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}
	// In a real application, you would retrieve the pet from a database
	if petID == 1 {
		pet := Pet{ID: petID, Name: "Fido", Status: "available"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pet)
	} else {
		http.Error(w, "Pet not found", http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pets", addPetHandler)
	mux.HandleFunc("/pets/", getPetHandler) // Note the trailing slash for path parameters

	// Serve the Swagger UI (optional, but helpful for development)
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.json"}
	sh := middleware.SwaggerUI(opts, nil)
	mux.Handle("/docs", sh)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
