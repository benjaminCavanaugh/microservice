package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	fileName := "config.yaml";
	// TODO: Error handling here.
	config, _ := ParseFromFile(fileName)
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

	var handler helloWorldhandler
	
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




/*

Contents from config.go.
These are here because I was unable to get importing working.

*/

//var fileName string

type Config struct {
	serviceConfig ServiceConfig
}

func NewConfig() Config {
	var c Config;
	return c;
}

func ParseFromFile(filename string) (*Config, error) {
    buf, err := ioutil.ReadFile(filename)

	if err != nil {
        return nil, err
    }

    c := &Config{}
    err = yaml.Unmarshal(buf, c)

	if err != nil {
        return nil, fmt.Errorf("in file %q: %w", filename, err)
    }

    return c, err
}

func (c Config) GetServiceConfig() (ServiceConfig, error) {
	// if c.serviceConfig == nil {
	// 	return nil, error.Error()
	// }

	return c.serviceConfig, nil
}

type ServiceConfig struct {
	htmlServerConfig HtmlServerConfig
}

func NewServiceConfig() ServiceConfig {
	var s ServiceConfig;
	return s;
}

func (s ServiceConfig) GetHtmlServerConfig() HtmlServerConfig {
	return s.htmlServerConfig;
}

type HtmlServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewHtmlServerConfig() HtmlServerConfig {
	var h HtmlServerConfig;
	return h;
}


/*

Contents from handler.go.
These are here because I was unable to get importing working.

*/

type helloWorldhandler struct {
	message string
}

func (s helloWorldhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
     fmt.Fprintf(w, "HeloWorld")
}