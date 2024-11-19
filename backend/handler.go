package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var gamesRepo *JogosRepo

func addGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var game Jogos
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&game)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = gamesRepo.InsertGame(&game)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(game)
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func listGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := gamesRepo.FindAllGames()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(games)
}

func deleteGamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		deletedCount, err := gamesRepo.DeleteAllGames()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int64{"deletedCount": deletedCount})
	} else {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func connectToMongo() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb+srv://guilherme:gui123456@jogos.dtv836u.mongodb.net/?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Erro ao conectar:", err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Falha no ping:", err)
	}

	log.Println("Conectado ao MongoDB")
	return client
}
