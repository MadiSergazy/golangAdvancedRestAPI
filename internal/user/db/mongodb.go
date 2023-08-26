package db

import (
	"context"
	"errors"
	"fmt"
	"mado/internal/user"
	"mado/pkg/logging"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	errCreateUser      = errors.New("failed to create user ")
	errConvertObjectID = errors.New("failed to create user ")
	errUserNotFound    = errors.New("not found user")
	errEntityNotFound  = errors.New("err entity not dound")
	errNotFound        = errors.New("not found")
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Create user")

	nCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	res, err := d.collection.InsertOne(nCtx, user)
	if err != nil {
		return "", errCreateUser
	}

	d.logger.Debug("converted insertedID to odjectID")
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", errConvertObjectID
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("Failed to convert hex to objectID hex: %s", id)
	}

	filter := bson.M{"_id": oid}

	res := d.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		// TODO 404

		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("err entity not dound")
		}
		return u, fmt.Errorf("Failed to find user by id: %s", id)
	}

	if err = res.Decode(&u); err != nil {
		return u, fmt.Errorf("Failed to decode user by id: %s from BD due to error: %s", id, err)
	}
	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID ID=%s", user.ID)
	}

	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user, err: %v", err)
	}

	var updateUserObj bson.M

	if err = bson.Unmarshal(userBytes, &updateUserObj); err != nil {
		return fmt.Errorf("failed to unmarshal user bytes, err: %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	res, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query, err: %v", err)
	}

	if res.MatchedCount == 0 {
		return errUserNotFound
	}

	d.logger.Tracef("Matched %d documents and modified %d documents", res.MatchedCount, res.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID ID=%s", id)
	}

	filter := bson.M{"_id": objectID}

	res, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID ID=%s", id)
	}
	if res.DeletedCount == 0 {
		return errNotFound
	}
	d.logger.Tracef("Deleted %d documents", res.DeletedCount)
	return nil
}
