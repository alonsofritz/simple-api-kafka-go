package akafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Ler e Receber dados do Apache Kafka
func Consume(topics []string, servers string, msgChan chan *kafka.Message) {
	kafkaConsumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		// servidor de comunicacao com kafka
		"bootstrap.servers": servers,
		// Definir Grupo de consumidores caso tenha varias instancias
		"group.id": "simple-api-kafka-go",
		// Pergar mensagens sempre do inicio
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}

	// Falar que preciso consumir os topicos
	kafkaConsumer.SubscribeTopics(topics, nil)
	// Vamos rodar um looping infinito: ficar pegando as msgs do kafka e ficar jogando elas no canal de mensagens aqui
	for {
		// ReadMessage() <- recebe um timeout, caso negativo, faz o tempo todo
		//fica vendo o tempo todo msg, toda vez que receber msg do kafka, ele vai manda para o Canal msgChan
		// dai pegaremos esse canal em outra thread e faremos o processamento dessas mensagens
		msg, err := kafkaConsumer.ReadMessage(-1)
		if err == nil {
			// consumo das msgs
			msgChan <- msg
		}
	}
}
