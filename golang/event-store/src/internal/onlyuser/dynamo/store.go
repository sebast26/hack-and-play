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

func (s Store) Save(ctx context.Context, user onlyuser.User) error {
	changes := user.Changes
	if len(changes) == 0 {
		return nil // nothing to do
	}

	var items []dbEventItem
	for _, change := range changes {
		serializedChange, err := json.Marshal(change)
		if err != nil {
			return err
		}

		var item dbEventItem
		switch change.(type) {
		case onlyuser.UserCreated:
			item = dbEventItem{
				key: key{
					ID:      fmt.Sprintf("user-%s", user.ID),
					Version: 0, // TODO
				},
				Type: "UserCreated",
				Data: string(serializedChange),
			}
		case onlyuser.UserEmailChanged:
			item = dbEventItem{
				key: key{
					ID:      fmt.Sprintf("user-%s", user.ID),
					Version: 0, // TODO
				},
				Type: "UserEmailChanged",
				Data: string(serializedChange),
			}
		}

		items = append(items, item)
	}

	// TODO: db.AppendEvents(streamName, dbEventItems);
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

// IsConditionalCheckFailed checks if generic error is AWS specific one
// for types.ConditionalCheckFailedException.
func IsConditionalCheckFailed(err error) bool {
	var conditionalCheckError *types.ConditionalCheckFailedException
	return errors.As(err, &conditionalCheckError)
}

func (s Store) Load(ctx context.Context, userID string) onlyuser.User {
	streamName := fmt.Sprintf("user-%s", userID)

	dbEvents, err := s.readEvents(ctx, streamName)
	if err != nil {
		return onlyuser.User{}
	}
	if len(dbEvents) == 0 {
		return onlyuser.User{} // TODO: is it properly handled? how to handle it?
	}

	var events []interface{}
	for _, dbEvent := range dbEvents {
		if dbEvent.Type == "UserCreated" {
			var e onlyuser.UserCreated
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				// TODO
				return onlyuser.User{}
			}
			events = append(events, e)
		}
		if dbEvent.Type == "UserEmailChanged" {
			var e onlyuser.UserEmailChanged
			err := json.Unmarshal([]byte(dbEvent.Data), &e)
			if err != nil {
				// TODO!!
				return onlyuser.User{}
			}
			events = append(events, e)
		}
	}

	var user = onlyuser.User{}
	for _, event := range events {
		user.When(event)
	}
	return user
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

		ConsistentRead: aws.Bool(true),
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
