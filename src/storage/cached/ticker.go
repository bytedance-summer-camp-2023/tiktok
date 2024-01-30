package cached

import (
	redis2 "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"tiktok/src/storage/redis"
	"tiktok/src/utils/logging"
	"time"
)

type TimeTicker struct {
	Tick *time.Ticker
	Work func(client redis2.UniversalClient) error
}

func (t *TimeTicker) Start() {
	for range t.Tick.C {
		err := t.Work(redis.Client)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error happens when dealing with cron job")
			continue
		}
	}
}

func NewTick(interval time.Duration, f func(client redis2.UniversalClient) error) *TimeTicker {
	return &TimeTicker{
		Tick: time.NewTicker(interval),
		Work: f,
	}
}