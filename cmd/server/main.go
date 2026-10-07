package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"wails-datastar-example/internal/server"
)

func main() {
	addr := flag.String("addr", "", "listen address (default DESKTOP_ADDR or :8080)")
	flag.Parse()

	listen := *addr
	if listen == "" {
		listen = os.Getenv("DESKTOP_ADDR")
	}
	if listen == "" {
		listen = ":8080"
	}

	token := os.Getenv("DESKTOP_TOKEN")
	if token == "" {
		token = "dev"
	}

	httpServer := &http.Server{
		Addr:              listen,
		Handler:           server.New(token),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("desktop server on %s", listen)
	log.Fatal(httpServer.ListenAndServe())
}
