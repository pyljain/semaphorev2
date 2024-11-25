package main

import (
	"flag"
	"fmt"
	"os"
	"semaright/internal/server"
)

func main() {
	var port int
	var connectionString string

	flag.IntVar(&port, "port", 8080, "port")
	flag.StringVar(&connectionString, "connection-string", "mongodb://127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019", "Connection string to mongodb")
	flag.Parse()

	s, err := server.New(port, connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create server: %v\n", err)
		os.Exit(1)
	}

	err = s.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run server: %v\n", err)
		os.Exit(1)
	}
}

/*
./sr -port=8080 -connectionStrings="db1.db,db2.db"
./sr -port=8081 -connectionStrings="db1.db,db2.db"
*/
