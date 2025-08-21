package dtos

import (
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
)

type ErrorMessage error

type JsonData interface {
	ToDataItem() dtos.DataItem
}


// Error Response
func NewErrorResponse(msg ErrorMessage) *ErrorResponse {
	return &ErrorResponse{
		Details: ErrorDetails{
			Message: msg.Error(),
		},
	}
}

type ErrorResponse struct {
	Details ErrorDetails  `json:"details"`
}

type ErrorDetails struct {
	Message string  `json:"message"`
}

// Normal Response
func NewJSONResponse(data JsonData) *JSONResponse {
	return &JSONResponse{Data: data.ToDataItem()}
}

type JSONResponse struct {
	Data dtos.DataItem  `json:"data"`
}


// Paginated Response
func NewPaginatedResponse(data []JsonData, limit, cursor int) *PaginatedJSONResponse {
	dataItems := structArrayToDataItemArray(data)
	return &PaginatedJSONResponse{
		Data:  dataItems,
		Limit: limit,
		Cursor: cursor,
	}
}

type PaginatedJSONResponse struct {
	Data   []dtos.DataItem  `json:"data"`
	Limit  int         `json:"limit"`
	Cursor int         `json:"cursor"`
}

func structArrayToDataItemArray(data []JsonData) []dtos.DataItem {
	response := make([]dtos.DataItem, len(data)) 
	for index, item := range data {
		response[index] = item.ToDataItem()
	}
	return response
}
