package iam

import "go.mongodb.org/mongo-driver/mongo"

type AttributeStore struct {
	collection *mongo.Collection
}
