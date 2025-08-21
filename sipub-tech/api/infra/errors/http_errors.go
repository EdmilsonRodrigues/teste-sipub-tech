package errors

import (
	"fmt"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/dtos"
)

var (
	BadRequestResponse = dtos.NewErrorResponse(dtos.ErrorMessage(fmt.Errorf("bad request")))
	UnprocessableEntityResponse = dtos.NewErrorResponse(
		dtos.ErrorMessage(fmt.Errorf("unprocessable Entity: ")),
	)
	InternalServerErrorResponse = dtos.NewErrorResponse(dtos.ErrorMessage(fmt.Errorf("internal Server Error")))
	MovieNotFoundErrorResponse = dtos.NewErrorResponse(dtos.ErrorMessage(fmt.Errorf("movie not found")))
)

func InternalServerError(message string) *dtos.ErrorResponse {
	err := InternalServerErrorResponse
	err.Details.Message += ": " + message
	return err
}

func UnprocessableEntity(message string) *dtos.ErrorResponse {
	err := UnprocessableEntityResponse
	err.Details.Message += ": " + message
	return err
}


