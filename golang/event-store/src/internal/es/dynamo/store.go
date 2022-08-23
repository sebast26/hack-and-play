package dynamo

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
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

var (
	ErrConcurrentUpdate = errors.New("concurrent update on item")
)

func (s Store) ReadEvents(ctx context.Context, streamName string) ([]DBEventItem, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(expression.KeyEqual(expression.Key(idField), expression.Value(streamName))).
		Build()
	if err != nil {
		return nil, err
	}

	var events []DBEventItem
	var lastEvaluatedKey map[string]types.AttributeValue
	for {
		out, err := s.db.Query(ctx, &dynamodb.QueryInput{
			TableName:                 aws.String(s.table),
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			ScanIndexForward:          aws.Bool(true), // important to read events in asc order
			ConsistentRead:            aws.Bool(true),
			ExclusiveStartKey:         lastEvaluatedKey,
		})
		if err != nil {
			return nil, err
		}

		lastEvaluatedKey = out.LastEvaluatedKey

		for _, rawItem := range out.Items {
			var item DBEventItem
			err = attributevalue.UnmarshalMap(rawItem, &item)
			if err != nil {
				return nil, err
			}
			events = append(events, item)
		}

		// break when crossed the limit or no more to scan
		if out.LastEvaluatedKey == nil {
			break
		}
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

		expr, err := expression.NewBuilder().
			WithCondition(expression.And(
				expression.AttributeNotExists(expression.Name(idField)),
				expression.AttributeNotExists(expression.Name(versionField)))).Build()
		if err != nil {
			return err
		}

		out, err := s.db.PutItem(ctx, &dynamodb.PutItemInput{
			ConditionExpression:      expr.Condition(),
			ExpressionAttributeNames: expr.Names(),
			Item:                     item,
			TableName:                aws.String(s.table),
			ReturnValues:             "ALL_OLD",
		})

		if isConditionalCheckFailed(err) {
			return ErrConcurrentUpdate
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

const (
	idField      = "id"
	versionField = "version"
)
