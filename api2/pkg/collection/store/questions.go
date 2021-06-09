package store

import (
	"context"
	v1 "github.com/nrc-no/core-kafka/pkg/collection/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Questions struct {
	mongoClient *mongo.Client
}

type StoredQuestion struct {
	Key         string
	Country     string
	Formulation string
	AnswerType  string
}

func NewQuestionsStore(mongoClient *mongo.Client) *Questions {
	return &Questions{
		mongoClient: mongoClient,
	}
}

func (q *Questions) GetQuestionsForCountry(ctx context.Context, country string) (*v1.CountrySpecificBeneficiaryQuestions, error) {

	cursor, err := q.mongoClient.Database("core").Collection("questions").Find(ctx, bson.M{
		"country": country,
	})
	if err != nil {
		return nil, err
	}

	var questions []v1.Question

	for {
		if !cursor.Next(ctx) {
			break
		}
		var storedQuestion StoredQuestion
		if err := cursor.Decode(&storedQuestion); err != nil {
			return nil, err
		}
		questions = append(questions, v1.Question{
			Formulation: storedQuestion.Formulation,
			Key:         storedQuestion.Key,
			AnswerType:  storedQuestion.AnswerType,
		})
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &v1.CountrySpecificBeneficiaryQuestions{
		Country:   country,
		Questions: questions,
	}, nil

}

func (q *Questions) GetGlobalQuestions(ctx context.Context) (*v1.GlobalBeneficiaryQuestions, error) {

	cursor, err := q.mongoClient.Database("core").Collection("questions").Find(ctx, bson.M{
		"country": "",
	})
	if err != nil {
		return nil, err
	}

	var questions []v1.Question

	for {
		if !cursor.Next(ctx) {
			break
		}
		var storedQuestion StoredQuestion
		if err := cursor.Decode(&storedQuestion); err != nil {
			return nil, err
		}
		questions = append(questions, v1.Question{
			Formulation: storedQuestion.Formulation,
			Key:         storedQuestion.Key,
			AnswerType:  storedQuestion.AnswerType,
		})
	}
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}

	return &v1.GlobalBeneficiaryQuestions{
		Questions: questions,
	}, nil

}
