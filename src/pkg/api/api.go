package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

// Error создает ErrorResponse
func Error(msg string) ErrorResponse {
	return ErrorResponse{msg}
}

// WriteJSON записывает Response в формате JSON в http.ResponseWriter
func WriteJSON(w http.ResponseWriter, statusCode int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

// ReadJSON считывает body из запроса в структуру
func ReadJSON[T any](r *http.Request) (*T, error) {
	var body T
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	return &body, err
}

func ParseUintQueryParam(r *http.Request, key string) (*uint, error) {
	// Get the query parameter value as a string
	queryParam := r.URL.Query().Get(key)

	// If the query parameter is not provided, return nil
	if queryParam == "" {
		return nil, nil
	}

	// Parse the query parameter to uint64
	parsedValue, err := strconv.ParseUint(queryParam, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid uint value for %s: %v", key, err)
	}

	// Convert the uint64 to uint
	value := uint(parsedValue)

	// Return a pointer to the uint
	return &value, nil
}

func ParseIntQueryParam(r *http.Request, key string) (*int, error) {
	// Get the query parameter value as a string
	queryParam := r.URL.Query().Get(key)

	// If the query parameter is not provided, return nil
	if queryParam == "" {
		return nil, nil
	}

	// Parse the query parameter to uint64
	parsedValue, err := strconv.ParseUint(queryParam, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid uint value for %s: %v", key, err)
	}

	// Convert the uint64 to uint
	value := int(parsedValue)

	// Return a pointer to the uint
	return &value, nil
}

func ParseStringQueryParam(r *http.Request, key string) *string {
	// Get the query parameter value as a string
	queryParam := r.URL.Query().Get(key)

	// If the query parameter is not provided, return nil
	if queryParam == "" {
		return nil
	}

	// Return a pointer to the queryParam string
	return &queryParam
}
