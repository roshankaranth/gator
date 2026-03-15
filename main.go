package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/roshankaranth/gator/internal/config"
	"github.com/roshankaranth/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {

	cfg, err := config.Read()

	if err != nil {
		log.Fatalf("%v", err)
	}

	dbURL := cfg.Db_url
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("%v", err)
	}

	dbQueries := database.New(db)

	programState := &state{
		cfg: &cfg,
		db:  dbQueries,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	err = cmds.register("login", handlerLogin)

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("register", handlerRegister)

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("reset", handlerReset)

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("users", handlerUsers)

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("agg", handlerAggregate)

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("feeds", handlerFeeds)

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("follow", middlewareLoggedIn(handlerFollow))

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("following", middlewareLoggedIn(handlerFollowing))

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	if err != nil {
		log.Fatal(err)
	}

	err = cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("Not enough args provided!\n")
	}

	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
