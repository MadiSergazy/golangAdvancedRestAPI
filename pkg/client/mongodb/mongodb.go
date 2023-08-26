package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, host, port, username, password, database, authDB string) (db *mongo.Database, err error) {
	var mongoDBURL string
	var isAuth bool
	if username == "" && password == "" {
		mongoDBURL = fmt.Sprintf("mongodb://%s:%s", host, port)
		isAuth = false
	} else {
		mongoDBURL = fmt.Sprintf("mongodb://%s%s@%s:%s", username, password, host, port)
		isAuth = true
	}

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	if isAuth {
		if authDB == "" {
			authDB = database
		}
		clientOptions.SetAuth(options.Credential{
			// AuthMechanism:           nil,
			// AuthMechanismProperties: nil,
			AuthSource: authDB,
			Username:   username,
			Password:   password,
			// PasswordSet:             nil,
		})
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New("Failed to connect to mongoDB: " + err.Error())
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, errors.New("Failed to ping in mongoDB: " + err.Error())
	}

	return client.Database(database), nil
}
