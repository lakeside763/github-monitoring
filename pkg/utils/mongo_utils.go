package utils

import "go.mongodb.org/mongo-driver/mongo"

func HandleMongoError(err error) error {
	if err == mongo.ErrNoDocuments {
		return nil
	}
	return err
}