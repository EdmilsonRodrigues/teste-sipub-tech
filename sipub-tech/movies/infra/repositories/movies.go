package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/domain"
	"github.com/EdmilsonRodrigues/teste-sipub-tech/sipub-tech/movies/core/ports"
)

const (
	movieTableName = "movies"
	searchMoviesByYearIndex = "search-by-year-index"
)

func NewMovieRepository(config *RepositoryConfig) *MovieRepository {
	return &MovieRepository{
		baseRepository: *newBaseRepository(
			config.endpoint,
			config.region,
			map[string]queryParserFunction{
				movieTableName: func(response *dynamodb.QueryOutput) ([]any, error) {
					return queryUnmarshaller[DBMovie](response)
				},
			},
			map[string]getParserFunction{
				movieTableName: func(response *dynamodb.GetItemOutput) (Item, error) {
					return getUnmarshaller[DBMovie](response)
				},
			},
			map[string]scanParserFunction{
				movieTableName: func(response *dynamodb.ScanOutput) ([]any, error) {
					return scanUnmarshaller[DBMovie](response)
				},
			},
		),
	}
}

type MovieRepository struct {
	baseRepository
}


func (repo *MovieRepository) CreateTables(ctx context.Context) error {	
	if err := repo.createMovieTable(ctx); err != nil {
		return fmt.Errorf("failed creating movie table: %w", err)
	}

	return repo.createAllTables(ctx)
}

func (repo *MovieRepository) GetOne(ctx context.Context, id int) (movie domain.Movie, err error) {
	query := DBMovie{Id: id}

	rawMovie, err := repo.getItem(ctx, movieTableName, query)
	if err != nil {
		err = fmt.Errorf("error getting movie with id %d: %w", id, err)
		return 
	} else if rawMovie == nil {
		err = ports.ErrMovieNotFound
		return 
	}

	movie, err = repo.parseMovie(rawMovie)
	return
}

func (repo *MovieRepository) GetAll(
	ctx context.Context, year string, limit int, lastMovieId int,
) (movies []domain.Movie, cursor int, err error) {
	var query map[string]types.AttributeValue
	if lastMovieId != 0 {
		query = DBMovie{Id: lastMovieId}.GetKey()
	}

	if year == "" {
		var fetchedMovies []any
		var cursorMap map[string]types.AttributeValue

		fetchedMovies, cursorMap, err = repo.scanItems(ctx, movieTableName, limit, query)
		if err != nil {
			err = fmt.Errorf("failed scanning for movies: %w", err)
			return 
		}

		movies, cursor, err = repo.parseFetches(cursorMap, fetchedMovies)

	} else {
		var fetchedMovies []any
		var cursorMap map[string]types.AttributeValue

		fetchedMovies, cursorMap, err = repo.queryItems(ctx, movieTableName, searchMoviesByYearIndex, "year", year, limit, query)
		if err != nil {
			err = fmt.Errorf("failed querying for movies: %w", err)
			return 
		}

		movies, cursor, err = repo.parseFetches(cursorMap, fetchedMovies)
	}

	return
}


func (repo *MovieRepository) Save(ctx context.Context, movie domain.Movie) error {
	id, err := repo.getNextId(ctx)
	if err != nil {
		return fmt.Errorf("error getting the id of new movie %w", err)
	}
	
	parsedMovie := DBMovie{
		Id: id,
		Title: movie.Title,
		Year: movie.Year,
	}
	if err := repo.addItem(ctx, movieTableName, parsedMovie); err != nil {
		return fmt.Errorf("failed saving movie %+v: %w", movie, err)
	}

	return nil
}

func (repo *MovieRepository) Delete(ctx context.Context, id int) error {
	query := DBMovie{Id: id}
	if err := repo.deleteItem(ctx, movieTableName, query); err != nil {
		return fmt.Errorf("failed deleting movie with id %d: %w", id, err)
	}
	return nil
}

func (repo *MovieRepository) createMovieTable(ctx context.Context) (error) {
	table, err := repo.createTable(ctx, &tableConfig{
		TableName: movieTableName,
		TableAttributes: []tableAttribute{
			{
				Name: "id",
				AttrType: attributeTypeNumber,
				KeyType: keyTypePartition,
			},
			{
				Name: "year",
				AttrType: attributeTypeString,
				KeyType: keyTypeNone,
			},
			{
				Name: "title",
				AttrType: attributeTypeString,
				KeyType: keyTypeNone,
			},
		},
		GlobalSecondaryIndexes: []globalSecondaryIndex{
			{
				IndexName: searchMoviesByYearIndex,
				IndexAttributes: []tableAttribute{
					{
						Name: "year",
						AttrType: attributeTypeString,
						KeyType: keyTypePartition,
					},
					{
						Name: "title",
						AttrType: attributeTypeString,
						KeyType: keyTypeSorting,
					},
				},
				ProjectionType: projectionTypeKeysOnly,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating movies table: %w", err)
	}

	repo.addWaiter(movieTableName, table)
	return nil
}

func (repo *MovieRepository) parseMovies(fetchedMovies []any) ([]domain.Movie, error) {
	movies := make([]domain.Movie, len(fetchedMovies))
	for index, movie := range fetchedMovies {
		parsedMovie, ok := movie.(DBMovie)
		if !ok {
			return nil, fmt.Errorf("failed parsing movie: %+v", movie)
		}
		movies[index] = domain.Movie{
			ID: parsedMovie.Id,
			Title: parsedMovie.Title,
			Year: parsedMovie.Year,
		}
	}

	return movies, nil
}

func (repo *MovieRepository) parseMovie(rawMovie Item) (movie domain.Movie, err error) {
	dbMovie, ok := rawMovie.(DBMovie) 
	if !ok {
		err = fmt.Errorf("error parsing received movie:i %+v", rawMovie)
		return 
	}

	movie = domain.Movie{
		ID: dbMovie.Id,
		Title: dbMovie.Title,
		Year: dbMovie.Year,
	}

	return
}

func (repo *MovieRepository) parseCursor(rawCursor map[string]types.AttributeValue) (int, error) {
	if rawCursor == nil {
		return 0, nil
	}
	
	var idMap map[string]int

	if err := attributevalue.UnmarshalMap(rawCursor, &idMap); err != nil {
		return  0, fmt.Errorf("couldn't unmarshal response. Here's why: %w", err)
	}
	
	return idMap["id"], nil	
}

func (repo *MovieRepository) parseFetches(
	rawCursor map[string]types.AttributeValue, rawMovies []any,
) (movies []domain.Movie, cursor int, err error) {
	cursor, err = repo.parseCursor(rawCursor)
	if err != nil {
		err = fmt.Errorf("failed parsing cursor %+v: %w", rawCursor, err)
		return
	}
	
	movies, err = repo.parseMovies(rawMovies)
	return
}
