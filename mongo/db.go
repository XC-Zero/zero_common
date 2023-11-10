package mongo

import (
	"context"
	"github.com/qiniu/qmgo/options"
	"gitlab.tessan.com/data-center/tessan-erp-common/config"
	"go.mongodb.org/mongo-driver/event"
	op "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)
import "github.com/qiniu/qmgo"

var timeout = time.Minute * 15

func InitMongoClient(config config.MongoDBConfig) (*qmgo.Client, error) {
	clientOptions := options.ClientOptions{ClientOptions: &op.ClientOptions{
		//Timeout:         &timeout,
		//MaxConnIdleTime: &timeout,
		Monitor: &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				{
					length := len(startedEvent.Command)
					if length > 500 {
						length = 500
					}
					log.Println("[INFO]" + string(startedEvent.Command[:length]))
				}

			},
			//Succeeded: func(ctx context.Context, startedEvent *event.CommandSucceededEvent) {
			//	log.Println("[INFO]" + time.Duration(startedEvent.DurationNanos).String())
			//
			//},
			Failed: func(ctx context.Context, startedEvent *event.CommandFailedEvent) {
				log.Println("[ERROR]" + startedEvent.Failure)

			},
		},
	}}
	// 连接到MongoDB
	client, err := qmgo.NewClient(context.TODO(), &qmgo.Config{
		Uri:      config.URL,
		Database: config.DBName,
	}, clientOptions)

	if err != nil {
		return nil, err
	}

	return client, nil
}
