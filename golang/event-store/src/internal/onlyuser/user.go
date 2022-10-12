package onlyuser

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"sgorecki.me/golang/event-store/src/internal/es"
	eventstore "sgorecki.me/golang/event-store/src/internal/es/dynamo"
)

type User struct {
	es.Entity
	Name  string
	Email string
}

func NewUser(name, email string) User {
	user := User{}
	id := uuid.New()
	user.Apply(&user, UserCreated{
		ID:    id.String(),
		Name:  name,
		Email: email,
	})
	return user
}

type UserCreated struct {
	ID    string
	Name  string
	Email string
}

func NewUserCreated(data string) (UserCreated, error) {
	var e UserCreated
	err := json.Unmarshal([]byte(data), &e)
	if err != nil {
		return UserCreated{}, err
	}
	return e, nil
}

type UserEmailChanged struct {
	UserID string
	Email  string
}

func NewUserEmailChanged(data string) (UserEmailChanged, error) {
	var e UserEmailChanged
	err := json.Unmarshal([]byte(data), &e)
	if err != nil {
		return UserEmailChanged{}, err
	}
	return e, nil
}

func (u *User) ChangeEmail(email string) {
	u.Apply(u, UserEmailChanged{
		UserID: u.ID,
		Email:  email,
	})
}

func (u *User) When2(event eventstore.DBEventItem) error {
	switch event.Type {
	case "UserCreated":
		v, err := NewUserCreated(event.Data)
		if err != nil {
			return err
		}
		u.ID = v.ID
		u.Name = v.Name
		u.Email = v.Email
	case "UserEmailChanged":
		v, err := NewUserEmailChanged(event.Data)
		if err != nil {
			return err
		}
		u.Email = v.Email
	default:
		return errors.New("unknown event type")
	}
	return nil
}

func (u *User) When(event interface{}) {
	switch v := event.(type) {
	case UserCreated:
		u.ID = v.ID
		u.Name = v.Name
		u.Email = v.Email
	case UserEmailChanged:
		u.Email = v.Email
	default:
		panic("unknown event type!")
	}
}
