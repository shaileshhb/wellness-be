package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/shaileshhb/websockets/db"
	"github.com/shaileshhb/websockets/log"
	"github.com/shaileshhb/websockets/server"
)

func main() {
	logger := log.InitializeLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error loading.env file")
		return
	}

	database := db.NewDatabase()
	ser := server.NewServer(logger, database)
	ser.InitializeRouter()

	ser.RegisterModuleRoutes()

	logger.Error().Err(ser.App.Listen(fmt.Sprintf(":%s", os.Getenv("PORT")))).Msg("")

	// Stop Server On System Call or Interrupt.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	os.Exit(0)
}
