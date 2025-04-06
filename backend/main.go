package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type TrafficResponse struct {
	Area         string `json:"area"`
	VehicleCount int    `json:"vehicle_count"`
	TrafficJam   bool   `json:"traffic_jam"`
	RoadStatus   string `json:"road_status"`
}

func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the area parameter from the URL query
	area := r.URL.Query().Get("area")
	if area == "" {
		http.Error(w, "Area parameter is required", http.StatusBadRequest)
		return
	}

	// Call FastAPI service (adjust the FastAPI URL accordingly)
	fastapiURL := fmt.Sprintf("http://localhost:8000/analyze?area=%s", url.QueryEscape(area))
	resp, err := http.Get(fastapiURL)
	if err != nil {
		http.Error(w, "Failed to contact analysis service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Parse the FastAPI response into the TrafficResponse struct
	var result TrafficResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")
	// Return the result as a JSON response
	json.NewEncoder(w).Encode(result)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the form HTML page
	http.ServeFile(w, r, "form.html")
}

func main() {
	// Define routes
	http.HandleFunc("/", homeHandler)       // Home page with the form
	http.HandleFunc("/analyze", analyzeHandler) // Analyze traffic and return data

	// Start the server
	fmt.Println("Golang server running on http://localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
