package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	defaultPortNumber = 9999
)

func main() {
	portNo := defaultPortNumber

	if len(os.Args) > 1 {
		// check if '-h' or '--help' was given
		for _, arg := range os.Args[1:] {
			if arg == "-h" || arg == "--help" {
				printUsage(0)
			}
		}

		var err error
		port := os.Args[1]
		if portNo, err = strconv.Atoi(port); err != nil {
			log.Printf("failed to parse port number: %s", err)

			printUsage(1)
		}
	}

	runHttp(portNo)
}

// print usage and exit with error code: `errorCode`
func printUsage(errorCode int) {
	fmt.Printf(`Usage:
	# print this help message
	$ %[2]s -h
	$ %[2]s --help

	# run http server on default port: %[1]d
	$ %[2]s

	# run http server on port number: PORT_NUMBER
	$ %[2]s PORT_NUMBER
`, defaultPortNumber, os.Args[0])

	os.Exit(errorCode)
}

// run http server on port number: `portNo`
func runHttp(portNo int) {
	http.HandleFunc("/", hello)

	addr := fmt.Sprintf(":%d", portNo)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Printf("failed to listen and serve: %s", err)
	}
}

// response with 'hello'
func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello\n")
}
