package main

import (
	"database/sql"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	//Read existing config
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("read config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.Dburl)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	st := &state{cfg: &cfg, db: dbQueries}

	cmds := &commands{
		m: make(map[string]func(*state, command) error),
	}

	cmds.register("login", middlewareLoggedIn(handlerLogin))
	cmds.register("register", middlewareLoggedIn(handlerRegister))
	cmds.register("reset", middlewareLoggedIn(handlerReset))
	cmds.register("users", middlewareLoggedIn(handlerGetAllUsers))
	cmds.register("agg", middlewareLoggedIn(handlerGetAllAggregations))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", middlewareLoggedIn(handlerGetAllFeeds))
	cmds.register("follow", middlewareLoggedIn(handlerfollowFeed))
	cmds.register("following", middlewareLoggedIn(handlerfollowingFeeds))
	cmds.register("unfollow", middlewareLoggedIn(handlerunfollowFeed))
	cmds.register("browse", middlewareLoggedIn(handleGetFeedFromPosts))

	args := os.Args

	if len(args) < 2 {
		log.Fatalf("no command provided")
		os.Exit(1)
	}

	cmdName := args[1]
	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:]
	}
	cmd := command{name: cmdName, args: cmdArgs}

	if err := cmds.run(st, cmd); err != nil {
		log.Fatalf("run command: %v", err)
	}
}
