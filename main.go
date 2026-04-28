package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

const (
	defaultPortNumber = 9999
	minPortNumber     = 1
	maxPortNumber     = 65535

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 3 * time.Second
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
		if portNo < minPortNumber || portNo > maxPortNumber {
			log.Printf("port number out of range (%d-%d): %d", minPortNumber, maxPortNumber, portNo)

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
`, defaultPortNumber, filepath.Base(os.Args[0]))

	os.Exit(errorCode)
}

// run http server on port number: `portNo`
func runHttp(portNo int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", portNo),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if err != nil {
			log.Printf("failed to listen and serve: %s", err)
		}
	case sig := <-stop:
		log.Printf("received signal: %s, shutting down", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("failed to shut down gracefully: %s", err)
		}
	}
}

// handle requests to "/*"
func hello(w http.ResponseWriter, r *http.Request) {
	// all other routes other than "/": http 404 error
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// respond with 'hello'
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := io.WriteString(w, "hello\n"); err != nil {
		log.Printf("failed to write hello: %s", err)
	}
}
