package repositories

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type queryParserFunction func(response *dynamodb.QueryOutput) ([]any, error)
type scanParserFunction func(response *dynamodb.ScanOutput) ([]any, error)
type getParserFunction func(response *dynamodb.GetItemOutput) (Item, error)

const (
	idTableName = "idCounter"
	idItemName = "current"
)

var (
	dynamodbClient *dynamodb.Client
)

func newBaseRepository(
	endpoint, region string,
	queryUnmarshallers map[string]queryParserFunction,
	getUnmarshallers map[string]getParserFunction,
	scanUnmarshallers map[string]scanParserFunction,
) *baseRepository {
	queryUnmarshallers[idTableName] = func(response *dynamodb.QueryOutput) ([]any, error) {
		return queryUnmarshaller[IdCounter](response)
	}
	getUnmarshallers[idTableName] = func(response *dynamodb.GetItemOutput) (Item, error) {
		return getUnmarshaller[IdCounter](response)
	}
	scanUnmarshallers[idTableName] = func(response *dynamodb.ScanOutput) ([]any, error) {
		return scanUnmarshaller[IdCounter](response)
	}

	return &baseRepository{
		endpoint: endpoint,
		region: region,

		queryUnmarshallers: queryUnmarshallers,
		getUnmarshallers: getUnmarshallers,
		scanUnmarshallers: scanUnmarshallers,
	}
}

type baseRepository struct {
	endpoint string
	region string

	client *dynamodb.Client
	waiter *dynamodb.TableExistsWaiter

	tableWaiters []tableWaiterFunc
	tables []*types.TableDescription

	queryUnmarshallers map[string]queryParserFunction
	scanUnmarshallers map[string]scanParserFunction
	getUnmarshallers map[string]getParserFunction

}

func (repo *baseRepository) Open() {
	if dynamodbClient != nil {
		repo.client = dynamodbClient
		repo.waiter = dynamodb.NewTableExistsWaiter(repo.client)
		return
	}

	awsConfig, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(repo.region),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	dynamodbClient = dynamodb.NewFromConfig(awsConfig, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(repo.endpoint)
	})

	repo.Open()
}

func (repo *baseRepository) createAllTables(ctx context.Context) error {
	if err := repo.createIdTable(ctx); err != nil {
		return fmt.Errorf("failed creating id counter table: %w", err)
	}

	if err := repo.waitTables(ctx); err != nil {
		return fmt.Errorf("failed waiting for tables to be created: %w", err)
	}

	if err := repo.createCurrentId(ctx); err != nil {
		return fmt.Errorf("failed creating current id: %w", err)
	}

	return nil
}

func (repo *baseRepository) addItem(ctx context.Context, tableName string, item Item) error {
	marshalled, err := attributevalue.MarshalMap(item)
	if err != nil {
		panic(err)
	}
	if _, err := repo.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: marshalled,
		},
	); err != nil {
		return fmt.Errorf("couldn't add item to table. Here's why: %w", err)
	}
	
	return nil
}

func (repo *baseRepository) getItem(ctx context.Context, tableName string, item Item) (Item, error) {
	response, err := repo.client.GetItem(
		ctx,
		&dynamodb.GetItemInput{
			Key: item.GetKey(),
			TableName: aws.String(tableName),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't get info about %+v. Here's why: %w", item, err)
	}

	item, err = repo.getUnmarshallers[tableName](response)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshall %+v. Here's why: %w", item, err)
	}

	return item, nil
}

func (repo *baseRepository) queryItems(
	ctx context.Context, tableName, indexName, key string, value any, limit int, cursor map[string]types.AttributeValue,
) (items []any, nextCursor map[string]types.AttributeValue, err error) {
	var index *string
	if indexName != "" {
		index = aws.String(indexName)
	}

	keyEx := expression.Key(key).Equal(expression.Value(value))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		err = fmt.Errorf("couldn't build expression for query. Here's why: %w", err)
		return
	}

	queryPaginator := dynamodb.NewQueryPaginator(repo.client, &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		IndexName:                 index,
		Limit:                     aws.Int32(int32(limit)),
		ExclusiveStartKey:         cursor,
	})

	if queryPaginator.HasMorePages() {
		var response *dynamodb.QueryOutput
		response, err = queryPaginator.NextPage(ctx)
		if err != nil {
	err = fmt.Errorf("couldn't query for %s with key: %q and value: %+v. Here's why: %w", tableName, key, value, err)
			return
		}
		
		items, err = repo.queryUnmarshallers[tableName](response)
		if err != nil {
			err = fmt.Errorf("couldn't parse items: %w", err)
			return
		}

		nextCursor = response.LastEvaluatedKey
	}
	
	return
}

func (repo *baseRepository) scanItems(
	ctx context.Context, tableName string, limit int, cursor map[string]types.AttributeValue,
) (items []any, nextCursor map[string]types.AttributeValue, err error) {
	scanPaginator := dynamodb.NewScanPaginator(
		repo.client,
		&dynamodb.ScanInput{
			TableName: aws.String(tableName),
			Limit: aws.Int32(int32(limit)),
			ExclusiveStartKey: cursor,
		},
	)
	
	if scanPaginator.HasMorePages() {
		var response *dynamodb.ScanOutput
		response, err = scanPaginator.NextPage(ctx)
		if err != nil {
			err = fmt.Errorf("couldn't scan for %s. Here's why: %w", tableName, err)
			return
		}
		
		items, err = repo.scanUnmarshallers[tableName](response)
		if err != nil {
			err = fmt.Errorf("couldn't parse items: %w", err)
			return
		}

		nextCursor = response.LastEvaluatedKey
	}

	return 
}

func (repo *MovieRepository) deleteItem(
	ctx context.Context, tableName string, item Item,
) error {
	_, err := repo.client.DeleteItem(
		ctx,
		&dynamodb.DeleteItemInput{
	TableName: aws.String(tableName),
	Key: item.GetKey(),
},
	)
	if err != nil {
		return fmt.Errorf("couldn't delete %+v from table. Here's why: %w", item, err)
	}
	return nil
}

func (repo *baseRepository) createTable(ctx context.Context, tableCfg *tableConfig) (*dynamodb.CreateTableOutput, error) {
	attrDefs, keySchema := repo.genTablePrimaryIndex(tableCfg.TableAttributes)
	globalSecondaryIndexes := repo.genTableSecondaryIndexes(tableCfg.GlobalSecondaryIndexes)
	if len(globalSecondaryIndexes) == 0 {
		globalSecondaryIndexes = nil
	}

	createTableInput := &dynamodb.CreateTableInput{
		AttributeDefinitions: attrDefs,
		KeySchema: keySchema,
		GlobalSecondaryIndexes: globalSecondaryIndexes,
		TableName:   aws.String(tableCfg.TableName),
		BillingMode: types.BillingModePayPerRequest,
	}

	table, err := repo.client.CreateTable(ctx, createTableInput)
	if err != nil {
		return nil, fmt.Errorf("couldn't create table %v with config %+v. Here's why: %w",
			tableCfg.TableName,
			createTableInput,
			err,
		)
	}

	return table, nil
}

func (repo *baseRepository) createIdTable(ctx context.Context) (error) {
	table, err := repo.createTable(ctx, &tableConfig{
		TableName: idTableName,
		TableAttributes: []tableAttribute{
			{
				Name: "name",
				AttrType: attributeTypeString,
				KeyType: keyTypePartition,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error creating id table: %w", err)
	}

	repo.addWaiter(idTableName, table)
	return nil
}

func (repo *baseRepository) createCurrentId(ctx context.Context) error {
	id := IdCounter{Name: idItemName}
	
	if idItem, err := repo.getItem(ctx, idTableName, id); err != nil {
		return fmt.Errorf("could not get current id: %w", err)
	} else if idItem == nil {
		id.Id = 1
		if err := repo.addItem(ctx, idTableName, id); err != nil {
			return fmt.Errorf("could not create current id: %w", err)
		}
	}
	return nil
}


func (repo *baseRepository) getNextId(ctx context.Context) (int, error) {
	update := expression.Set(expression.Name("id"), expression.Name("id").Plus(expression.Value(1)))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return 0, fmt.Errorf("could't build expression to upgrade id: %w", err)
	}

	id := IdCounter{Name: idItemName}
	response, err := repo.client.UpdateItem(
		ctx,
		&dynamodb.UpdateItemInput{
			TableName:                   aws.String(idTableName),
			Key:                         id.GetKey(),
			ExpressionAttributeNames:    expr.Names(),
			ExpressionAttributeValues:   expr.Values(),
			UpdateExpression:            expr.Update(),
			ReturnValues:                types.ReturnValueUpdatedNew,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("couldn't update id: %w", err)
	}

	
	if err = attributevalue.UnmarshalMap(response.Attributes, &id); err != nil {
		return 0, fmt.Errorf("couldn't unmarshall update response. Here's why: %w", err)
	}

	return id.Id, nil
}


func (repo *baseRepository) waitTables(ctx context.Context) error {
	for _, waiter := range repo.tableWaiters {
		if err := waiter(ctx); err != nil {
			return fmt.Errorf("error waiting for table creation: %w", err)
		}
	}
	repo.tableWaiters = []tableWaiterFunc{}
	return nil
}



func (repo *baseRepository) addWaiter(tableName string, table *dynamodb.CreateTableOutput) {
	tableWaiter := func (ctx context.Context) error {
		if err := repo.waiter.Wait(
			ctx,
			&dynamodb.DescribeTableInput{
				TableName: aws.String(tableName),
			},
			5*time.Minute,
		); err != nil {
			return fmt.Errorf("wait for table exists failed. Here's why: %w", err)
		}
		tableDesc := table.TableDescription
		
		repo.tables = append(repo.tables, tableDesc)
		return nil
	}

	repo.tableWaiters = append(repo.tableWaiters, tableWaiter)
}

func (repo *baseRepository) genTablePrimaryIndex(
	attributes []tableAttribute,
) ([]types.AttributeDefinition, []types.KeySchemaElement) {
	attributeDefinitions := make([]types.AttributeDefinition, len(attributes))
	keySchema := []types.KeySchemaElement{}

	for index, attribute := range(attributes) {
		attributeDefinitions[index] = types.AttributeDefinition{
			AttributeName: aws.String(attribute.Name),
			AttributeType: types.ScalarAttributeType(attribute.AttrType),
		}
		if attribute.KeyType != keyTypeNone {
			keySchema = append(keySchema, types.KeySchemaElement{
				AttributeName: aws.String(attribute.Name),
				KeyType: types.KeyType(attribute.KeyType),
			})
		}
	}

	return attributeDefinitions, keySchema
}

func (repo *baseRepository) genTableSecondaryIndexes(GSIs []globalSecondaryIndex) []types.GlobalSecondaryIndex {
	globalSecondaryIndexes := make([]types.GlobalSecondaryIndex, len(GSIs))
	for index, gsi := range(GSIs) {
		globalSecondaryIndexes[index] = repo.genTableSecondaryIndex(gsi)
	}
	return globalSecondaryIndexes
}

func (repo *baseRepository) genTableSecondaryIndex(gsi globalSecondaryIndex) types.GlobalSecondaryIndex {
	keySchema := make([]types.KeySchemaElement, len(gsi.IndexAttributes))
	for index, attribute := range(gsi.IndexAttributes) {
		keySchema[index] = types.KeySchemaElement{
			AttributeName: aws.String(attribute.Name),
			KeyType: types.KeyType(attribute.KeyType),
		}
	}
	
	return types.GlobalSecondaryIndex{
		IndexName: aws.String(gsi.IndexName),
		KeySchema: keySchema,
		Projection: &types.Projection{
			ProjectionType: types.ProjectionType(gsi.ProjectionType),
		},
	}
}



func queryUnmarshaller[itemType any](response *dynamodb.QueryOutput) ([]any, error) {
	var items []itemType
	if err := attributevalue.UnmarshalListOfMaps(response.Items, &items); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal query response. Here's why: %v", err)
	}
	anyItems := make([]any, len(items))
	for index, item := range items {
		anyItems[index] = item
	}
	return anyItems, nil
}

func scanUnmarshaller[itemType any](response *dynamodb.ScanOutput) ([]any, error) {
	var items []itemType
	if err := attributevalue.UnmarshalListOfMaps(response.Items, &items); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal query response. Here's why: %v", err)
	}
	anyItems := make([]any, len(items))
	for index, item := range items {
		anyItems[index] = item
	}
	return anyItems, nil
}

func getUnmarshaller[itemType Item](response *dynamodb.GetItemOutput) (Item, error) {
	if response.Item == nil {
		return nil, nil
	}

	var movie itemType
	if err := attributevalue.UnmarshalMap(response.Item, &movie); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response. Here's why: %w", err)
	}
	return movie, nil	
}

