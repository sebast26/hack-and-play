package service

import (
	"context"

	"sgorecki.me/golang/event-store/src/internal/onlyuser/dynamo"
)

type UserService struct {
	store dynamo.UserStore
}

// Handle shows how the APIs for Load/Save from store does not change when
// switching from state-persistence into event sourcing
func (us UserService) Handle(ctx context.Context, command ChangeUserEmailCommand) {
	// this should stay the same both in an old and a new way
	user, err := us.store.Load(ctx, command.userID)
	if err != nil {
		panic("cannot load user")
	}
	newEmail := command.email
	user.ChangeEmail(newEmail)
	us.store.Save(ctx, user)
}

type ChangeUserEmailCommand struct {
	userID string
	email  string
}
