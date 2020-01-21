package model

import (
	"encoding/json"
	"strconv"

	"github.com/lew/persona/app/model/basemgo"
	"github.com/lew/persona/app/model/redis"
	"gopkg.in/mgo.v2/bson"
)

// "time"
// "encoding/json"

// 静态信息:
// 名字,性别,地区,头像
// { "_id" : 200021, "nickname" : "%u5927%u54E5%u522B%u6740%u6211", "headimgurl" : "http://thirdwx.qlogo.cn/mmopen/vi_32/Y0TIjQQnfqTqIAj2ib8XMfCibxNgm65NLtHHzQQBlDfpaQlvRIzR7MLAHaBojA8fgkg9xfziaj3RuxDib83okiayLicA/132", "openid" : "o76rawS5e2R_N7DKgO-ZveEF_maQ", "unionid" : "oA-7h0h8xeyar7Wz11rRHZwuOcM8", "lType" : "wx", "refresh_token" : "22_Xn8bxbaurzDo9gWtcOgRFmqWp1nMgieU78bEAQuelfPingchaSlknbv0TBossTUdRMkej6tafC-7PysbxoX1VZs6scRPRgOR3aknD1dAFUY", "app" : { "appid" : "com.coolgamebox.majiang", "os" : "Android" }, "gamekind" : "default", "uid" : 200021, "name" : "8klq9v", "loginCode" : "8i38uz", "face" : "avata:65", "members" : {  }, "sendTime" : ISODate("2019-05-22T09:44:46.691Z"), "resVersion" : "2.57.001", "gid" : 36016268 }
type CgbUser struct {
	Gid float64 `json:"gid" bson:"gid"`
	// AppId        float64            `json:"appId" bson:"appId"`
	UserId     float64 `json:"uid" bson:"uid"`
	Nickname   string  `json:"nickname" bson:"nickname"`     //昵称
	Name       string  `json:"name" bson:"name"`             //名字
	Headimgurl string  `json:"headimgurl" bson:"headimgurl"` //头像
}

func FindCgbUser(appId float64, userId float64) (*CgbUser, error) {
	var user CgbUser
	if basemgo.HasDBCon(appId) {
		con := basemgo.GetDBCon(appId, "cgbuser")
		if err := con.Find(bson.M{"_id": userId}).One(&user); err != nil {
			// if err.Error() != basemgo.GetErrNotFound().Error() {
			if err != nil {
				return &user, err
			}
		}
	}

	return &user, nil
}

// redis中获取cgb数据
func FindCgbUserByRD(appId float64, userId float64) (*map[string]string, error) {
	var user CgbUser = CgbUser{}
	cgbByte, err := redis.GetRDIns(appId).HGETALL(string("cgbuser." + strconv.FormatFloat(userId, 'f', 0, 64)))
	if err == nil {
		// cgbByte.Decode(mapInstance, &person)
		// 转map string 会导致 number类型数据丢失,fuck
		jsonStr, err := json.Marshal(cgbByte)
		err = json.Unmarshal(jsonStr, &user)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return &cgbByte, nil
}

// type PersonaDB struct { // mongo对象
// 	Persona       `bson:",inline"`
// 	ID            bson.ObjectId `json:"_id" bson:"_id"`
// 	GameCountS    float64       `json:"gameCountS" bson:"gameCountS"`       //游戏把数 评分
// 	PlayDurationS float64       `json:"playDurationS" bson:"playDurationS"` //打牌时长 评分
// 	DurationS     float64       `json:"durationS" bson:"durationS"`         //游戏总时长 评分
// 	FriendsS      float64       `json:"friendsS" bson:"friendsS"`           //好友列表 评分
// 	TotalS        float64       `json:"totalS" bson:"totalS"`               //总分

// }

// func AddOrUpdatePersona(p *PersonaDB) (doneUser string, err error) {

// 	con := basemgo.AnalyseDB.GetCon("persona", p.makeIndex)
// 	testIndex, theErr := con.Indexes()
// 	// con.EnsureIndexKey
// 	log.Println(testIndex, theErr, len(testIndex[:]), len(testIndex))

// 	//可以添加一个或多个文档
// 	/* 对应mongo命令行
// 	   db.user.insert({username:"13888888888",summary:"code",
// 	   age:20,phone:"13888888888"})*/
// 	personaDB, currErr := FindPersona(p.Gid, p.DayStamp)
// 	if currErr != nil { //如果没有则插入新的
// 		p.ID = bson.NewObjectId()
// 		err = con.Insert(p)
// 		// con.Find()
// 		jsonBytes, _ := json.Marshal(p)
// 		log.Printf("save ok", string(jsonBytes))
// 		doneUser = string(jsonBytes)
// 	} else {
// 		p.ID = personaDB.ID
// 		err = con.Update(bson.M{"_id": p.ID}, p)

// 		jsonBytes, _ := json.Marshal(p)
// 		log.Printf("save ok", string(jsonBytes))
// 		doneUser = string(jsonBytes)

// 	}

// 	return
// }

// func FindPersona(gid float64, dayStamp string) (*PersonaDB, error) {
// 	var user PersonaDB
// 	con := basemgo.AnalyseDB.GetCon("persona", user.makeIndex)
// 	//通过bson.M(是一个map[string]interface{}类型)进行
// 	//条件筛选，达到文档查询的目的
// 	/* 对应mongo命令行
// 	   db.user.find({username:"13888888888"})*/
// 	if err := con.Find(bson.M{"gid": gid, "dayStamp": dayStamp}).One(&user); err != nil {
// 		// if err.Error() != basemgo.GetErrNotFound().Error() {
// 		if err != nil {
// 			return &user, err
// 		}

// 	}
// 	return &user, nil
// }

// // 查询该用户的所有数据
// func FindPersonaAll(gid float64) (*[]PersonaDB, error) {
// 	var users []PersonaDB
// 	con := basemgo.AnalyseDB.GetCon("persona", nil)
// 	//通过bson.M(是一个map[string]interface{}类型)进行
// 	//条件筛选，达到文档查询的目的
// 	/* 对应mongo命令行
// 	   db.user.find({username:"13888888888"})*/
// 	if err := con.Find(bson.M{"gid": gid}).All(&users); err != nil {
// 		// if err.Error() != basemgo.GetErrNotFound().Error() {
// 		if err != nil {
// 			return &users, err
// 		}

// 	}
// 	return &users, nil
// }
