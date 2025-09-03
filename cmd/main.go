package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samothreesixty/EmailScraper/internal/api"
	"github.com/samothreesixty/EmailScraper/internal/config"
	"github.com/samothreesixty/EmailScraper/internal/imapclient"
	"github.com/samothreesixty/EmailScraper/internal/service"
)

func main() {
	// Open the database connection
	dbConn, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}

	c, err := imapclient.NewGmailClient()
	if err != nil {
		log.Fatal(err)
	}

	go service.StartEmailScraper(5, c, dbConn)
	api.InitAPI(dbConn)
	go api.StartAPIService()

	// Wait for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
