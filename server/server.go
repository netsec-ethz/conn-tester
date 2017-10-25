package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/netsec-ethz/conn-tester/lib/httputils"
)

func main() {
	fmt.Println("Starting server...")

	if len(os.Args) < 2 {
		fmt.Println("Error! You must specify servers port")
		os.Exit(-1)
	}

	rm := httputils.CreateHttpRequestMultiplexer(httputils.DefaultHostAddressExtractor)

	registerHandlers(rm)

	server := rm.StartHttpServer(os.Args[1])

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop // Wait for Ctrl+C

	fmt.Println("Shutting down server...")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	server.Shutdown(ctx)
}
