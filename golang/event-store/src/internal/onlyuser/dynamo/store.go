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
	"sgorecki.me/golang/event-store/src/internal/onlyuser"
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

func (s Store) Load(ctx context.Context, userID string) (onlyuser.User, error) {
	streamName := fmt.Sprintf("user-%s", userID)
	dbEvents, err := s.readEvents(ctx, streamName)
	if err != nil {
		return onlyuser.User{}, fmt.Errorf("%v: cannot read events", err)
	}
	if len(dbEvents) == 0 {
		return onlyuser.User{}, nil // TODO: is it properly handled? how to handle it?
	}

	events, err := loadEvents(dbEvents)
	if err != nil {
		return onlyuser.User{}, fmt.Errorf("%v: cannot load events", err)
	}
	var user = onlyuser.User{}
	for _, event := range events {
		user.When(event)
	}
	return user, nil
}

func (s Store) Save(ctx context.Context, user onlyuser.User) error {
	if len(user.Changes) == 0 {
		return nil // nothing to do
	}

	dbItems, err := toDBItems(user, user.Changes)
	if err != nil {
		return fmt.Errorf("%v: error converting to DB items", err)
	}

	return s.appendEvents(ctx, dbItems)
}

func toDBItems(user onlyuser.User, changes []interface{}) ([]dbEventItem, error) {
	var items []dbEventItem
	for i, change := range changes {
		serializedChange, err := json.Marshal(change)
		if err != nil {
			return nil, err
		}

		var item dbEventItem
		switch change.(type) {
		case onlyuser.UserCreated:
			key := toKey(user, i)
			item = dbEventItem{
				key:  key,
				Type: "UserCreated",
				Data: string(serializedChange),
			}
		case onlyuser.UserEmailChanged:
			key := toKey(user, i)
			item = dbEventItem{
				key:  key,
				Type: "UserEmailChanged",
				Data: string(serializedChange),
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func loadEvents(dbEvents []dbEventItem) ([]interface{}, error) {
	var events []interface{}
	for _, dbEvent := range dbEvents {
		if dbEvent.Type == "UserCreated" {
			var e onlyuser.UserCreated
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
		if dbEvent.Type == "UserEmailChanged" {
			var e onlyuser.UserEmailChanged
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}
	}
	return events, nil
}

// IsConditionalCheckFailed checks if generic error is AWS specific one
// for types.ConditionalCheckFailedException.
func IsConditionalCheckFailed(err error) bool {
	var conditionalCheckError *types.ConditionalCheckFailedException
	return errors.As(err, &conditionalCheckError)
}

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

func (s Store) appendEvents(ctx context.Context, items []dbEventItem) error {
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

		if IsConditionalCheckFailed(err) {
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

func toKey(user onlyuser.User, i int) key {
	return key{
		ID:      fmt.Sprintf("user-%s", user.ID),
		Version: user.Version + i + 1,
	}
}
