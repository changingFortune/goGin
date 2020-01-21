package controller

import (
	log "github.com/sirupsen/logrus"

	"github.com/lew/persona/app/model"
)

// 数据分析相关方法
// 计算玩家得分
func CalculateS(p *model.Persona) {
	if p.Gid == 0 {
		log.Warningf("the persona gid == 0 %v", *p)
		// return
	}
	// 对数据进行评分,然后存入mongo数据库
	// 游戏把数	3 打牌时长	3游戏总时长	2好友列表	2
	personaDB := model.PersonaDB{*p, "", 0, 0, 0, 0, 0}
	personaDB.GameCountS = makeScore(3, p.GameCount, 3, 8, 9)
	personaDB.PlayDurationS = makeScore(3, p.PlayDuration, 45*60*1000, 120*60*1000, 9)
	personaDB.DurationS = makeScore(3, p.Duration, 45*60*1000*7, 120*60*1000*7, 9)
	personaDB.FriendsS = makeScore(3, float64(len(p.Friends)), 4, 12, 9)
	personaDB.TotalS = personaDB.GameCountS + personaDB.PlayDurationS + personaDB.DurationS + personaDB.FriendsS
	_, err := model.AddOrUpdatePersona(&personaDB)
	if err != nil {
		log.Error("save err", err.Error())
	}
	log.Debugf("gameover save to db!!!! err:%s ", err)

}

// // 单人评分
// func FindAnalyseOne(appId float64, userId float64) (cgbUser *model.CgbUser, personas *[]model.PersonaDB, err error) {

// 	// mongo中获取
// 	cgbUser, err = model.FindCgbUser(appId, userId)
// 	// redis中获取
// 	// cgbUser, err = model.FindCgbUserByRD(appId, userId)

// 	if err != nil {
// 		log.Error("analyzeOne err", err.Error())
// 	} else {
// 		// 获取该用户的所有评分
// 		personas, err = model.FindPersonaAll(cgbUser.Gid)
// 		if err != nil {
// 			log.Error("analyzeOne err", err.Error())
// 		}

// 	}
// 	return
// }

//weight:权重
func makeScore(weight float64, sourceVal float64, low float64, middl float64, high float32) float64 {
	rate := 1.0
	if sourceVal <= low {
		rate = (sourceVal / low) * 0.6
	} else if sourceVal > low && sourceVal <= middl {
		rate = (sourceVal / middl) * 0.8
	} else {
		rate = 1
	}
	return weight * rate
}
