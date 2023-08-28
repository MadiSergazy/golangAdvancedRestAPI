package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"

	author "mado/internal/author/db"
	"mado/internal/config"
	"mado/internal/user"
	"mado/pkg/client/postgresql"
	"mado/pkg/logging"
)

func main() {
	logger := logging.GetLogger() //^ we call out logging package there (outside of it is own package and because of it func init will call automaticly)
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}

	authorRepository := author.NewRepository(postgreSQLClient, logger)

	// cfgMongo := cfg.MongoDB
	// mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username, cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	// if err != nil {
	// 	panic(err)
	// }
	// storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	// user1 := user.User{
	// 	ID: "", Email: "myemail@example.com", Username: "mado", PasswordHash: "123456",
	// }
	// user1ID, err := storage.Create(context.TODO(), user1)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// logger.Info("user1ID: ", user1ID)

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
		// ^Unix sockets are used for faster communication between processes on the same computer
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			// logger.Info(err)
			logger.Fatal(err)
		}

	} else {

		logger.Info("Listen tcp socket")
		var err error
		// for working with websockets in future because websockets based on tcp protocol
		// TCP is used for reliable communication over a network
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
