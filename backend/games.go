package main

type Jogos struct {
	ID        string  `json: "id, omitempty" bson:"id"`
	Nome      string  `json:"nome, omitempty" bson:"nome"`
	Descricao string  `json:"descricao, omitempty" bson:"descricao"`
	Preco     float64 `json:"preco, omitempty" bson:"preco"`
	Nota      float64 `json:"nota, omitempty" bson:"nota"`
}
