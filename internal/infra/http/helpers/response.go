package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

type PaginatedResponse[T any] struct {
	Data       []*T           `json:"data"`
	Extra      map[string]any `json:"extra,omitempty"`
	Pagination PaginationMeta `json:"pagination"`
}

type PaginationMeta struct {
	Total       int64 `json:"total"`
	Offset      int   `json:"offset"`
	Limit       int   `json:"limit"`
	TotalPages  int   `json:"totalPages"`
	HasNext     bool  `json:"hasNext"`
	HasPrevious bool  `json:"hasPrevious"`
}

func NewPaginatedResponse[T any](data []*T, total int64, limit, offset int) *PaginatedResponse[T] {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &PaginatedResponse[T]{
		Data:  data,
		Extra: make(map[string]any),
		Pagination: PaginationMeta{
			Total:       total,
			Offset:      offset,
			Limit:       limit,
			TotalPages:  totalPages,
			HasNext:     offset+1 < totalPages,
			HasPrevious: offset+1 > 1,
		},
	}
}
