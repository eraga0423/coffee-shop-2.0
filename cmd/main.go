package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

// Define command-line flags for directory path and port number
var (
	dir = flag.String("dir", "data", "Path to the directory")

	port = flag.String("port", "8080", "Port number")
)

func main() {
	flag.Parse()
	// Start the server with the specified port, handling any errors
	if *port == "0" {
		slog.Error("Port number must be between 1024 and 65535")
		os.Exit(1)
	}

	err := StartServer(*port)
	if err != nil {
		log.Fatal(err)
	}
}

// init sets up the usage information displayed when the help flag is used
func init() {
	flag.Usage = func() {
		fmt.Println(
			`Frappuccino Management System

Usage:
	frappuccino [--port <N>] [--dir <S>] 
	frappuccino --help
			
Options:
	--help       Show this screen.
	--port N     Port number.
	--dir S      Path to the data directory.`)
	}
}
