package myconnect

import (
	"context"
	"crypto-rate/libs/myfunc"

	kafka "github.com/segmentio/kafka-go"
)

func KafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func KafkaReader(kafkaURL string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaURL},
		Partition: 0,
		Topic:     topic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
}

func KafkaProducer(kafkaURL string, topic string, key string, value string) error {
	// to produce messages
	writer := KafkaWriter(kafkaURL, topic)
	defer writer.Close()

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	err := writer.WriteMessages(context.TODO(), msg)
	if err != nil {
		myfunc.MyErrFormat(err)
	}
	return nil
}
