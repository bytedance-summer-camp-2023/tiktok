package service

import (
	"context"
	"encoding/json"
	"github.com/bytedance-summer-camp-2023/tiktok/dal/redis"
	"github.com/bytedance-summer-camp-2023/tiktok/pkg/gocron"
)

const frequency = 10

// 点赞服务消息队列消费者
func consume() error {
	msgs, err := RelationMq.ConsumeSimple()
	if err != nil {
		logger.Errorf("RelationMQ Err: %s", err.Error())
		return err
	}
	// 将消息队列的消息全部取出
	for msg := range msgs {
		//err := redis.LockByMutex(context.Background(), redis.RelationMutex)
		//if err != nil {
		//	logger.Errorf("Redis mutex lock error: %s", err.Error())
		//	return err
		//}
		rc := new(redis.RelationCache)
		// 解析json
		if err = json.Unmarshal(msg.Body, &rc); err != nil {
			logger.Errorf("RelationMQ json unmarshal Err: %s", err.Error())
			//err = redis.UnlockByMutex(context.Background(), redis.RelationMutex)
			//if err != nil {
			//	logger.Errorf("Redis mutex unlock error: %s", err.Error())
			//	return err
			//}
			continue
		}
		//logger.Errorf("RelationMQ json unmarshal Err: %s", err.Error())
		// 将结构体存入redis
		if err = redis.UpdateRelation(context.Background(), rc); err != nil {
			logger.Errorf("RelationMQ Err: %s", err.Error())
			//err = redis.UnlockByMutex(context.Background(), redis.RelationMutex)
			//if err != nil {
			//	logger.Errorf("Redis mutex unlock error: %s", err.Error())
			//	return err
			//}
			continue
		}
		//err = redis.UnlockByMutex(context.Background(), redis.RelationMutex)
		//if err != nil {
		//	logger.Errorf("Redis mutex unlock error: %s", err.Error())
		//	return err
		//}
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

// gocron定时任务,每隔一段时间就让Consumer消费消息队列的所有消息
func GoCron() {
	s := gocron.NewSchedule()
	s.Every(frequency).Tag("relationMQ").Seconds().Do(consume)
	s.StartAsync()
}
