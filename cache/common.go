package cache

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func init() {
	//加载配置信息文件
	file, err := ini.Load("conf/config.ini")
	if err != nil {
		fmt.Println("load config.ini failed", err)
	}
	//读取配置信息文件内容
	LoadRedis(file)
	Redis()
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64) //string to uint64
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})
	_, err := client.Ping().Result()
	if err != nil {
		logging.Info("redis connect failed", err)
		panic(err)
	}
	RedisClient = client
	logging.Info("redis connect success")
}
