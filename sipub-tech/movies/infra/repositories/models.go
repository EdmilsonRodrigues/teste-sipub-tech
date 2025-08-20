package repositories

import (
	"fmt"
	
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
)

type Item interface {
	GetKey() map[string]types.AttributeValue
}

func NewDBMovie(movie *domain.Movie, id int) *DBMovie {
	return &DBMovie{
		Id: id,
		Year: movie.Year,
		Title: movie.Title,
	}
}

type DBMovie struct {
	Title string `dynamodbav:"title"`
	Year string  `dynamodbav:"year"`
	Id int       `dynamodbav:"id"`
}

func (movie DBMovie) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(movie.Id)
	if err != nil {
		panic(fmt.Errorf("failed marshalling id: %w", err))
	}
	return map[string]types.AttributeValue{"id": id}
}

type IdCounter struct {
	Name string `dynamodbav:"name"`
	Id int      `dynamodbav:"id"`
}

func (counter IdCounter) GetKey() map[string]types.AttributeValue {
	current, err := attributevalue.Marshal(counter.Name)
	if err != nil {
		panic(fmt.Errorf("failed marshalling counter name: %w", err))
	}
	return map[string]types.AttributeValue{"name": current}
}


