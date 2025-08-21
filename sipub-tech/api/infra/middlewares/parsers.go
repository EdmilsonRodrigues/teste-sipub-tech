package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
	
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	infraDtos "github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/dtos"
)

const (
	DtoKey         = "dto"
	QueryYearKey   = "year"
	QueryLimitKey  = "limit"
	QueryCursorKey = "cursor"
)

func ParseQueryParameters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		year, limit, cursor, err := parseQueryParams(
			ctx.Query(QueryYearKey),
			ctx.Query(QueryLimitKey),
			ctx.Query(QueryCursorKey),
		)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, infraDtos.NewErrorResponse(infraDtos.ErrorMessage(err)))
			ctx.Abort()
			return
		}
		
		dto := &dtos.MoviesQueryDTO{
			Year:   year,
			Limit:  limit,
			Cursor: cursor,
		}

		ctx.Set(DtoKey, dto)
		ctx.Next()
	}
}

func parseQueryParams(year, limit, cursor string) (parsedYear string, parsedLimit, parsedCursor int, err error) {	
	parsedLimit, err = parseLimit(limit)
	if err != nil {
		err = parseErrorAndLog("limit", limit, "a positive integer.", err)
		return
	}

	const minimumYear = 1880
	currentYear := time.Now().Year()
	parsedYear, err = parseYear(year, minimumYear, currentYear)
	if err != nil {
		err = parseErrorAndLog("year", year, fmt.Sprintf("an year between %d and %d", minimumYear, currentYear), err)
		return
	}
	
	parsedCursor, err = parseCursor(cursor)
	if err != nil {
		err = parseErrorAndLog("cursor", cursor, "a positive integer.", err)
		return
	}
	return
}

func parseLimit(limit string) (int, error) {
	if limit == "" {
		return 0, nil
	}
	
	parsed, err := strconv.Atoi(limit)
	if err != nil {
		return 0, fmt.Errorf("failed parsing limit: %w", err)
	}
	if parsed < 1 {
		return 0, fmt.Errorf("limit cannot be zero nor negative")
	}
	return parsed, nil
}

func parseYear(year string, min, max int) (string, error) {
	if year == "" {
		return "", nil
	}
	parsed, err := strconv.Atoi(year)
	if err != nil {
		return "", fmt.Errorf("failed parsing limit: %w", err)
	}
	if parsed < min {
		return "", fmt.Errorf("year cannot be before %d", min)
	}
	if parsed > max {
		return "", fmt.Errorf("year cannot be after %d", max)
	}

	return strconv.Itoa(parsed), nil
}

func parseCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(cursor)
	if err != nil {
		return 0, fmt.Errorf("failed parsing limit: %w", err)
	}
	if parsed < 1 {
		return 0, fmt.Errorf("cursor cannot be zero nor negative")
	}
	return parsed, nil
}


func parseErrorAndLog(field, original, mustBe string, err error) error {
	log.Printf("Failed parsing %s %q: %v", field, original, err)
	return fmt.Errorf("error parsing request. Query %q malformed. It must be %s", field, mustBe) 
}
