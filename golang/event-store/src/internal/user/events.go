package user

type UserCreated struct {
	Name    string
	Email   string
	Address string
}

type UserEmailUpdated struct {
	Email string
}
