package consumer

import (
	"context"
	"github.com/cenkalti/backoff"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/superwhys/superGo/superLog"
	"github.com/ugorji/go/codec"
	"io"
	"time"
)

type OptionKafkaReaderFunc func(writer *SuperKafkaReader)
type SuperKafkaReader struct {
	SuperReaderConfig *kafka.ReaderConfig
	SuperReader       *kafka.Reader
}

const LocalKafkaIps = "localhost:9092"

func InitReader(kafkaIps, topic, groupId string, opts ...OptionKafkaReaderFunc) *SuperKafkaReader {
	superReader := &SuperKafkaReader{
		SuperReaderConfig: &kafka.ReaderConfig{
			Brokers: []string{kafkaIps},
			GroupID: groupId,
			Topic:   topic,
			Dialer:  nil,
		},
	}
	for _, opt := range opts {
		opt(superReader)
	}
	superReader.SuperReader = kafka.NewReader(*superReader.SuperReaderConfig)
	return superReader
}

func AddDialer(timeOut time.Duration) OptionKafkaReaderFunc {
	return func(skw *SuperKafkaReader) {
		skw.SuperReaderConfig.Dialer = &kafka.Dialer{
			ClientID:        "",
			Timeout:         timeOut,
			Deadline:        time.Time{},
			LocalAddr:       nil,
			DualStack:       false,
			FallbackDelay:   0,
			KeepAlive:       0,
			Resolver:        nil,
			TLS:             nil,
			SASLMechanism:   nil,
			TransactionalID: "",
		}
	}
}

func AddMaxAttempts(maxAttempts int) OptionKafkaReaderFunc {
	return func(skw *SuperKafkaReader) {
		skw.SuperReaderConfig.MaxAttempts = maxAttempts
	}
}

func AddMaxPartition(partition int) OptionKafkaReaderFunc {
	return func(skw *SuperKafkaReader) {
		skw.SuperReaderConfig.Partition = partition
	}
}

func (sr *SuperKafkaReader) ReadMessageWithProto(protoMsg proto.Message) (*kafka.Message, error) {
	bf := backoff.NewExponentialBackOff()
	for {
		ctx := context.Background()
		msg, err := sr.SuperReader.ReadMessage(ctx)
		if err == io.EOF {
			return nil, errors.New("Reader closed")
		}
		if err != nil {
			superLog.Error("Read kafka", err)
			time.Sleep(bf.NextBackOff())
			continue
		}
		if err := proto.Unmarshal(msg.Value, protoMsg); err != nil {
			superLog.Errorf("Decode kafka message. Partition=%d Offset=%d %s", msg.Partition, msg.Offset, err)
			continue
		}
		msg.Value = nil
		return &msg, nil
	}
}

func (sr *SuperKafkaReader) ReadMessageWithJson(out interface{}) (*kafka.Message, error) {
	bf := backoff.NewExponentialBackOff()
	for {
		ctx := context.Background()
		msg, err := sr.SuperReader.ReadMessage(ctx)
		if err == io.EOF {
			return nil, errors.New("Reader closed")
		}
		if err != nil {
			superLog.Error("Read kafka", err)
			time.Sleep(bf.NextBackOff())
			continue
		}
		if err := codec.NewDecoderBytes(msg.Value, &codec.JsonHandle{}).Decode(out); err != nil {
			superLog.Errorf("Decode kafka message. Partition=%d Offset=%d %s", msg.Partition, msg.Offset, err)
			continue
		}
		msg.Value = nil
		return &msg, nil
	}
}
