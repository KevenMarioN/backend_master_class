package main

import (
	"database/sql"
	"log"

	"github.com/kevenmarion/backend_master_class/api"
	db "github.com/kevenmarion/backend_master_class/db/sqlc"
	"github.com/kevenmarion/backend_master_class/util"
	_ "github.com/lib/pq"
)

func main() {
	var (
		conn   *sql.DB
		err    error
		config *util.Config
	)

	if config, err = util.LoadConfig("./"); err != nil {
		log.Fatal("cannot load config:", err)
	}

	if conn, err = sql.Open(config.DBDRIVER, config.DBSource); err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	server := api.NewServer(db.NewStore(conn))
	server.LoadRouters()
	server.LoadValidators()

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

}
