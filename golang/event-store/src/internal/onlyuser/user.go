package onlyuser

import (
	"github.com/google/uuid"
	"sgorecki.me/golang/event-store/src/internal/es"
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

type UserEmailChanged struct {
	UserID string
	Email  string
}

func (u *User) ChangeEmail(email string) {
	u.Apply(u, UserEmailChanged{
		UserID: u.ID,
		Email:  email,
	})
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
