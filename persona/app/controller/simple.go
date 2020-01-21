package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lew/persona/app/model"
	"github.com/lew/persona/app/model/redis"
)

type Result struct {
	Code string       `json:"code"`
	Data []model.User `json:"data"`
}

func Home(c *gin.Context) {
	// value, exist := c.GetQuery("key")
	// if !exist {
	// 	value = "the key is not exist!"
	// }
	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("welcome to go persona")))
	return
}

// func Middleware(c *gin.Context) {
// 	fmt.Println("this is a middleware!")
// }

func GetHandler(c *gin.Context) {
	value, exist := c.GetQuery("key")
	if !exist {
		value = "the key is not exist!"
	}
	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("get success! %s\n", value)))
	return
}
func PostHandler(c *gin.Context) {
	type JsonHolder struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	holder := JsonHolder{Id: 1, Name: "my name"}
	//若返回json数据，可以直接使用gin封装好的JSON方法
	c.JSON(http.StatusOK, holder)
	return
}
func PutHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("put success!\n"))
	return
}
func DeleteHandler(c *gin.Context) {
	c.Data(http.StatusOK, "text/plain", []byte("delete success!\n"))
	return
}

func RedisGet(c *gin.Context) {
	value, exist := c.GetQuery("key")
	fmt.Println(c.Request.URL.Query())
	if !exist {
		value = "the key is not exist!"
	} else {
		redis.PersonaRD.GetStringValue("testKey")
	}
	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("get success! %s\n", value)))
	return
}
func RedisSet(c *gin.Context) {
	value, exist := c.GetQuery("key")

	if !exist {
		value = "the key is not exist!"
	} else {
		redis.PersonaRD.Set("testKey", value)
	}
	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("get success! %s\n", value)))
	return
}

func MongoSet(c *gin.Context) {
	password, exist1 := c.GetQuery("password")
	username, exist2 := c.GetQuery("username")
	var theuser string
	var err error
	if !exist1 || !exist2 {
		password = "the key is not exist!"
	} else {
		theuser, err = model.AddUser(password, username)
		// theuser = fuck
		// log.Println("mongoget", err)
	}
	log.Println("mongoget", theuser)
	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("set success! %s\n", theuser, err)))
	return
}
func MongoGet(c *gin.Context) {
	username, exist := c.GetQuery("username")
	var theuser model.User
	var err error

	if !exist {
		username = "the key is not exist!"
	} else {
		theuser, err = model.FindUser(username)
	}

	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("get success! %s\n", theuser, err)))
	return
}
func MongoGetAll(c *gin.Context) {
	// var theuser []model.User
	// var err error

	theuser, err := model.FindAllUser()
	// jsonBytes, err2 := json.Marshal(theuser)
	var code string = "OK"
	if err != nil {
		code = err.Error()
	}
	result := Result{code, theuser}
	resultStr, _ := json.Marshal(result)

	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("%s", resultStr)))
	return
}

func MongoUpdate(c *gin.Context) {
	// var theuser []model.User
	// var err error
	username, exist := c.GetQuery("username")
	err := model.UpdateUser(username)
	var result string = "get success!"
	if !exist {
		result = "query key not exist"
	}
	if err != nil {
		result = "have err" + err.Error()
	}

	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("s%", result)))
	return
}
func MongoDel(c *gin.Context) {
	// var theuser []model.User
	// var err error
	username, exist := c.GetQuery("username")
	err := model.DelUser(username)
	var result string = "get success!"
	if !exist {
		result = "query key not exist"
	}
	if err != nil {
		result = "have err" + err.Error()
	}

	c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("s%", result)))
	return
}

func CrossTest() {
	log.Print("fuck simple.go")
}
