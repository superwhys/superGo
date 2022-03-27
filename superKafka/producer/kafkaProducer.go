package producer

import (
	"bytes"
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/superwhys/superGo/superLog"
	"github.com/ugorji/go/codec"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type OptionKafkaWriterFunc func(writer *SuperKafkaWriter)
type SuperKafkaWriter struct {
	SuperWriterConfig *kafka.WriterConfig
	SuperWriter       *kafka.Writer
}

const LocalKafkaIps = "localhost:9092"

func InitWriter(kafkaIps, topic string, opts ...OptionKafkaWriterFunc) *SuperKafkaWriter {
	superWriter := &SuperKafkaWriter{
		SuperWriterConfig: &kafka.WriterConfig{
			Brokers:  []string{kafkaIps},
			Topic:    topic,
			Dialer:   nil,
			Balancer: &kafka.LeastBytes{},
		},
	}
	for _, opt := range opts {
		opt(superWriter)
	}
	superWriter.SuperWriter = kafka.NewWriter(*superWriter.SuperWriterConfig)
	return superWriter
}

func AddDialer(timeOut time.Duration) OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.Dialer = &kafka.Dialer{
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

func AddWriteTimeout(timeOut time.Duration) OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.WriteTimeout = timeOut
	}
}

func AddBalancerHash() OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.Balancer = &kafka.Hash{}
	}
}

func AddBalancerCRC32() OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.Balancer = &kafka.CRC32Balancer{}
	}
}

func AddMaxAttempts(maxAttempts int) OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.MaxAttempts = maxAttempts
	}
}

func AddMaxBatchSize(batchSize int) OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.BatchSize = batchSize
	}
}

func AddMaxBatchBytes(batchBytes int) OptionKafkaWriterFunc {
	return func(skw *SuperKafkaWriter) {
		skw.SuperWriterConfig.BatchBytes = batchBytes
	}
}

func (sw *SuperKafkaWriter) WriteMessageWithJSON(key string, msg interface{}) error {
	out := bytes.Buffer{}
	enc := codec.NewEncoder(&out, &codec.JsonHandle{})
	if err := enc.Encode(msg); err != nil {
		return errors.Wrap(err, "Encode message to JSON")
	}
	return writeMessage(sw.SuperWriter, []byte(key), out.Bytes())
}

func (sw *SuperKafkaWriter) WriteMessageWithProto(key string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "Encode message")
	}
	return writeMessage(sw.SuperWriter, []byte(key), data)
}

func writeMessage(writer *kafka.Writer, key, value []byte) error {
	if len(key) == 0 {
		key = []byte(strconv.Itoa(rand.Int()))
	}
	cnt := 0
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err := writer.WriteMessages(ctx, kafka.Message{
			Key:   key,
			Value: value,
		})
		if err == nil {
			break
		}
		superLog.Errorf("Write kafka message, key=%s err=%s", string(key), err.Error())
		if !strings.Contains(err.Error(), "broken pipe") {
			return err
		}
		time.Sleep(time.Second)
		cnt++
		superLog.Infof("Retrying to write kafka, key=%s attempts=%d", string(key), cnt)
	}
	return nil
}
