package repositories

import (
	"context"
	
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type attributeType types.ScalarAttributeType
type keyType types.KeyType
type projectionType types.ProjectionType

const (
	attributeTypeString = attributeType(types.ScalarAttributeTypeS)
	attributeTypeNumber = attributeType(types.ScalarAttributeTypeN)

	keyTypePartition = keyType(types.KeyTypeHash)
	keyTypeSorting = keyType(types.KeyTypeRange)
	keyTypeNone = keyType("")

	projectionTypeAll = projectionType(types.ProjectionTypeAll)
	projectionTypeKeysOnly = projectionType(types.ProjectionTypeKeysOnly)
	projectionTypeInclude = projectionType(types.ProjectionTypeInclude)
)

func NewRepositoryConfig(region, endpoint string) *RepositoryConfig {
	return &RepositoryConfig{
		region: region,
		endpoint: endpoint,
	}
}

type RepositoryConfig struct {
	region string
	endpoint string
}

type tableConfig struct {
	TableName string
	TableAttributes []tableAttribute
	GlobalSecondaryIndexes []globalSecondaryIndex
}


type tableAttribute struct {
	Name string
	AttrType attributeType
	KeyType keyType
}

type globalSecondaryIndex struct {
	IndexName string
	IndexAttributes []tableAttribute
	ProjectionType projectionType
}

type tableWaiterFunc func(ctx context.Context) error
