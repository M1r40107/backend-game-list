package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JogosRepo struct {
	MongoCollection *mongo.Collection
}

func (r *JogosRepo) InsertGame(jgs *Jogos) (interface{}, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), jgs)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func (r *JogosRepo) FindGameByName(gameID string) (*Jogos, error) {
	var jgs Jogos
	err := r.MongoCollection.FindOne(context.Background(), bson.D{{Key: "ID", Value: gameID}}).Decode(&jgs)
	if err != nil {
		return nil, err
	}
	return &jgs, nil
}

func (r *JogosRepo) FindAllGames() ([]Jogos, error) {
	cursor, err := r.MongoCollection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var jogos []Jogos
	err = cursor.All(context.Background(), &jogos)
	if err != nil {
		return nil, fmt.Errorf("result decode error: %s", err.Error())
	}
	return jogos, nil
}

func (r *JogosRepo) UpdateGameByName(gameID string, updateGm *Jogos) (int64, error) {
	result, err := r.MongoCollection.UpdateOne(
		context.Background(),
		bson.D{{Key: "ID", Value: gameID}},
		bson.D{{Key: "$set", Value: updateGm}},
	)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (r *JogosRepo) DeleteGameByName(gameID string) (int64, error) {
	result, err := r.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "ID", Value: gameID}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (r *JogosRepo) DeleteAllGames() (int64, error) {
	result, err := r.MongoCollection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
