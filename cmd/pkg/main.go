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

	"github.com/thearyanahmed/nordsec/pkg/config"
	"github.com/thearyanahmed/nordsec/pkg/handler"
	"github.com/thearyanahmed/nordsec/pkg/logger"
	"github.com/thearyanahmed/nordsec/pkg/service"
)

func main() {
	conf, err := config.FromENV()

	if err != nil {
		log.Fatal(err)
	}

	logger.Setup(conf)

	svcAggregator, err := service.NewServiceAggregator(conf, logger.Logger())

	if err != nil {
		log.Fatal(err)
	}

	router := handler.NewRouter(conf, svcAggregator, logger.Logger())

	address := conf.AppAddress()

	server := &http.Server{Addr: address, Handler: router}

	serverCtx, cancelFunc := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		// @todo should be configurable from .env
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		cancelFunc()
	}()

	// Run the server

	fmt.Printf("will be serving on: %s\n", conf.AppAddress())

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
