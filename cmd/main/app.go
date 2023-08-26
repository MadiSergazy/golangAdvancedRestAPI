package main

import (
	"fmt"
	"log"
	"mado/internal/config"
	"mado/internal/user"
	"mado/pkg/logging"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger() //^ we call out logging package there (outside of it is own package and because of it func init will call automaticly)
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("Start the application")
	var listener net.Listener
	if cfg.Listen.Type == "sock" {
		// /path/to/binary
		// Dir() -- path/to
		// Abs() absolute path
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			// logger.Info(err)
			logger.Fatal(err)
		}

		logger.Info("Create socket")

		socketPath := path.Join(appDir, "app.sock")

		logger.Info("Listen unix socket %s", socketPath)
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			// logger.Info(err)
			logger.Fatal(err)
		}

	} else {

		logger.Info("Listen tcp socket")
		var err error
		// for working with websockets in future because websockets based on tcp protocol
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		if err != nil {
			log.Fatal(err)
		}
		logger.Info("Started listening server on port ", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))

	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := server.Serve(listener); err != nil {
		logger.Fatal(err)
	}
}
