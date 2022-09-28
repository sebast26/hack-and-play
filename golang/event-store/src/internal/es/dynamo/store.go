package dynamo

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"sgorecki.me/golang/event-store/src/internal/clock"
)

// Store keeps dependencies.
type Store struct {
	db    *dynamodb.Client
	table string
	clock clock.Clock
}

// NewStore creates Store instance.
func NewStore(dynamoClient *dynamodb.Client, table string, clock clock.Clock) *Store {
	return &Store{
		db:    dynamoClient,
		table: table,
		clock: clock,
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
	if len(items) == 0 {
		return nil
	} else if len(items) == 1 {
		return s.appendEvent(ctx, items[0])
	} else {
		return s.appendEventsTransaction(ctx, items)
	}
}

func (s Store) appendEvent(ctx context.Context, event DBEventItem) error {
	item, err := s.prepareEvent(event)
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

	return nil
}

func (s Store) appendEventsTransaction(ctx context.Context, events []DBEventItem) error {
	expr, err := expression.NewBuilder().
		WithCondition(expression.And(
			expression.AttributeNotExists(expression.Name(idField)),
			expression.AttributeNotExists(expression.Name(versionField)))).Build()
	if err != nil {
		return err
	}

	transactItems := make([]types.TransactWriteItem, len(events))
	for i, event := range events {
		item, err := s.prepareEvent(event)
		if err != nil {
			return err
		}
		transactItems[i] = types.TransactWriteItem{
			Put: &types.Put{
				Item:                     item,
				TableName:                aws.String(s.table),
				ConditionExpression:      expr.Condition(),
				ExpressionAttributeNames: expr.Names(),
			},
		}
	}

	_, err = s.db.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactItems,
	})

	return err
}

func (s Store) prepareEvent(event DBEventItem) (map[string]types.AttributeValue, error) {
	event.CreatedAt = s.clock.Now().Format(dateTimeFormat)
	return attributevalue.MarshalMap(event)
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

	Type      string `dynamodbav:"event_type"`
	Data      string `dynamodbav:"event_data"`
	CreatedAt string `dynamodbav:"created_at"`
}

const (
	idField        = "id"
	versionField   = "version"
	dateTimeFormat = "2006-01-02T15:04:05.000000000Z"
)
