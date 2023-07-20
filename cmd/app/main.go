package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/alonsofritz/simple-api-kafka-go/internal/infra/akafka"
	"github.com/alonsofritz/simple-api-kafka-go/internal/infra/repository"
	"github.com/alonsofritz/simple-api-kafka-go/internal/infra/web"
	"github.com/alonsofritz/simple-api-kafka-go/internal/usecase"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// criar conexao com banco
	db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/products")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// criar meu objeto do meu repository
	repository := repository.NewProductRepositoryMysql(db)
	createProductUseCase := usecase.NewCreateProductUseCase(repository)
	listProductsUseCase := usecase.NewListProductsUseCase(repository)

	productHandlers := web.NewProductHandlers(createProductUseCase, listProductsUseCase)

	// Me ajuda a mapear minhas rotas quando criar o servidor web
	router := chi.NewRouter()
	router.Post("/products", productHandlers.CreateProductHandler)
	router.Get("/products", productHandlers.ListProductsHandler)

	// Concorda se eu deixar meu server rodando aqui, ele vai travar
	// e nao vai dar continuidade ao codigo, dessa forma impedindo
	// a implementação do kafka, dessa forma, eu uso novamente
	// go routines. Dessa forma ele começa a servir minha aplicação em
	// outra rotina
	// Sobe Servidor WEB
	go http.ListenAndServe(":8000", router)

	// cria canal de comunicação para ajudar a trabalhar para se comunicar em trheads diferentes
	msgChan := make(chan *kafka.Message)

	// Utilizar minha funcao do kafka para comecar a consumir as msgs
	// se comunica com o kafka, recebendo os dados do kafka, e jogando para o canal msgChan
	// para que eu consiga seguir com a execução. mesmo que a execucao anterior seja "infinita", basta jogarmos para outra Thread
	// GO Magic: go routines
	go akafka.Consume([]string{"products"}, "host.docker.internal:9094", msgChan)

	// Perceba que a parte a seguir nao tem conhecimento nenhum da
	// area de dominio/ do domonio da minha aplicação
	// ele nao sabe se tenho entidades, utilizamos apenas o DTO,
	// que ele sabe que vai "mandar o dto para o usecases" e vai receber
	// uma resposta

	// Fica lendo as mensagem que recebemos do Kafka
	// os dados que recebemos do kafka vamos "hidratar" o DTO
	for msg := range msgChan {
		dto := usecase.CreateProductInputDTO{}

		// Pega o json e converte para Struct (similar ao GSON faz com Classes)
		// pega os dados que recebeu no Kafka e joga no DTO
		err := json.Unmarshal(msg.Value, &dto)
		if err != nil {
			// log erro
			continue
		}

		_, err = createProductUseCase.Execute(dto)
	}

}
