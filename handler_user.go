package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(state *State, cmd Command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("expect arguments to not be empty")
	}

	username := cmd.args[0]

	user, err := state.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user \"%v\" not found", username)
	}

	if err := state.cfg.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("Username %v setted\n", user.Name)
	return nil
}

func handlerRegister(state *State, cmd Command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("expect arguments to not be empty")
	}

	newUser, err := state.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return err
	}
	err = state.cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%v have been registered", newUser.Name)
	fmt.Println(newUser)
	return nil
}

func handlerReset(state *State, _ Command) error {
	err := state.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Reset all users")
	return nil
}

func handlerUsers(state *State, _ Command) error {
	users, err := state.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		fmt.Printf("* %v", user.Name)
		if user.Name == state.cfg.CurrentUserName {
			fmt.Print(" (current)")
		}
		fmt.Println()
	}
	return nil
}
