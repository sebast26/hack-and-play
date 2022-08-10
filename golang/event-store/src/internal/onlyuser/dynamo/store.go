package dynamo

import (
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

func (s Store) Load(userID string) onlyuser.User {
	return onlyuser.User{}
}

func (s Store) Save(user onlyuser.User) {

}

type key struct {
	ID      string `dynamodbav:"id"`
	Version int    `dynamodbav:"version"`
}

// eventItem is dynamoDB struct for event.
type eventItem struct {
	key

	Type string `dynamodbav:"event_type"`
	Data string `dynamodbav:"event_data"`
}
