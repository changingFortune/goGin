package model

import (
	"encoding/json"
	"time"

	// "time"
	// "encoding/json"

	"github.com/lew/persona/app/model/basemgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 静态信息:
// 名字,性别,地区,头像
type Persona struct {
	Gid          float64            `json:"gid" bson:"gid"`
	AppId        float64            `json:"appId" bson:"appId"`
	UserId       float64            `json:"userId" bson:"userId"`
	FirstLogTime string             `json:"firstLogTime" bson:"firstLogTimed"` //第一次登录时间
	LogoutTime   string             `json:"logoutTime" bson:"logoutTime"`      //登出时间
	Duration     float64            `json:"duration" bson:"duration"`          //游戏总时长	2
	PlayDuration float64            `json:"playDuration" bson:"playDuration"`  //打牌时长	3
	GameCount    float64            `json:"gameCount" bson:"gameCount"`        //游戏把数	3
	WinCount     float64            `json:"winCount" bson:"winCount"`          //输赢次数
	WinScore     float64            `json:"winScore" bson:"winScore"`          //输赢分数
	Friends      map[string]float64 `json:"friends" bson:"friends"`            //好友列表	2
	DayStamp     string             `json:"dayStamp" bson:"dayStamp"`          //日期戳

}

type PersonaDB struct { // mongo对象
	Persona       `bson:",inline"`
	ID            bson.ObjectId `json:"_id" bson:"_id"`
	GameCountS    float64       `json:"gameCountS" bson:"gameCountS"`       //游戏把数 评分
	PlayDurationS float64       `json:"playDurationS" bson:"playDurationS"` //打牌时长 评分
	DurationS     float64       `json:"durationS" bson:"durationS"`         //游戏总时长 评分
	FriendsS      float64       `json:"friendsS" bson:"friendsS"`           //好友列表 评分
	TotalS        float64       `json:"totalS" bson:"totalS"`               //总分
}

// 生成表索引
func (temp *PersonaDB) makeIndex() []*mgo.Index {
	var indexs = []*mgo.Index{
		basemgo.GetIndex([]string{"gid", "dayStamp"}, true, true, true),
		basemgo.GetIndex([]string{"dayStamp", "totalS"}, false, true, true),
		basemgo.GetIndex([]string{"gid"}, false, true, true),
		basemgo.GetIndex([]string{"appId"}, false, true, true),
		basemgo.GetIndex([]string{"userId"}, false, true, true),
	}
	return indexs
}

// var _con *mgo.Collection

// func getCon() *mgo.Collection {
// 	if _con == nil {
// 		_con = basemgo.AnalyseDB.GetCon("persona", PersonaDB{}.makeIndex)
// 	}
// 	return _con
// }

// func (persona Persona) transToMongo() *PersonaDB {

// }

func AddOrUpdatePersona(p *PersonaDB) (doneUser string, err error) {

	con := basemgo.AnalyseDB.GetCon("persona", p.makeIndex)

	//可以添加一个或多个文档
	/* 对应mongo命令行
	   db.user.insert({username:"13888888888",summary:"code",
	   age:20,phone:"13888888888"})*/
	personaDB, currErr := FindPersona(p.Gid, p.DayStamp)
	if currErr != nil { //如果没有则插入新的
		// log.Warningf("AddOrUpdatePersona currErr:%s basemgoErr:%s", currErr.Error(), basemgo.GetErrNotFound().Error())
		if currErr.Error() == basemgo.GetErrNotFound().Error() {
			p.ID = bson.NewObjectId()
			err = con.Insert(p)
			jsonBytes, _ := json.Marshal(p)
			// log.Printf("save ok", string(jsonBytes))
			doneUser = string(jsonBytes)
		} else {
			err = currErr
		}

	} else {
		p.ID = personaDB.ID
		err = con.Update(bson.M{"_id": p.ID}, p)

		jsonBytes, _ := json.Marshal(p)
		// log.Printf("save ok", string(jsonBytes))
		doneUser = string(jsonBytes)

	}

	return
}

func FindPersona(gid float64, dayStamp string) (*PersonaDB, error) {
	var user PersonaDB
	con := basemgo.AnalyseDB.GetCon("persona", user.makeIndex)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   db.user.find({username:"13888888888"})*/
	if err := con.Find(bson.M{"gid": gid, "dayStamp": dayStamp}).One(&user); err != nil {
		// if err.Error() != basemgo.GetErrNotFound().Error() {
		if err != nil {
			return &user, err
		}

	}
	return &user, nil
}

// 查询该用户的所有数据
func FindPersonaAll() (*[]PersonaDB, error) {
	var users []PersonaDB
	con := basemgo.AnalyseDB.GetCon("persona", nil)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   db.user.find({username:"13888888888"})*/
	if err := con.Find(bson.M{}).All(&users); err != nil {
		// if err.Error() != basemgo.GetErrNotFound().Error() {
		if err != nil {
			return &users, err
		}

	}
	return &users, nil
}
func FindPersonaByGid(gid float64) (*[]PersonaDB, error) {
	var users []PersonaDB
	con := basemgo.AnalyseDB.GetCon("persona", nil)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   db.user.find({username:"13888888888"})*/
	if err := con.Find(bson.M{"gid": gid}).All(&users); err != nil {
		// if err.Error() != basemgo.GetErrNotFound().Error() {
		if err != nil {
			return &users, err
		}

	}
	return &users, nil
}

func FindPersonaByTotal(totalS float64, dayStamp string) (*[]PersonaDB, error) {
	var users []PersonaDB
	con := basemgo.AnalyseDB.GetCon("persona", nil)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   query := collection.Find(bson.M{"_id": bson.M{"$gt": lastId}})*/

	if err := con.Find(bson.M{"totalS": bson.M{"$lt": totalS}, "dayStamp": dayStamp}).All(&users); err != nil {
		// if err.Error() != basemgo.GetErrNotFound().Error() {
		if err != nil {
			return &users, err
		}

	}
	return &users, nil
}

// // 每日数据删除,不知道为毛无效
// func DailyClean() {
// 	crontab := cron.New()
// 	log.Debugln("what the fucking...")
// 	entryId, _ := crontab.AddFunc("1 * * * * *", func() {
// 		log.Debugln("Run models.CleanAllTag...")
// 		fmt.Println("Every hour on the half hour")
// 	}) // 每天 8:20:00 定时执行 myfunc 函数
// 	crontab.AddFunc("30 * * * * *", func() { fmt.Println("Every hour on the half hour") })
// 	log.Debugln("what the fucking222...", entryId)
// 	crontab.Start()
// 	log.Debugln("what the fucking333...")
// }

// DailyClean 每日数据删除
func DailyClean() {
	go func() {
		for {
			if basemgo.AnalyseDB != nil {
				con := basemgo.AnalyseDB.GetCon("persona", nil)
				if con != nil {
					// 24*60 = 1440
					during, _ := time.ParseDuration("-1440h")
					newTimeStr := time.Now().Add(during).Format("20060102")
					log.Infoln("执行DailyClean 日期:", newTimeStr)
					changInfo, delErr := con.RemoveAll(bson.M{"dayStamp": bson.M{"$lt": newTimeStr}})
					log.Infoln("RemoveAll 执行结果:", changInfo, delErr)
					if delErr != nil && delErr.Error() != basemgo.GetErrNotFound().Error() {
						// 直接打印错误日志 直接打印日志
						log.Error("DailyClean err", delErr.Error())
					}

				}
			}
			time.Sleep(time.Hour * 1440)
		}
	}()
}
