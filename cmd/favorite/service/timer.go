package service

import (
	"context"
	"encoding/json"
	"fmt"
	"tiktok/dal/redis"
)

const frequency = 10

// 点赞服务消息队列消费者
func consume() error {
	msgs, err := FavoriteMq.ConsumeSimple()
	if err != nil {
		fmt.Println(err.Error())
		logger.Errorf("FavoriteMQ Err: %s", err.Error())
	}

	// 将消息队列的消息全部取出
	for msg := range msgs {

		fc := new(redis.FavoriteCache)
		// 解析json
		if err = json.Unmarshal(msg.Body, &fc); err != nil {
			logger.Errorf("json unmarshal error: %s", err.Error())
			fmt.Println("json unmarshal error:" + err.Error())

			continue
		}
		fmt.Printf("==> Get new message: %v\n", fc)
		// 将结构体存入redis
		if err = redis.UpdateFavorite(context.Background(), fc); err != nil {
			logger.Errorf("json unmarshal error: %s", err.Error())
			fmt.Println("json unmarshal error:" + err.Error())
			continue
		}

		if !autoAck {
			err := msg.Ack(true)
			if err != nil {
				logger.Errorf("ack error: %s", err.Error())
				return err
			}
		}
	}
	return nil
}
