package user

import (
	"errors"
	"sgorecki.me/golang/event-store/src/internal/es"
)

// User represents a user entity.
type User struct {
	es.Applier
	Name    string
	Email   string
	Address string
}

func (u User) AddUser(name, email, address string) error {
	if !canAddUser(email) {
		return errors.New("user already exists!")
	}

	u.Name = name
	u.Email = email
	u.Address = address

	u.Apply(UserCreated{
		Name:    name,
		Email:   email,
		Address: address,
	})

	return nil
}

//func (u User) When(event interface{}) {
//
//}

func canAddUser(email string) bool {
	// TODO
	return true
}

func testInh(user es.Applier) {

}

func testUserInh() {
	user := User{
		Name:    "Sebastian",
		Email:   "sgorecki@gmail",
		Address: "Some address",
	}

	// todo: how does it work? does it work? When should be implemented, but it is not
	testInh(user)
}
