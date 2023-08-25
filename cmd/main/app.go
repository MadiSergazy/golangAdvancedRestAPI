package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"mado/internal/user"
	"mado/pkg/logging"
)

func main() {
	logger := logging.GetLogger() //^ we call out logging package there (outside of it is own package and because of it func init will call automaticly)
	logger.Info("Create router")
	router := httprouter.New()

	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {

	logger := logging.GetLogger()
	logger.Info("Start the application")

	//for working with websockets in future because websockets based on tcp protocol
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Started listening server on port")
}
