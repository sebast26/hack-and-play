package onlyuser

type User struct {
	Name    string
	Email   string
	Address string
}

func (u User) ChangeEmail(email string) {
	// this is an old way
	u.Email = email
}
