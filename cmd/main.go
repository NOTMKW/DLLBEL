package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/NOTMKW/DLLBEL/internal/config"
	"github.com/NOTMKW/DLLBEL/internal/server"
)

func main () {
	cfg := config.Load()

	srv := server.NewServer(cfg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func () {
		<-c
		log.Println("Shutting down gracefully....")
		srv.Shutdown()
		os.Exit(0)
	}()

	log.Printf("starting MT5 websocket server on port %s", cfg.Port)
	if err := srv.Start(); err != nil{
		log.Fatal("failed to start server", err)
	}
}