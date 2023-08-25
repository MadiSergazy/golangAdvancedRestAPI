package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"mado/internal/user"
)

func main() {
	router := httprouter.New()
	handler := user.NewHandler()
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
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
		log.Fatalln(err)
	}
	fmt.Println("Started listening server on port")
}
