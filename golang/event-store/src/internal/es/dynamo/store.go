package dynamo

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Store keeps dependencies.
type Store struct {
	db    *dynamodb.Client
	table string
}

// NewStore creates Store instance.
func NewStore(dynamoClient *dynamodb.Client, table string) *Store {
	return &Store{
		db:    dynamoClient,
		table: table,
	}
}

func (s Store) ReadEvents(ctx context.Context, streamName string) ([]DBEventItem, error) {
	// TODO: Solve issue with max 1MB result (pagination)
	out, err := s.db.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.table),
		KeyConditionExpression: aws.String("#id = :hashKey"),
		ExpressionAttributeNames: map[string]string{
			"#id": "id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":hashKey": &types.AttributeValueMemberS{Value: streamName},
		},
		ScanIndexForward: aws.Bool(true), // important to read events in asc order
		ConsistentRead:   aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	if len(out.Items) == 0 {
		return nil, nil // TODO: should it be here or in invoking function?
	}

	var events []DBEventItem
	for _, rawItem := range out.Items {
		var item DBEventItem
		err = attributevalue.UnmarshalMap(rawItem, &item)
		if err != nil {
			return nil, err
		}
		events = append(events, item)
	}

	return events, nil
}

func (s Store) AppendEvents(ctx context.Context, items []DBEventItem) error {
	//TODO: should use transact_write_items
	for _, e := range items {
		item, err := attributevalue.MarshalMap(e)
		if err != nil {
			return err
		}

		out, err := s.db.PutItem(ctx, &dynamodb.PutItemInput{
			ConditionExpression:       nil,
			ExpressionAttributeNames:  nil,
			ExpressionAttributeValues: nil,
			Item:                      item,
			TableName:                 aws.String(s.table),
			ReturnValues:              "ALL_OLD",
		})

		if isConditionalCheckFailed(err) {
			// TODO: concurrent update
			// should retry
			return err
		}

		if err != nil {
			return err
		}

		if out.Attributes == nil { // successfully added

		}
	}

	return nil
}

// isConditionalCheckFailed checks if generic error is AWS specific one
// for types.ConditionalCheckFailedException.
func isConditionalCheckFailed(err error) bool {
	var conditionalCheckError *types.ConditionalCheckFailedException
	return errors.As(err, &conditionalCheckError)
}

type EventKey struct {
	ID      string `dynamodbav:"id"`
	Version int    `dynamodbav:"version"`
}

// DBEventItem is dynamoDB struct for event.
type DBEventItem struct {
	EventKey

	Type string `dynamodbav:"event_type"`
	Data string `dynamodbav:"event_data"`
}
