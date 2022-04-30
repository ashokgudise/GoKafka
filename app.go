// Example function-based Apache Kafka producer
package main

/**
 * Copyright 2016 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"strconv"
	"time"
)

const (
	topicName     = "topic_for_go_lang"
	brokerAddress = "localhost:9092"
)

func main() {

	produce()
	consume()

}

func produce() {

	i := 0

	// Produce messages to topic (asynchronously)
	topic := topicName

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokerAddress})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	for {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            []byte(strconv.Itoa(i)),
			Value:          []byte("This is message" + strconv.Itoa(i)),
		}, nil)

		i++

		time.Sleep(time.Second)
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}

func consume() {

	consumer, error := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "golang.group",
		"auto.offset.reset": "earliest",
	})

	if error != nil {
		panic(error)
	}

	consumer.SubscribeTopics([]string{topicName}, nil)

	for {
		message, error := consumer.ReadMessage(-1)
		if error == nil {
			fmt.Printf("Message on %s: %s\n", message.TopicPartition, string(message.Value))
		} else {
			fmt.Printf(" Consumer error: %v (%v)\n", error, message)
		}

	}

	consumer.Close()
}
