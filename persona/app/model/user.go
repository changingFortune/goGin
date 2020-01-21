package model

import (

	// "time"
	"encoding/json"

	"github.com/lew/persona/app/model/basemgo"
	"gopkg.in/mgo.v2/bson"
)

//个人项目部分代码
type User struct {
	ID       bson.ObjectId `json:"_id" bson:"_id"`
	UserName string        `json:"username" bson:"username"`
	Summary  string        `json:summary bson:"summary"`
	Age      int           `json:"age" bson:"age"`
	Phone    int           `json:"phone" bson:"phone"`
	PassWord string        `json:"password" bson:"password"`
	Sex      int           `json:"sex" bson:"sex"`
	Name     string        `json:"name" bson:"name"`
	Email    string        `json:"email" bson:"email"`
}

func AddUser(password string, username string) (doneUser string, err error) {
	con := basemgo.AnalyseDB.GetCon("user", nil)
	//可以添加一个或多个文档
	/* 对应mongo命令行
		   db.user.insert({username:"13888888888",summary:"code",
	       age:20,phone:"13888888888"})*/
	tempU := &User{ID: bson.NewObjectId(), UserName: username, PassWord: password}
	err = con.Insert(tempU)

	jsonBytes, err := json.Marshal(tempU)

	// log.Printf("save ok", string(jsonBytes))
	doneUser = string(jsonBytes)
	return
}

func FindUser(username string) (User, error) {
	var user User
	con := basemgo.AnalyseDB.GetCon("user", nil)

	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   db.user.find({username:"13888888888"})*/
	if err := con.Find(bson.M{"username": username}).One(&user); err != nil {
		if err.Error() != basemgo.GetErrNotFound().Error() {
			return user, err
		}

	}
	return user, nil
}
func FindAllUser() (users []User, err error) {
	// var users []User
	con := basemgo.AnalyseDB.GetCon("user", nil)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   db.user.find({})*/
	if findErr := con.Find(bson.M{}).All(&users); findErr != nil {
		if findErr.Error() != basemgo.GetErrNotFound().Error() {
			return users, findErr
		}

	}
	return users, nil
}

func UpdateUser(username string) (err error) {
	con := basemgo.AnalyseDB.GetCon("user", nil)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行 $set
	 */
	//  bson.ObjectIdHex("5204af979955496907000001")
	if err := con.Update(bson.M{"username": username},
		bson.M{"$set": bson.M{
			"name": "Jimmy Gu",
			"age":  34,
		}}); err != nil {
		if err.Error() != basemgo.GetErrNotFound().Error() {
			return err
		}
	}
	return nil
}

func DelUser(username string) error {
	// var users []User
	con := basemgo.AnalyseDB.GetCon("user", nil)
	//通过bson.M(是一个map[string]interface{}类型)进行
	//条件筛选，达到文档查询的目的
	/* 对应mongo命令行
	   db.user.find({})*/
	if err := con.Remove(bson.M{"username": username}); err != nil {
		if err.Error() != basemgo.GetErrNotFound().Error() {
			return err
		}

	}
	return nil
}
