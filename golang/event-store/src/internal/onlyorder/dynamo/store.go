package dynamo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"sgorecki.me/golang/event-store/src/internal/onlyorder"
)

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

func (s Store) Load(ctx context.Context, orderID string) (onlyorder.Order, error) {
	streamName := fmt.Sprintf("order-%s", orderID)
	dbEvents, err := s.readEvents(ctx, streamName)
	if err != nil {
		return onlyorder.Order{}, fmt.Errorf("%v: cannot read events", err)
	}
	if len(dbEvents) == 0 {
		return onlyorder.Order{}, nil // TODO: is it properly handled? how to handle it?
	}

	events, err := loadEvents(dbEvents)
	if err != nil {
		return onlyorder.Order{}, fmt.Errorf("%v: cannot load events", err)
	}
	var order = onlyorder.Order{}
	for _, event := range events {
		order.When(event)
	}
	return order, nil
}

func (s Store) Save(ctx context.Context, order onlyorder.Order) error {
	if len(order.Changes) == 0 {
		return nil
	}

	dbItems, err := toDBItems(order, order.Changes)
	if err != nil {
		return fmt.Errorf("%v: error converting to DB items", err)
	}

	return s.appendEvents(ctx, dbItems)
}

func loadEvents(dbEvents []dbEventItem) ([]interface{}, error) {
	var events []interface{}
	for _, dbEvent := range dbEvents {
		if dbEvent.Type == "OrderCreated" {
			var e onlyorder.OrderCreated
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
		if dbEvent.Type == "ItemAdded" {
			var e onlyorder.ItemAdded
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
	}
	return events, nil
}

func (s Store) appendEvents(ctx context.Context, items []dbEventItem) error {
	// TODO: should use transact_write_items
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

// TODO: good candidate to place this method on general store (abstracting away entities)
func (s Store) readEvents(ctx context.Context, streamName string) ([]dbEventItem, error) {
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

	var events []dbEventItem
	for _, rawItem := range out.Items {
		var item dbEventItem
		err = attributevalue.UnmarshalMap(rawItem, &item)
		if err != nil {
			return nil, err
		}
		events = append(events, item)
	}

	return events, nil
}

// isConditionalCheckFailed checks if generic error is AWS specific one
// for types.ConditionalCheckFailedException.
func isConditionalCheckFailed(err error) bool {
	var conditionalCheckError *types.ConditionalCheckFailedException
	return errors.As(err, &conditionalCheckError)
}

func toDBItems(order onlyorder.Order, changes []interface{}) ([]dbEventItem, error) {
	var items []dbEventItem
	for i, change := range changes {
		serializedChange, err := json.Marshal(change)
		if err != nil {
			return nil, err
		}

		var item dbEventItem
		switch change.(type) {
		case onlyorder.OrderCreated:
			key := toKey(order, i)
			item = dbEventItem{
				key:  key,
				Type: "OrderCreated",
				Data: string(serializedChange),
			}
		case onlyorder.ItemAdded:
			key := toKey(order, i)
			item = dbEventItem{
				key:  key,
				Type: "ItemAdded",
				Data: string(serializedChange),
			}
		}

		items = append(items, item)
	}
	return items, nil
}

func toKey(order onlyorder.Order, i int) key {
	return key{
		ID:      fmt.Sprintf("order-%s", order.ID),
		Version: order.Version + i + 1,
	}
}

type key struct {
	ID      string `dynamodbav:"id"`
	Version int    `dynamodbav:"version"`
}

// dbEventItem is dynamoDB struct for event.
type dbEventItem struct {
	key

	Type string `dynamodbav:"event_type"`
	Data string `dynamodbav:"event_data"`
}
