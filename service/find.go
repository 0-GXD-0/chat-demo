package service

import (
	"chat/conf"
	"chat/model/ws"
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

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

func FindMany(database, sendID, id string, time int64, pageSize int) (results []ws.Result, err error) {
	var resultMe []ws.Trainer  //id
	var resultYou []ws.Trainer //sendID
	sendIDCollection := conf.MongoDBClient.Database(database).Collection(sendID)
	idCollection := conf.MongoDBClient.Database(database).Collection(id)

	// 使用 time 参数进行过滤
	filter := bson.M{"startTime": bson.M{"$lt": time}}
	log.Printf("查询条件: %v", filter)

	sendIDTimeCurcor, err := sendIDCollection.Find(context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		log.Printf("查询 sendIDCollection 失败: %v", err)
		return nil, err
	}

	defer sendIDTimeCurcor.Close(context.TODO())

	idTimeCurcor, err := idCollection.Find(context.TODO(),
		filter,
		options.Find().SetSort(bson.D{{"startTime", -1}}),
		options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		log.Printf("查询 idCollection 失败: %v", err)
		return nil, err
	}

	defer idTimeCurcor.Close(context.TODO())

	err = sendIDTimeCurcor.All(context.TODO(), &resultYou)
	if err != nil {
		log.Printf("解析 sendIDTimeCursor 失败: %v", err)
		return nil, err
	}
	err = idTimeCurcor.All(context.TODO(), &resultMe)
	if err != nil {
		log.Printf("解析 idTimeCursor 失败: %v", err)
		return nil, err
	}
	results, err = AppendAndSort(resultMe, resultYou)
	if err != nil {
		log.Printf("合并和排序结果失败: %v", err)
		return nil, err
	}
	log.Printf("查询到的历史消息数量: %d", len(results))
	return results, err
}

func AppendAndSort(resultMe []ws.Trainer, resultYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultMe {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}

	for _, r := range resultYou {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	// 按 StartTime 排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].StartTime > results[j].StartTime
	})
	log.Printf("合并后的历史消息数量: %d", len(results))
	return results, err
}
