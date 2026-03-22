/*
 * Copyright (c) 2026. AshokShau <github.com/AshokShau>
 */

package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/AshokShau/Auto-Approve-Bot/src/config"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	tdCtx       = context.TODO()
	mongoClient *mongo.Client

	chatColl    *mongo.Collection
	userColl    *mongo.Collection
	disableColl *mongo.Collection
)

func Connect() error {
	var err error

	mongoClient, err = mongo.Connect(options.Client().ApplyURI(config.DbUrl))
	if err != nil {
		return fmt.Errorf("connect db err: %v", err)
	}

	db := mongoClient.Database(config.DbName)
	chatColl = db.Collection("chats")
	userColl = db.Collection("users")
	disableColl = db.Collection("disabled")
	return nil
}

// Utility function to find a single document
func findOne(collection *mongo.Collection, filter bson.M) *mongo.SingleResult {
	return collection.FindOne(tdCtx, filter)
}

// GetServedChats Gets all served chats with chat_id < 0
func GetServedChats() ([]bson.M, error) {
	cursor, err := chatColl.Find(tdCtx, bson.M{"chat_id": bson.M{"$lt": 0}})
	if err != nil {
		return nil, fmt.Errorf("error while retrieving served chats: %w", err)
	}

	defer func() { _ = cursor.Close(tdCtx) }()

	var chats []bson.M
	for cursor.Next(tdCtx) {
		var chat bson.M
		if err = cursor.Decode(&chat); err != nil {
			return nil, fmt.Errorf("error while decoding chat: %w", err)
		}
		chats = append(chats, chat)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over cursor: %w", err)
	}

	return chats, nil
}

// IsServedChat Checks if a chat is served
func IsServedChat(chatID int64) (bool, error) {
	var result bson.M
	err := findOne(chatColl, bson.M{"chat_id": chatID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("error while checking if chat is served: %w", err)
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
		return nil
	}

	_, err = chatColl.InsertOne(tdCtx, bson.M{"chat_id": chatID})
	if err != nil {
		return fmt.Errorf("error while adding served chat: %w", err)
	}

	return nil
}

// GetServedUsers Retrieves all served users
func GetServedUsers() ([]bson.M, error) {
	cursor, err := userColl.Find(tdCtx, bson.M{"user_id": bson.M{"$gt": 0}})
	if err != nil {
		return nil, fmt.Errorf("error while retrieving served users: %w", err)
	}
	defer func() { _ = cursor.Close(tdCtx) }()

	var users []bson.M
	for cursor.Next(tdCtx) {
		var user bson.M
		if err = cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("error while decoding user: %w", err)
		}
		users = append(users, user)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over cursor: %w", err)
	}

	return users, nil
}

// IsServedUser Checks if a user is served
func IsServedUser(userID int64) (bool, error) {
	var result bson.M
	err := findOne(userColl, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, fmt.Errorf("error while checking if user is served: %w", err)
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
		return nil
	}

	_, err = userColl.InsertOne(tdCtx, bson.M{"user_id": userID})
	if err != nil {
		return fmt.Errorf("error while adding served user: %w", err)
	}

	return nil
}

// GetChatCount Gets the count of chats
func GetChatCount() (int, error) {
	count, err := chatColl.CountDocuments(tdCtx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("error while counting chats: %w", err)
	}
	return int(count), nil
}

// GetUserCount Gets the count of users
func GetUserCount() (int, error) {
	count, err := userColl.CountDocuments(tdCtx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("error while counting users: %w", err)
	}
	return int(count), nil
}

// IsApproveEnabled Checks if a chat is enabled for auto-approve
func IsApproveEnabled(chatID int64) (bool, error) {
	var result bson.M
	err := findOne(disableColl, bson.M{"chat_id": chatID}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return true, nil
		}
		return false, fmt.Errorf("error while checking if chat is disabled: %w", err)
	}
	return false, nil
}

// DisableApprove Disables auto-approving for a chat
func DisableApprove(chatID int64) error {
	enabled, err := IsApproveEnabled(chatID)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}

	_, err = disableColl.InsertOne(tdCtx, bson.M{"chat_id": chatID})
	if err != nil {
		return fmt.Errorf("error while adding disabled chat: %w", err)
	}

	return nil
}

// EnableApprove Enables auto-approving for a chat
func EnableApprove(chatID int64) error {
	enabled, err := IsApproveEnabled(chatID)
	if err != nil {
		return err
	}
	if enabled {
		return nil
	}

	_, err = disableColl.DeleteOne(tdCtx, bson.M{"chat_id": chatID})
	if err != nil {
		return fmt.Errorf("error while enabling chat: %w", err)
	}

	return nil
}
