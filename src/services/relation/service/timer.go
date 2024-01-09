package service

import (
	"context"
	"encoding/json"
	"tiktok/storage/redis"
)

const frequency = 10

func consume() error {
	msgs, err := RelationMq.ConsumeSimple()
	if err != nil {
		logger.Errorf("RelationMQ Err: %s", err.Error())
		return err
	}
	// 将消息队列的消息全部取出
	for msg := range msgs {
		rc := new(redis.RelationCache)
		// 解析json
		if err = json.Unmarshal(msg.Body, &rc); err != nil {
			logger.Errorf("RelationMQ json unmarshal Err: %s", err.Error())
			continue
		}
		// 将结构体存入redis
		if err = redis.UpdateRelation(context.Background(), rc); err != nil {
			logger.Errorf("RelationMQ Err: %s", err.Error())
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
