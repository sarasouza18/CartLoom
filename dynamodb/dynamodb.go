package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// NewDynamoDBClient initializes a DynamoDB client for the specified region.
func NewDynamoDBClient(ctx context.Context, region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

// CreateTable creates a DynamoDB table with the specified schema.
func CreateTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	input := buildCreateTableInput(tableName)

	_, err := client.CreateTable(ctx, input)
	if err != nil {
		log.Printf("Failed to create table %s: %v", tableName, err)
		return err
	}

	log.Printf("Table %s created successfully", tableName)
	return nil
}

// EnableGlobalReplication adds global replication to an existing table.
func EnableGlobalReplication(ctx context.Context, client *dynamodb.Client, tableName, region string) error {
	input := buildGlobalReplicationInput(tableName, region)

	_, err := client.UpdateTable(ctx, input)
	if err != nil {
		log.Printf("Failed to add global replication for table %s to region %s: %v", tableName, region, err)
		return err
	}

	log.Printf("Global replication added for table %s in region %s", tableName, region)
	return nil
}

// buildCreateTableInput constructs the CreateTableInput for a DynamoDB table.
func buildCreateTableInput(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("OrderID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("OrderID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}
}

// buildGlobalReplicationInput constructs the UpdateTableInput for global table replication.
func buildGlobalReplicationInput(tableName, region string) *dynamodb.UpdateTableInput {
	return &dynamodb.UpdateTableInput{
		TableName: aws.String(tableName),
		ReplicaUpdates: []types.ReplicationGroupUpdate{
			{
				Create: &types.CreateReplicationGroupMemberAction{
					RegionName: aws.String(region),
				},
			},
		},
	}
}
