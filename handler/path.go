// Package handler contains http handlers func implementation.
package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PathFinder is the finder which is search for the path based on the input flights.
type PathFinder interface {
	Find(flights [][]string) ([]string, error)
}

// PostFlightRequestBody is the request body for flight/path Post.
type PostFlightRequestBody struct {
	Flights [][]string `json:"flights"`
}

// PostFlightResponseBody is the response body for flight/path Post.
type PostFlightResponseBody struct {
	Path []string `json:"path"`
}

// FlightPathHandler defines the flight path handler.
type FlightPathHandler struct {
	pathFinder PathFinder
}

// NewFlightPathHandler creates a new instance if FlightPathHandler.
func NewFlightPathHandler(pathFinder PathFinder) *FlightPathHandler {
	return &FlightPathHandler{
		pathFinder: pathFinder,
	}
}

// PostFlightPath returns handler func for the flight/path Post.
func (h *FlightPathHandler) PostFlightPath() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var flightsBody PostFlightRequestBody

		err := json.NewDecoder(r.Body).Decode(&flightsBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("can't decode body, %s", err.Error()), http.StatusBadRequest)
			return
		}

		path, err := h.pathFinder.Find(flightsBody.Flights)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid flights, %s", err.Error()), http.StatusBadRequest)
			return
		}

		payload, err := json.Marshal(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("can't marshal response, %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(payload); err != nil {
			http.Error(w, fmt.Sprintf("can't write response, %s", err.Error()), http.StatusInternalServerError)
			return
		}
	}
}
