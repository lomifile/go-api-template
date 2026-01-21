package utils

import (
	"encoding/json"
	"fmt"
)

// IDResponse Id only response from endpoing or database
type IDResponse struct {
	ID int64 `db:"id" json:"id"`
}

// Scan Database bind for ID response
func (r *IDResponse) Scan(src any) error {
	if src == nil {
		*r = IDResponse{}
		return nil
	}
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("IdResponse: cannot scan %T", src)
	}
	return json.Unmarshal(b, r)
}

// DatabaseResponse base database response
type DatabaseResponse[T any] struct {
	Data          T      `db:"data"`
	StatusCode    int64  `db:"statusCode"`
	StatusMessage string `db:"statusMessage"`
	TS            string `db:"ts"`
}

// SuccessResponseMap Base succeess response
type SuccessResponseMap[T any] struct {
	RequestID string `json:"request_id"`
	Status    int    `json:"status"`
	Data      T      `json:"data"`
	TS        string `json:"ts"`
}

// ErrorResponseMap Base error response
type ErrorResponseMap struct {
	RequestID string `json:"request_id"`
	Status    int    `json:"status"`
	Error     string `json:"error"`
	TS        string `json:"ts"`
}

// PaginationResponseMeta base paginations response
type PaginationResponseMeta struct {
	Next        *int `json:"next"`
	Prev        *int `json:"previous"`
	HasNextPage bool `json:"has_next_page"`
	HasPrevPage bool `json:"has_previous_page"`
}

// PaginationResponse Base pagination response
type PaginationResponse[T any] struct {
	Items T                      `json:"items"`
	Total int                    `json:"total"`
	Meta  PaginationResponseMeta `json:"meta"`
}
