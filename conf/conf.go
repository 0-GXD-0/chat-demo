package conf

import (
	"chat/model"
	"fmt"
	"strings"

	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"gopkg.in/ini.v1"
)

var (
	MongoDBClient *mongo.Client

	AppMode  string
	HttpPort string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	MongoDBName string
	MongoDBAddr string
	MongoDBPwd  string
	MongoDBPort string
)

func Init() {
	//从本地读取环境
	file, err := ini.Load("conf/config.ini")
	if err != nil {
		fmt.Println("load config.ini failed", err)
	}
	LoadServer(file)
	LoadMySQL(file)
	LoadMongoDB(file)
	logging.Info("运行到MongoDB前一步")
	MongoDB()
	logging.Info("运行到MongoDB后一步")

	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	model.Database(path)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMySQL(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadMongoDB(file *ini.File) {
	logging.Info("Loading MongoDB configuration")
	logging.Info("MongoDBName from config file:", file.Section("MongoDB").Key("MongoDBName").String())
	logging.Info("MongoDBAddr from config file:", file.Section("MongoDB").Key("MongoDBAddr").String())
	logging.Info("MongoDBPwd from config file:", file.Section("MongoDB").Key("MongoDBPwd").String())
	logging.Info("MongoDBPort from config file:", file.Section("MongoDB").Key("MongoDBPort").String())

	MongoDBName = file.Section("MongoDB").Key("MongoDBName").String()
	logging.Info("MongoDBName:", MongoDBName)
	MongoDBAddr = file.Section("MongoDB").Key("MongoDBAddr").String()
	MongoDBPwd = file.Section("MongoDB").Key("MongoDBPwd").String()
	MongoDBPort = file.Section("MongoDB").Key("MongoDBPort").String()

}

func MongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoDBClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info("mongodb connect failed", err)
		panic(err)
	}
	logging.Info("mongodb connect success")
}
