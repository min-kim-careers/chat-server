package chat

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

var streamCtx = context.Background()

type Stream struct {
	writer *kafka.Writer
}

func NewStream() *Stream {
	return &Stream{
		writer: initWriter(),
	}
}

func initWriter() *kafka.Writer {
	brokerAddress := os.Getenv("KAFKA_BROKER_ADDRESS")
	topicName := os.Getenv("KAFKA_TOPIC_NAME")

	if brokerAddress == "" || topicName == "" {
		log.Fatal("Missing required Kafka configuration")
	}

	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topicName,
	})
}

func generateStreamKey(roomID string) string {
	return "chat:room:" + roomID
}

func (stream *Stream) Write(roomID string, msg *Message) error {
	key := generateStreamKey(roomID)
	value := SerializeMessage(msg)

	err := stream.writer.WriteMessages(streamCtx,
		kafka.Message{
			Key:   []byte(key),
			Value: value,
		})
	if err != nil {
		log.Printf("Error writing to stream for room <%s>: %v", key, err)
		return err
	}

	return nil
}
