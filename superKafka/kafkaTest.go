package superKafka

import (
	"fmt"
	"github.com/superwhys/superGo/superKafka/consumer"
	"github.com/superwhys/superGo/superKafka/producer"
)

type person struct {
	Name string `json:"name,omitempty"`
}

func main() {
	writer := producer.InitWriter(producer.LocalKafkaIps, "kafka_test")
	err := writer.WriteMessageWithJSON("", &person{
		Name: "why",
	})
	if err != nil {
		return
	}

	reader := consumer.InitReader(consumer.LocalKafkaIps, "kafka_test", "kafka_test_group")
	message, err := reader.ReadMessageWithJson(&person{})
	if err != nil {
		return
	}

	fmt.Println(message.Value)

}
