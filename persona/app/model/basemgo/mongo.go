package basemgo

import (
	// "log"

	"time"

	"github.com/lew/persona/app/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var AnalyseDB *DBInstal

// var ServerDB *DBInstal
var games map[float64]*DBInstal = map[float64]*DBInstal{}

type DBInstal struct {
	database *mgo.Database
}

// // mongo生成索引接口
// type DBInterface interface {
// 	makeIndex() []*mgo.Index
// 	// transToMongo(interface{}) interface{}
// }

func InitMongo() {
	// 初始化并保存mongo的session和database
	AnalyseDB = &DBInstal{
		initDBByStr(
			config.Cfg.MustValue("mongo.persona", "url", ""))}
	// 初始化游戏服的各个数据库
	gamesArr := config.Cfg.MustValueArray("mongos", "games", ",")
	for _, j := range gamesArr {
		tempDb := &DBInstal{
			initDB(
				config.Cfg.MustValueArray(j, "url", ","),
				config.Cfg.MustInt(j, "poolLimit", 4096),
				config.Cfg.MustValue(j, "db", "persona"),
				config.Cfg.MustValue(j, "username", ""),
				config.Cfg.MustValue(j, "password", ""),
				config.Cfg.MustValue(j, "replicaSet", ""))}

		if tempDb.database != nil {
			games[config.Cfg.MustFloat64(j, "appId", 200005)] = tempDb
		} else {
			sectionVal, _ := config.Cfg.GetSection(j)
			log.Fatal("mongo init Fatal ", j, sectionVal)
		}

	}
}

func initDB(addrs []string, poolLimit int, dbName string, un string, pw string, rs string) *mgo.Database {

	var err error

	// Dial格式 mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	// 数据库：  persona
	// 用户名：     person
	// 密码：         pNdxd8n&61gEx
	// mongodb://person:pNdxd8n&61gEx@dds-2ze5b4e0c4ef1c941.mongodb.rds.aliyuncs.com:3717,dds-2ze5b4e0c4ef1c942.mongodb.rds.aliyuncs.com:3717/persona?replicaSet=mgset-9000737

	dialInfo := &mgo.DialInfo{
		Addrs:          addrs,
		Direct:         false,
		Timeout:        time.Second * 1,
		PoolLimit:      poolLimit,
		Username:       un,
		Password:       pw,
		ReplicaSetName: rs,
		Database:       dbName,
	}
	//创建一个维护套接字池的session

	session, err := mgo.DialWithInfo(dialInfo)
	var database *mgo.Database = nil
	if err == nil {
		session.SetMode(mgo.Monotonic, true)
		//使用指定数据库
		// database = session.DB(dbName)
		database = session.DB(dbName)
	} else {
		log.Panicf("initDB err:%s \n dialInfo:%+v", err.Error(), dialInfo)
	}

	return database
}
func initDBByStr(dbStr string) *mgo.Database {

	var err error

	// Dial格式 mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	// 数据库：  persona
	// 用户名：     person
	// 密码：         pNdxd8n&61gEx
	// mongodb://person:pNdxd8n&61gEx@dds-2ze5b4e0c4ef1c941.mongodb.rds.aliyuncs.com:3717,dds-2ze5b4e0c4ef1c942.mongodb.rds.aliyuncs.com:3717/persona?replicaSet=mgset-9000737;maxPoolSize=4096
	session, err := mgo.Dial(dbStr)
	var database *mgo.Database = nil
	if err == nil {
		session.SetMode(mgo.Monotonic, true)
		//使用指定数据库
		// database = session.DB(dbName)
		database = session.DB("")

	} else {
		log.Panicf("initDB err:%s \n dbStr:%s", err.Error(), dbStr)
	}

	return database
}

// var session *mgo.Session
// func GetSession() *mgo.Session {
// 	return session
// }

func GetDBCon(appId float64, collName string) *mgo.Collection {
	return games[appId].database.C(collName)
}
func HasDBCon(appId float64) bool {
	if len(games) == 0 || games[appId] == nil {
		return false
	} else {
		return true
	}
}

// var database *mgo.Database
// func GetDataBase() *mgo.Database {
// 	return database
// }

func GetErrNotFound() error {
	return mgo.ErrNotFound
}
func GetIndex(keys []string, unique bool, dropDups bool, background bool) *mgo.Index {
	index := mgo.Index{
		Key:        keys,       // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     unique,     // 唯一索引 同mysql唯一索引
		DropDups:   dropDups,   // 索引重复替换旧文档,Unique为true时失效
		Background: background, // 后台创建索引
	}
	return &index
}

func (dbInstal *DBInstal) GetCon(cName string, makeIndex func() []*mgo.Index) *mgo.Collection {
	con := dbInstal.database.C(cName)
	if makeIndex != nil {
		haveIndex, _ := con.Indexes()
		indexs := makeIndex()
		// 索引数不匹配,则重建索引
		if len(haveIndex) != (len(indexs) + 1) {
			for index := 0; index < len(indexs); index++ {
				con.EnsureIndex(*indexs[index])
			}
		}
	}
	return con
}

// func GetCon(cName string, dbT DBInterface) *mgo.Collection {
// 	con := GetDataBase().C(cName)
// 	if dbT != nil {
// 		haveIndex, _ := con.Indexes()
// 		indexs := dbT.makeIndex()
// 		// 索引数不匹配,则重建索引
// 		if len(haveIndex) != (len(indexs) + 1) {
// 			for index := 0; index < len(indexs); index++ {
// 				con.EnsureIndex(*indexs[index])
// 			}
// 		}
// 	}
// 	return con
// }

/**
import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	RedisURL            = "redis://127.0.0.1:6379/2"
	redisMaxIdle        = 3   //最大空闲连接数
	redisIdleTimeoutSec = 240 //最大空闲连接时间
	RedisPassword       = ""
)

func logTtest() {
	log.Println("fuck")
}

type Person struct {
	Name  string
	Phone string
}

func dbCon() {
	session, err := mgo.Dial("")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		panic(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println("Phone:", result.Phone)
}
*/
