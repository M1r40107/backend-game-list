package main

import (
	"context"
	"log"
	"net/http"
)

func main() {

	client := connectToMongo()
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Erro ao desconectar do MongoDB: %v", err)
		}
	}()

	db := client.Database("menudb")
	gamesRepo = &JogosRepo{MongoCollection: db.Collection("gamesdb")}

	http.HandleFunc("/addgame", addGameHandler)
	http.HandleFunc("/games", listGamesHandler)
	http.HandleFunc("/deletegames", deleteGamesHandler)

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
