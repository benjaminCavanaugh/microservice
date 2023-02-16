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

	configuration "microservice/src/config"
	handler "microservice/src/handler"
	connector "microservice/src/service"
)

func main() {
	fileName := "config.yaml";
	// TODO: Error handling here.
	config, _ := configuration.ParseFromFile(fileName)
	fmt.Printf("Full yaml.conf file:\n%v\n", config);

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serviceConfig, err := config.GetServiceConfig()

	if err != nil {
		log.Fatal(err)
	}

  	// TODO: set up logging


	// Connect to the database.

	databaseConfig := serviceConfig.GetDatabaseConfig()
	fmt.Printf("Database config:\n%v\n\n\n", databaseConfig);

	// FIXME: Make this use the values read-in from the DatabaseConfig instead of hard-coding them here.
	databaseConfig.Host= "localhost"
	databaseConfig.Port= 5432
	databaseConfig.User= "postgres"
	databaseConfig.Password= "password"
	databaseConfig.DBname= "golang-db"

	var connection connector.Connection = connector.Connect(databaseConfig);

	if connection.Err != nil {
		// TODO: Handle this error here.
	}

	defer connection.Database.Close();
	fmt.Println(connection.QueryUsersByName("a"));


	// set up the service from config

	htmlServerConfig := serviceConfig.GetHtmlServerConfig()
	fmt.Printf("\n\nServer config:\n%v\n", htmlServerConfig);

	// FIXME: Make this use the values read-in from the HtmlServerConfig instead of hard-coding them here.
	address := "localhost:3333"
	var readTimeout time.Duration = 55000000000
	var writeTimeout time.Duration = 55000000000
	// address := htmlServerConfig.Addr
	// var readTimeout time.Duration = htmlServerConfig.ReadTimeout * time.Second
	// var writeTimeout time.Duration = htmlServerConfig.WriteTimeout * time.Second

	handler := handler.NewHandler(connection);
	
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