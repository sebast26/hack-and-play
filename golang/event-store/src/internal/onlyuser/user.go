package onlyuser

import "github.com/google/uuid"

type User struct {
	// TODO: this should be generic! belongs to Entity class
	Changes []interface{}

	ID    string
	Name  string
	Email string
}

func NewUser(name, email string) User {
	user := User{}
	id := uuid.New()
	user.Apply(UserCreated{
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
	// this is important, to not mutate state here!!
	newEmail := email

	u.Apply(UserEmailChanged{
		UserID: u.ID,
		Email:  newEmail,
	})
}

// TODO: this should be generic! belongs to Entity class
func (u *User) Apply(event interface{}) {
	u.When(event)
	u.Changes = append(u.Changes, event)
}

// TODO: this should be abstract method reimplemented in all subclasses!
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
