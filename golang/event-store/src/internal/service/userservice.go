package service

import (
	"context"

	"sgorecki.me/golang/event-store/src/internal/onlyuser/dynamo"
)

type UserService struct {
	store dynamo.Store
}

func NewUserService(store dynamo.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (us UserService) Handle(ctx context.Context, command ChangeUserEmailCommand) {
	// this should stay the same both in an old and a new way
	user := us.store.Load(ctx, command.userID)
	newEmail := command.email
	user.ChangeEmail(newEmail)
	us.store.Save(ctx, user)
}

type ChangeUserEmailCommand struct {
	userID string
	email  string
}
