package controller

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lew/persona/app/config"
	"github.com/lew/persona/app/model"
	"github.com/lew/persona/app/model/redis"
	log "github.com/sirupsen/logrus"
)

// type KDRespBody struct {
//       Errcode int `json:"errcode"`
//       Desc string `json:"description"`
//       Data []services.KdSearchBack `json:"data"`
// }
const (
	Login    = "login"
	Register = "register"
	Logout   = "logout"
	GameOver = "gameOver"
)

var redisHname = "persona"
var expireHname = ""

// 创建时间
// 每天游戏时长
// 每天把数
// 每天输赢次数
// 游戏类型(ddz,mj)
// 充值记录
// 消耗元宝或房卡数

// 次要信息:
// 性别
// 同桌信息
// 俱乐部相关信息(人数,活跃,消耗速度)
// 反馈信息
// 比赛参与度
// 页游和活动等点击度
// GPS开启状态

func CommonEvent(c *gin.Context) {
	// buf := make([]byte, 1024000)
	// n, _ := c.Request.Body.Read(buf)
	// log.Printf("数据不对? 老数据转换:%d %s", n, string(buf[0:n]))
	// ------------array 解析------------
	body, _ := ioutil.ReadAll(c.Request.Body)
	// log.Printf("数据不对? 收到数据:%s", string(body))
	unCompressBody := DoZlibUnCompress(body)
	// log.Printf("数据不对? 解压后数据:%s", string(unCompressBody))
	//json str 转map
	// var arr []string
	// err := json.Unmarshal(unCompressBody, &arr)
	arr := strings.Split(string(unCompressBody), "\n")
	// string(buf[0:n])
	var code string = "OK"
	// var eventMap map[string]interface{}
	log.Printf("消息长度 %d", len(arr))
	for _, strVal := range arr {
		var jsonVal map[string]interface{}
		err := json.Unmarshal([]byte(strVal), &jsonVal)
		if err != nil {
			code = err.Error()
			log.Error(code)
			break
		}
		timeStr := jsonVal["time"].(string)
		if timeStr == "" {
			timeStr = string(time.Now().UnixNano() / 1e6)
		}
		beginUnix := time.Now().UnixNano() / 1e6

		// 分发事件
		var eventType = jsonVal["type"].(string)
		// if eventType == "Login" || eventType == "Logout" {
		distribute(eventType, timeStr, jsonVal["properties"].(map[string]interface{}))
		// } else {
		// 	go distribute(eventType, timeStr, jsonVal["properties"].(map[string]interface{}))
		// }

		endUnix := time.Now().UnixNano() / 1e6
		log.Printf("耗时 %d", endUnix-beginUnix)
		// 设置过期时间 2天
		if expireHname != getHname() {
			err := redis.PersonaRD.SetExpire(getHname(), 2*24*60*60, false)
			if err == nil {
				expireHname = getHname()
			}
		}
	}

	resp := map[string]interface{}{"code": code}
	c.JSON(http.StatusOK, resp)

	// -------------json转换和解析代码---------------
	// if err == nil {
	// 	fmt.Println("==============json str 转map=======================")
	// 	fmt.Println(dat)
	// 	fmt.Println(dat["properties"].(map[string]interface{})["user"].(map[string]interface{})["userId"])
	// } else {
	// 	fmt.Println("Unmarshal have err", err.Error())
	// }
	// resp := map[string]interface{}{"code": "OK", "data": dat}
	// c.JSON(http.StatusOK, resp)
	// c.JSON(200, gin.H{"errcode": 400, "description": "Post Data Err"})
	/*post_gwid := c.PostForm("name")
	fmt.Println(post_gwid)*/
}
func GetPersonaAll(c *gin.Context) {
	var code string = "OK"
	users, err := model.FindPersonaAll()
	if err != nil {
		code = err.Error()
	}

	resp := map[string]interface{}{"code": code, "users": users}
	c.JSON(http.StatusOK, resp)
}
func GetPersonaByGid(c *gin.Context) {
	var code string = "OK"
	gid := c.Query("gid")
	gidVal, _ := strconv.ParseFloat(gid, 64)
	personas, err := model.FindPersonaByGid(gidVal)
	if err != nil {
		code = err.Error()
	}

	resp := map[string]interface{}{"code": code, "personas": personas}
	c.JSON(http.StatusOK, resp)
}

// // 单人评分
// func GetAnalyseOne(c *gin.Context) {
// 	var code string = "OK"
// 	var resp interface{}
// 	appId := c.Query("appId")
// 	userId := c.Query("userId")

// 	appIdVal, _ := strconv.ParseFloat(appId, 64)
// 	userIdVal, _ := strconv.ParseFloat(userId, 64)
// 	cgbUser, personas, err := FindAnalyseOne(appIdVal, userIdVal)
// 	if err != nil {
// 		code = err.Error()
// 	}
// 	resp = map[string]interface{}{"code": code, "cgbUser": cgbUser, "persona": personas}
// 	c.JSON(http.StatusOK, resp)
// }
func GetAnalyseByGid(c *gin.Context) {
	var code string = "OK"
	var resp interface{}
	gid := c.Query("gid")

	// cgbUser, personas := FindAnalyseOne(appId, userId)
	gidVal, _ := strconv.ParseFloat(gid, 64)
	personas, err := model.FindPersonaByGid(gidVal)
	var cgbUser *map[string]string
	if personas != nil && len(*personas) > 0 {
		tempVal := (*personas)[0]
		cgbUser, err = model.FindCgbUserByRD(tempVal.AppId, tempVal.UserId)
	}

	if err != nil {
		code = err.Error()
	}
	resp = map[string]interface{}{"code": code, "cgbUser": cgbUser, "persona": personas}
	c.JSON(http.StatusOK, resp)
}

func GetRemoteTest(c *gin.Context) {
	appId := c.Query("appId")
	userId := c.Query("userId")

	appIdVal, _ := strconv.ParseFloat(appId, 64)
	// userIdVal, _ := strconv.ParseFloat(userId, 64)

	var code string = "OK"

	// fuckingVal := redis.GetRDIns(200005).GetStringValue("activity:20190219::keyv:default:1000054")
	// fuckingVal2 := redis.GetRDIns(200003).GetStringValue("activity:20190219::keyv:default:1000036")

	var user model.CgbUser = model.CgbUser{}
	// cgbMapStr, err := redis.PersonaRD.HGETALL("cgbuser.200021")
	cgbMapStr, err := redis.GetRDIns(appIdVal).HGETALL("cgbuser." + userId)
	if err != nil {
		code = err.Error()
	} else {
		jsonStr, _ := json.Marshal(cgbMapStr)
		err = json.Unmarshal(jsonStr, &user)
	}

	var resp = map[string]interface{}{"code": code, "cgbByte": cgbMapStr, "user": user}
	c.JSON(http.StatusOK, resp)
}

// 检测报警
func CheckAlarm(c *gin.Context) {
	totalS := c.Query("totalS")
	dayStamp := c.Query("dayStamp")
	if dayStamp == "" {
		dayStamp = getDayStamp()
	}
	totalSVal, _ := strconv.ParseFloat(totalS, 64)

	var code string = "OK"
	personas, err := model.FindPersonaByTotal(totalSVal, dayStamp)
	if err != nil {
		code = err.Error()
	} else {
		// jsonStr, _ := json.Marshal(cgbMapStr)
		// err = json.Unmarshal(jsonStr, &user)
	}
	// if personas != nil && len(*personas) > 0 {
	// 	tempVal := (*personas)[0]
	// 	cgbUser, err = model.FindCgbUserByRD(tempVal.AppId, tempVal.UserId)
	// }

	var resp = map[string]interface{}{"code": code, "personas": personas}
	c.JSON(http.StatusOK, resp)
}

// 分发事件
func distribute(eventType string, timeStr string, properties map[string]interface{}) {
	currTime, _ := strconv.ParseFloat(timeStr, 64)
	log.Debugf("事件类型%s", eventType)
	switch eventType {
	case Login: //	采集登录事件,创建当天用户数据

		user := properties["user"].(map[string]interface{})
		persona, _ := getOrCreateP(user)

		if persona.FirstLogTime == "" {
			persona.FirstLogTime = timeStr
		}
		redis.PersonaRD.SetHSET(getHname(), strconv.FormatFloat(user["gid"].(float64), 'f', 0, 64), persona)
	case Register:
	case GameOver: //	1.把数 2.游戏时长	3.输赢次数	4.输赢分数	5.好友数据
		startTime, _ := strconv.ParseFloat(properties["startTime"].(string), 64)
		playLen := currTime - startTime

		// 取出相关数据
		totalScores := properties["totalScores"].([]interface{})
		users := properties["users"].(map[string]interface{})
		gids := users["gids"].([]interface{})
		userIds := users["userIds"].([]interface{})
		appId := users["appId"].(float64)

		for i, j := range gids {
			// var tempUser map[string]interface{}
			tempUser := make(map[string]interface{})

			// 获取redis数据
			jVal := j.(float64)
			tempUser["gid"] = jVal
			tempUser["userId"] = userIds[i].(float64)
			tempUser["appId"] = appId
			tempScore := totalScores[i].(float64)
			var win float64 = -1
			if tempScore >= 0 {
				win = 1
			}

			persona, _ := getOrCreateP(tempUser)
			persona.GameCount = persona.GameCount + 1
			persona.PlayDuration = persona.PlayDuration + playLen
			persona.WinScore = persona.WinScore + tempScore
			persona.WinCount = persona.WinCount + win
			for _, l := range gids {
				if l != j {
					log.Println(persona.Friends, len(persona.Friends))
					if len(persona.Friends) == 0 {
						persona.Friends = make(map[string]float64)
					}
					log.Println(persona.Friends, len(persona.Friends))
					// redis的map[float64]float64格式无法解析到json,试一下map[string]float64
					gidStr := strconv.FormatFloat(float64(l.(float64)), 'f', 0, 64)
					persona.Friends[gidStr] = persona.Friends[gidStr] + 1
				}
			}

			err := redis.PersonaRD.SetHSET(getHname(), strconv.FormatFloat(jVal, 'f', 0, 64), persona)
			if err != nil {
				log.Println("GameOver change errr", err.Error())
			}

		}

		// log.Println(playLen, users, totalScores, gids, userIds, appId)

	case Logout: // 1.一分钟以上的间隔(防止频繁掉线导致的性能消耗)  2.判断时间,进行数据入库转移 3.后续可添加数据删除和分析触发

		// currTime, _ := strconv.ParseFloat(timeStr, 64)
		user := properties["user"].(map[string]interface{})
		persona, _ := getOrCreateP(user)
		outTime, _ := strconv.ParseFloat(persona.LogoutTime, 64)
		var limit float64 = config.Cfg.MustFloat64("common", "saveLimit", 0) //1 * 60 * 1000
		// log.Debugf("离线限制: currTime:%f outTime:%f  %f 配置限制:%f", currTime, outTime, (currTime - outTime), limit)
		if (currTime - outTime) >= limit {
			log.Debugf("通过限制!")
			persona.LogoutTime = timeStr
			// 累加游戏时间
			loginTime, _ := strconv.ParseFloat(properties["loginTime"].(string), 64)
			persona.Duration = persona.Duration + (currTime - loginTime)
			// log.Println(persona)
			redis.PersonaRD.SetHSET(getHname(), strconv.FormatFloat(user["gid"].(float64), 'f', 0, 64), persona)
			// 登出时数据移库到mongo
			// tempStr, err := model.AddOrUpdatePersona(&model.PersonaDB{*persona, ""})
			// log.Println("to mongo tempStr", tempStr, err)
			// 触发数据分析逻辑
			CalculateS(persona)
		}

	}
}
func getHname() string {
	return redisHname + "_" + getDayStamp()
}
func getDayStamp() string {
	return time.Now().Format("20060102")
}
func getOrCreateP(user map[string]interface{}) (persona *model.Persona, have bool) {
	have = false
	persona = getPersona(user["gid"].(float64))
	if persona == nil {
		persona = &model.Persona{}
		persona.Gid = user["gid"].(float64)
		persona.UserId = user["userId"].(float64)
		persona.AppId = user["appId"].(float64)
		persona.DayStamp = getDayStamp()
		have = true
	}
	return
}
func getPersona(gid float64) *model.Persona {
	personaByte, err := redis.PersonaRD.GetHGET(getHname(), strconv.FormatFloat(gid, 'f', 0, 64))
	// stu:=StuRead{}

	if err == nil {
		persona := model.Persona{}
		err = json.Unmarshal(personaByte, &persona)
		return &persona
	} else {
		log.Println("getPersona Unmarshal", err.Error())
		return nil
	}

}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	r.Close()
	return out.Bytes()
}
