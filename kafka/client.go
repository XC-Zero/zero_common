package kafka

import (
	"github.com/IBM/sarama"
	"github.com/XC-Zero/zero_common/config"
	"github.com/pkg/errors"
	"log"
	"time"
)

func InitKafkaClient(kafka config.KafkaConfig) (sarama.Client, error) {
	conf := sarama.NewConfig()
	conf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	conf.ChannelBufferSize = 100 * 100 * 10
	conf.Consumer.MaxProcessingTime = time.Millisecond * 10
	conf.Consumer.Fetch.Min = 1 * (1 << 27)
	conf.Consumer.Fetch.Default = 1 * (1 << 28)
	conf.Consumer.Fetch.Max = 1 * (1 << 29)
	//// 批量同步的最大批次时间是五分钟 所以这里设置6分钟
	conf.Consumer.Group.Session.Timeout = 30 * time.Second

	//conf.Consumer.MaxPollRecords = 100
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(kafka.Address, conf)
	if err != nil {
		return nil, err
	}
	return client, err
}

func InitKafkaProducer(kafka config.KafkaConfig) (sarama.SyncProducer, error) {

	conf := sarama.NewConfig()

	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(kafka.Address, conf)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	log.Println("Create producer!")

	var s = sarama.ProducerMessage{
		Topic:     "TEST_TOPIC",
		Key:       sarama.StringEncoder("TEST"),
		Value:     sarama.StringEncoder("INIT KAFKA PRODUCER!"),
		Timestamp: time.Now(),
	}
	p, n, err := producer.SendMessage(&s)

	if err != nil {
		panic(err)
	}
	log.Println(p, n)

	return producer, err
}
