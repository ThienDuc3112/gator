package main

import (
	"context"
	"fmt"

	"github.com/ThienDuc3112/gator/internal/database"
)

func middlewareLoggedIn(handler func(*State, Command, database.User) error) func(*State, Command) error {
	return func(state *State, cmd Command) error {
		user, err := state.db.GetUser(context.Background(), state.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("cannot get the current user: %v", err)
		}
		return handler(state, cmd, user)
	}
}
