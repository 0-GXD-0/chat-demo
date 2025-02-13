package service

import (
	"chat/conf"
	"chat/model/ws"
	"context"
	"time"
)

func InsertMsg(database, id string, content string, read uint, expire int64) error {
	//插入到mongodb
	collection := conf.MongoDBClient.Database(database).Collection(id)
	comment := ws.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
		Read:      read,
	}
	_, err := collection.InsertOne(context.TODO(), comment)
	return err
}
