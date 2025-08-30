package main

import (
	"fmt"
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
	conf, err := config.NewConfig()
	if err != nil || conf == nil {
		log.Fatal(err)
	}

	// Open the database connection
	dbConn, err := config.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Start the IMAP client
	fmt.Println("Connecting to server...")
	c, err := imapclient.Connect(conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	fmt.Println("Connected and logged in!")

	go service.StartEmailScraper(5, c, dbConn)
	api.InitAPI(dbConn)
	go api.StartAPIService()

	// Wait for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
