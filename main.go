package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal("Cannot read user config", err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	state := State{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := Commands{
		cmd: make(map[string]func(*State, Command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	name := os.Args[1]
	args := os.Args[2:]

	if err = cmds.run(&state, Command{name: name, args: args}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
