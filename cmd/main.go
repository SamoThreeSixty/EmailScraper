package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/samothreesixty/EmailScraper/internal/config"
	"github.com/samothreesixty/EmailScraper/internal/imapclient"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil || conf == nil {
		log.Fatal(err)
	}

	fmt.Println("Connecting to server...")
	c, err := imapclient.Connect(conf.Host, conf.Port, conf.Username, conf.Password)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()
	fmt.Println("Connected and logged in!")

	// Wait for shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
