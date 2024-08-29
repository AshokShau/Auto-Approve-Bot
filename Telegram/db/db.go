/*
 * Copyright © 2024 AshokShau <github.com/AshokShau>
 */

package db

import (
	"context"
	"github.com/AshokShau/Auto-Approve-Bot/Telegram/config"

	"github.com/go-faster/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// Todo: improve more, Ofc i can do better but i'm lazy.

var (
	tdCtx       = context.TODO()
	mongoClient *mongo.Client

	chatColl    *mongo.Collection
	userColl    *mongo.Collection
	disableColl *mongo.Collection
)

func init() {
	var err error

	mongoClient, err = mongo.Connect(tdCtx, options.Client().ApplyURI(config.DbUrl))
	if err != nil {
		log.Fatalf("Error while connecting to MongoDB: %v", err)
	}

	db := mongoClient.Database(config.DbName)
	chatColl = db.Collection("chats")
	userColl = db.Collection("users")
	disableColl = db.Collection("disabled")
}

// Utility function to find a single document
func findOne(collection *mongo.Collection, filter bson.M) *mongo.SingleResult {
	return collection.FindOne(tdCtx, filter)
}

// GetServedChats Gets all served chats with chat_id < 0
func GetServedChats() ([]bson.M, error) {
	cursor, err := chatColl.Find(tdCtx, bson.M{"chat_id": bson.M{"$lt": 0}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while retrieving served chats")
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, tdCtx)

	var chats []bson.M
	for cursor.Next(tdCtx) {
		var chat bson.M
		if err = cursor.Decode(&chat); err != nil {
			return nil, errors.Wrap(err, "Error while decoding chat")
		}
		chats = append(chats, chat)
	}

	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error while iterating over cursor")
	}

	return chats, nil
}

// IsServedChat Checks if a chat is served
func IsServedChat(chatID int64) (bool, error) {
	var result bson.M
	err := findOne(chatColl, bson.M{"chat_id": chatID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil // Chat not found
		}
		return false, errors.Wrap(err, "Error while checking if chat is served")
	}
	return true, nil
}

// AddServedChat Adds a chat as served
func AddServedChat(chatID int64) error {
	isServed, err := IsServedChat(chatID)
	if err != nil {
		return err
	}
	if isServed {
		return nil // Chat already served
	}

	_, err = chatColl.InsertOne(tdCtx, bson.M{"chat_id": chatID})
	if err != nil {
		return errors.Wrap(err, "Error while adding served chat")
	}

	return nil
}

// GetServedUsers Retrieves all served users
func GetServedUsers() ([]bson.M, error) {
	cursor, err := userColl.Find(tdCtx, bson.M{"user_id": bson.M{"$gt": 0}})
	if err != nil {
		return nil, errors.Wrap(err, "Error while retrieving served users")
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		_ = cursor.Close(ctx)
	}(cursor, tdCtx)

	var users []bson.M
	for cursor.Next(tdCtx) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			return nil, errors.Wrap(err, "Error while decoding user")
		}
		users = append(users, user)
	}

	if err = cursor.Err(); err != nil {
		return nil, errors.Wrap(err, "Error while iterating over cursor")
	}

	return users, nil
}

// IsServedUser Checks if a user is served
func IsServedUser(userID int64) (bool, error) {
	var result bson.M
	err := findOne(userColl, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil // User not found
		}
		return false, errors.Wrap(err, "Error while checking if user is served")
	}
	return true, nil
}

// AddServedUser Adds a user as served
func AddServedUser(userID int64) error {
	isServed, err := IsServedUser(userID)
	if err != nil {
		return err
	}
	if isServed {
		return nil // User already served
	}

	_, err = userColl.InsertOne(tdCtx, bson.M{"user_id": userID})
	if err != nil {
		return errors.Wrap(err, "Error while adding served user")
	}

	return nil
}

// GetChatCount Gets the count of chats
func GetChatCount() (int, error) {
	count, err := chatColl.CountDocuments(tdCtx, bson.M{})
	if err != nil {
		return 0, errors.Wrap(err, "Error while counting chats")
	}
	return int(count), nil
}

// GetUserCount Gets the count of users
func GetUserCount() (int, error) {
	count, err := userColl.CountDocuments(tdCtx, bson.M{})
	if err != nil {
		return 0, errors.Wrap(err, "Error while counting users")
	}
	return int(count), nil
}

// IsDisabledChat Checks if a chat is disabled
func IsDisabledChat(chatID int64) (bool, error) {
	var result bson.M
	err := findOne(disableColl, bson.M{"chat_id": chatID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil // Chat not found
		}
		return true, errors.Wrap(err, "Error while checking if chat is disabled")
	}

	return false, nil
}

// DisableApprove Disables auto-approving for a chat
func DisableApprove(chatID int64) error {
	isServed, err := IsDisabledChat(chatID)
	if err != nil {
		return err
	}

	if isServed {
		return nil // Already disabled
	}

	_, err = disableColl.InsertOne(tdCtx, bson.M{"chat_id": chatID})
	if err != nil {
		return errors.Wrap(err, "Error while adding disabled chat")
	}

	return nil
}

// EnableApprove Enables auto-approving for a chat
func EnableApprove(chatID int64) error {
	isServed, err := IsDisabledChat(chatID)
	if err != nil {
		return err
	}
	if !isServed {
		return nil // Already enabled
	}

	_, err = disableColl.DeleteOne(tdCtx, bson.M{"chat_id": chatID})
	if err != nil {
		return errors.Wrap(err, "Error while enabling chat")
	}

	return nil
}
