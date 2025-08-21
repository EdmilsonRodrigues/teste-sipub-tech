package middlewares_test

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"	

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/core/dtos"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/api/infra/middlewares"
)

func TestParseQueryParameters(t *testing.T) {
	queryParser := middlewares.ParseQueryParameters()
	minimumYear := 1880
	maximumYear := time.Now().Year()

	t.Run("should parse correctly if given parameters are correct", func(t *testing.T) {
		limit, cursor := rand.Int32(), rand.Int32()
		year := rand.IntN(maximumYear - minimumYear) + minimumYear
		
		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%d", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if ctx.IsAborted() {
			t.Errorf("Context was aborted even given right input")
		}
		
		dto, exists := ctx.Get("dto")
		if !exists {
			t.Errorf("DTO not set in context.")
		}

		result, ok := dto.(*dtos.MoviesQueryDTO)
		if !ok {
			t.Errorf("DTO malformed.")
		}		

		expected := dtos.MoviesQueryDTO{
			Year: strconv.Itoa(year),
			Limit: int(limit),
			Cursor: int(cursor),
		}

		assert.Equal(t, expected, *result)
	})
	
	t.Run("should parse correctly if given parameters are correct but year is missing", func(t *testing.T) {
		limit, cursor := rand.Int32(), rand.Int32()

		path := fmt.Sprintf("/movies/?&limit=%d&cursor=%d", limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if ctx.IsAborted() {
			t.Errorf("Context was aborted even given right input")
		}
		
		dto, exists := ctx.Get("dto")
		if !exists {
			t.Errorf("DTO not set in context.")
		}

		result, ok := dto.(*dtos.MoviesQueryDTO)
		if !ok {
			t.Errorf("DTO malformed.")
		}		

		expected := dtos.MoviesQueryDTO{
			Limit: int(limit),
			Cursor: int(cursor),
		}

		assert.Equal(t, expected, *result)
	})

	t.Run("should parse correctly if given parameters are correct but limit is missing", func(t *testing.T) {
		cursor := rand.Int32()
		year := rand.IntN(maximumYear - minimumYear) + minimumYear
		
		path := fmt.Sprintf("/movies/?year=%d&cursor=%d", year, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if ctx.IsAborted() {
			t.Errorf("Context was aborted even given right input")
		}
		
		dto, exists := ctx.Get("dto")
		if !exists {
			t.Errorf("DTO not set in context.")
		}

		result, ok := dto.(*dtos.MoviesQueryDTO)
		if !ok {
			t.Errorf("DTO malformed.")
		}		

		expected := dtos.MoviesQueryDTO{
			Year: strconv.Itoa(year),
			Cursor: int(cursor),
		}

		assert.Equal(t, expected, *result)
	})

	t.Run("should parse correctly if given parameters are correct but cursor is missing", func(t *testing.T) {
		limit := rand.Int32()
		year := rand.IntN(maximumYear - minimumYear) + minimumYear
		
		path := fmt.Sprintf("/movies/?year=%d&limit=%d", year, limit)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if ctx.IsAborted() {
			t.Errorf("Context was aborted even given right input")
		}
		
		dto, exists := ctx.Get("dto")
		if !exists {
			t.Errorf("DTO not set in context.")
		}

		result, ok := dto.(*dtos.MoviesQueryDTO)
		if !ok {
			t.Errorf("DTO malformed.")
		}		

		expected := dtos.MoviesQueryDTO{
			Year: strconv.Itoa(year),
			Limit: int(limit),
		}

		assert.Equal(t, expected, *result)
	})

	t.Run("should fail to parse if given year before minimum", func(t *testing.T) {
		limit, cursor := rand.Int32(), rand.Int32()
		year := rand.IntN(minimumYear)
		
		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%d", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if !ctx.IsAborted() {
			t.Errorf("Context was not aborted when year was incorrect")
		}
	})

	t.Run("should fail to parse if given year after maximum", func(t *testing.T) {
		limit, cursor := rand.Int32(), rand.Int32()
		year := rand.IntN(minimumYear) + maximumYear

		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%d", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if !ctx.IsAborted() {
			t.Errorf("Context was not aborted when year was incorrect")
		}
	})

	t.Run("should fail to parse if given limit was not an integer", func(t *testing.T) {
		limit, cursor := "abacate", rand.Int32()
		year := rand.IntN(maximumYear - minimumYear) + minimumYear

		path := fmt.Sprintf("/movies/?year=%d&limit=%s&cursor=%d", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if !ctx.IsAborted() {
			t.Errorf("Context was not aborted when year was incorrect")
		}
	})

	t.Run("should fail to parse if given limit was a negative integer", func(t *testing.T) {
		limit, cursor := -rand.Int32(), rand.Int32()
		year := rand.IntN(maximumYear - minimumYear) + minimumYear

		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%d", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if !ctx.IsAborted() {
			t.Errorf("Context was not aborted when year was incorrect")
		}
	})

	t.Run("should fail to parse if given cursor was not an integer", func(t *testing.T) {
		limit, cursor := rand.Int32(), "abacate"
		year := rand.IntN(maximumYear - minimumYear) + minimumYear

		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%s", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if !ctx.IsAborted() {
			t.Errorf("Context was not aborted when year was incorrect")
		}
	})

	t.Run("should fail to parse if given cursor was a negative integer", func(t *testing.T) {
		limit, cursor := rand.Int32(), -rand.Int32()
		year := rand.IntN(maximumYear - minimumYear) + minimumYear

		path := fmt.Sprintf("/movies/?year=%d&limit=%d&cursor=%d", year, limit, cursor)
		req, _ := http.NewRequest("GET", path, nil)
		ctx := getContext(req)

		queryParser(ctx)
		if !ctx.IsAborted() {
			t.Errorf("Context was not aborted when year was incorrect")
		}
	})
}

func getContext(req *http.Request) *gin.Context {
	return &gin.Context{
		Request: req,
		Writer: &FakeWriter{HeadersMapping: make(http.Header)},
	}
}

type FakeWriter struct {
	gin.ResponseWriter

	StatusCode     int
	HeadersMapping http.Header
	Body           []byte
}

func (w *FakeWriter) WriteHeader(code int) {
	w.StatusCode = code
}

func (w *FakeWriter) Header() http.Header {
	return w.HeadersMapping
}

func (w *FakeWriter) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return len(b), nil
}

func (w *FakeWriter) Status() int {
	return w.StatusCode
}

