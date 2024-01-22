package log

import (
	"github.com/IBM/sarama"
	"github.com/XC-Zero/zero_common/config"
	"github.com/XC-Zero/zero_common/kafka"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"
	"time"
)

func CreateFileWriter(fileName string, filePath ...string) (writer io.Writer, err error) {
	paths := path.Join(filePath...)

	writer, err = os.Create(path.Join(paths, fileName))
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

type kafkaWriter struct {
	topic string
	sarama.SyncProducer
}

func (k *kafkaWriter) Write(p []byte) (n int, err error) {
	_, _, err = k.SendMessage(&sarama.ProducerMessage{
		Topic:     k.topic,
		Value:     sarama.ByteEncoder(p),
		Timestamp: time.Now(),
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return len(p), nil
}

func CreateKafkaWriter(config config.KafkaConfig, topic string) (writer io.Writer, err error) {
	pro, err := kafka.InitKafkaProducer(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &kafkaWriter{
		topic:        topic,
		SyncProducer: pro,
	}, nil
}
