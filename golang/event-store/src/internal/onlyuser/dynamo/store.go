package dynamo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func (s Store) Save(user onlyuser.User) error {
	changes := user.Changes
	if len(changes) == 0 {
		return nil // nothing to do
	}

	var dbEvents []dbEvent
	for _, event := range changes {
		serializedEvent, err := json.Marshal(event)
		if err != nil {
			return err
		}
		dbEvents = append(dbEvents, dbEvent{
			key: key{
				ID:      fmt.Sprintf("user-%s", user.ID),
				Version: 0, // TODO: how to check version?
			},
			Type: "TODO",
			Data: string(serializedEvent),
		})

		// TODO: should be included in function parameters
		ctx := context.Background()

		// TODO: db.AppendEvents(streamName, dbEvents);
		for _, e := range dbEvents {
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
	}

	return nil
}

// IsConditionalCheckFailed checks if generic error is AWS specific one
// for types.ConditionalCheckFailedException.
func IsConditionalCheckFailed(err error) bool {
	var conditionalCheckError *types.ConditionalCheckFailedException
	return errors.As(err, &conditionalCheckError)
}

func (s Store) Load(userID string) onlyuser.User {
	//streamName := fmt.Sprintf("user-%s", userID)

	// TODO: db.ReadEvents(streamName);
	var dbEvents []dbEvent
	if len(dbEvents) == 0 {
		return onlyuser.User{} // TODO: is it properly handled? how to handle it?
	}

	var user = onlyuser.User{}
	for _, event := range dbEvents {
		user.When(event)
	}
	return user
}

type key struct {
	ID      string `dynamodbav:"id"`
	Version int    `dynamodbav:"version"`
}

// dbEvent is dynamoDB struct for event.
type dbEvent struct {
	key

	Type string `dynamodbav:"event_type"`
	Data string `dynamodbav:"event_data"`
}
