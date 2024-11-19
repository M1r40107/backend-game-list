package main

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewJogo(nome, descricao string, preco, nota float64) Jogos {
	return Jogos{
		ID:        uuid.New().String(),
		Nome:      nome,
		Descricao: descricao,
		Preco:     preco,
		Nota:      nota,
	}
}

func NewMongoClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://guilherme:gui123456@jogos.dtv836u.mongodb.net/?retryWrites=true&w=majority&appName=jogos")
	mongoTestClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error while connecting:", err)
	}

	log.Println("Successfully connected.")

	err = mongoTestClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Ping failed:", err)
	}

	log.Println("Ping success.")

	return mongoTestClient
}

func InsertGame(coll *mongo.Collection, game *Jogos) (interface{}, error) {
	result, err := coll.InsertOne(context.Background(), game)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func FindGameByName(coll *mongo.Collection, gameID string) (*Jogos, error) {
	var jgs Jogos
	err := coll.FindOne(context.Background(), bson.D{{Key: "id", Value: gameID}}).Decode(&jgs)
	if err != nil {
		return nil, err
	}
	return &jgs, nil
}

func DeleteGameByName(coll *mongo.Collection, gameID string) (int64, error) {
	result, err := coll.DeleteOne(context.Background(), bson.D{{Key: "id", Value: gameID}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func TestMongoOperations(t *testing.T) {
	mongoTestClient := NewMongoClient()
	defer func() {
		if err := mongoTestClient.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	coll := mongoTestClient.Database("menu-jogos-avaliacao").Collection("game_test")

	t.Run("Insert Game", func(t *testing.T) {
		game := NewJogo("Bloodborne", "Action RPG", 39.99, 9.7)

		result, err := InsertGame(coll, &game)
		if err != nil {
			t.Fatal("Insert operation failed:", err)
		}

		t.Log("Insert successful, result ID:", result)

		insertedGame, err := FindGameByName(coll, game.ID)
		if err != nil {
			t.Fatal("Find operation failed:", err)
		}

		if insertedGame.Nome != game.Nome {
			t.Fatalf("Expected game name %s, got %s", game.Nome, insertedGame.Nome)
		}

		t.Log("Find operation successful, game:", insertedGame)

		deleteCount, err := DeleteGameByName(coll, game.ID)
		if err != nil {
			t.Fatal("Delete operation failed:", err)
		}

		if deleteCount == 0 {
			t.Fatal("No documents were deleted")
		}

		t.Log("Delete successful, deleted count:", deleteCount)
	})
}
