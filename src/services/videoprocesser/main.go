package main

import (
	"github.com/go-kit/kit/transport/amqp"
	"github.com/sirupsen/logrus"
	"sync"
	"tiktok/src/constant/config"
	"tiktok/src/extra/tracing"
	"tiktok/src/models"
	"tiktok/src/utils/logging"
)

func main() {
	conn, err := amqp.Dial(rabbitmq.BuildMQConnAddr())
	if err != nil {
		panic(err)
	}

	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	tp, err := tracing.SetTraceProvider(config.VideoPicker)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer func(ch *amqp.Channel) {
		err := ch.Close()
		if err != nil {
			panic(err)
		}
	}(ch)

	if _, err = ch.QueueDeclare(strings.VideoPicker, true, false, false, false, nil); err != nil {
		panic(err)
	}

	if err = ch.Qos(1, 0, false); err != nil {
		panic(err)
	}

	go Consume(ch)
	logger := logging.LogService("VideoPicker")
	logger.Infof(strings.VideoPicker + "is running now")

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func Consume(channel *amqp.Channel) {
	msg, err := channel.Consume(strings.VideoPicker, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case d := <-msg:
			//解包 Otel Context
			ctx := rabbitmq.ExtractAMQPHeaders(context.Background(), d.Headers)
			ctx, span := tracing.Tracer.Start(ctx, "VideoPickerService")
			logger := logging.LogService("VideoPicker.Picker").WithContext(ctx)

			var raw models.RawVideo
			if err := json.Unmarshal(d.Body, &raw); err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when unmarshaling the prepare json body.")
			}

			span.End()
			err = d.Ack(false)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when dealing with the video...")
			}
		}
	}
}
