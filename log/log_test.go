package log

import (
	"github.com/XC-Zero/zero_common/config"
	"testing"
)

func TestLog_Send(t *testing.T) {
	writer, err := CreateKafkaWriter(config.KafkaConfig{
		Address: []string{"192.168.185.42:39092"},
	}, "know_weather_xml_log")
	if err != nil {
		panic(err)
	}
	log := New(LogOptions{
		ChannelSize: 100,
		FormatJson:  true,
		Writer:      writer,
	})
	log.Send("???", nil, INFO)
}
