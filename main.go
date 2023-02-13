package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	helper "microservice/src/config"
	handler "microservice/src/handler"
	connector "microservice/src/service"
)

func main() {
	helper.HelloWorld();

	db, connectionError := connector.Connect();

	if( connectionError != nil) {
		// TODO: Handle this error here.
	}

	defer db.Close();
	connector.QueryUsers(db);
	
	fileName := "config.yaml";
	// TODO: Error handling here.
	config, _ := helper.ParseFromFile(fileName)
	fmt.Println(config);

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serviceConfig, err := config.GetServiceConfig()

	if err != nil {
		log.Fatal(err)
	}

  	// TODO: set up logging

	// set up the service from config
	htmlServerConfig := serviceConfig.GetHtmlServerConfig()
	fmt.Println(htmlServerConfig);

	// FIXME: Make this use the values read-in from the HtmlServerConfig instead of hard-coding them here.
	address := "localhost:3333"
	var readTimeout time.Duration = 55000000000
	var writeTimeout time.Duration = 55000000000
	// address := htmlServerConfig.Addr
	// var readTimeout time.Duration = htmlServerConfig.ReadTimeout * time.Second
	// var writeTimeout time.Duration = htmlServerConfig.WriteTimeout * time.Second

	var handler handler.HelloWorldhandler
	
	htmlServer := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	// trap lifecycle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on %s", address)
		err := htmlServer.ListenAndServe()

		if err == http.ErrServerClosed {
			log.Print(err)
		} else {
			log.Fatal(err)
		}
		
		signals <- syscall.SIGQUIT
	}()

	sig := <-signals

	switch sig {
		case os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			err := htmlServer.Shutdown(ctx)
			if err != nil {
				log.Print(err)
			}
	}
}