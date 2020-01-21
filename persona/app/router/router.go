package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lew/persona/app/controller"
	"github.com/lew/persona/app/middleware"
)

func Register(env *string, ports *string) {
	fmt.Println("Register 参数 env:%s ports:%s", *env, *ports)
	// cfgEvn := config.Cfg.MustValue("evn", "model", "debug")
	if *env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else if *env == gin.TestMode {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 禁用控制台颜色
	gin.ForceConsoleColor()

	// // LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// // By default gin.DefaultWriter = os.Stdout
	// router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

	// 	// your custom format
	// 	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
	// 		param.ClientIP,
	// 		param.TimeStamp.Format(time.RFC1123),
	// 		param.Method,
	// 		param.Path,
	// 		param.Request.Proto,
	// 		param.StatusCode,
	// 		param.Latency,
	// 		param.Request.UserAgent(),
	// 		param.ErrorMessage,
	// 	)
	// }))

	//监听端口
	// router := createRouter()
	// http.ListenAndServe(":9008", router)

	// 多端口
	// ports := []string{":9008", ":9009"}
	// ports := config.Cfg.MustValueArray("go", "ports", ",")
	portsArr := strings.Split(string(*ports), ",")
	for _, v := range portsArr {
		go func(port string) { //每个端口都扔进一个goroutine中去监听
			fmt.Println("参数 port:" + port + " v:" + v)
			router := createRouter()
			http.ListenAndServe(":"+port, router)
		}(v)
	}
	fmt.Println("后续代码不会执行", portsArr)
	select {}
	fmt.Println("后续代码不会执行")
}

func createRouter() *gin.Engine {
	// router := gin.Default() //获得路由实例
	router := gin.New() //获得路由实例

	//添加中间件
	router.Use(middleware.TheTest)
	router.Use(gin.Recovery())
	//注册接口
	router.GET("/", controller.Home)
	router.GET("/simple/server/get", controller.GetHandler)
	router.POST("/simple/server/post", controller.PostHandler)
	router.PUT("/simple/server/put", controller.PutHandler)
	router.DELETE("/simple/server/delete", controller.DeleteHandler)

	router.GET("/redis/get", controller.RedisGet)
	router.GET("/redis/set", controller.RedisSet)

	router.GET("/mongo/get", controller.MongoGet)
	router.GET("/mongo/set", controller.MongoSet)
	router.GET("/mongo/getall", controller.MongoGetAll)
	router.GET("/mongo/update", controller.MongoUpdate)
	router.GET("/mongo/del", controller.MongoDel)

	// 正式接口
	router.POST("/persona/commonEvent", controller.CommonEvent)
	router.GET("/persona/getPersonaAll", controller.GetPersonaAll)
	router.GET("/persona/getPersonaByGid", controller.GetPersonaByGid)
	router.GET("/persona/getAnalyseByGid", controller.GetAnalyseByGid)
	router.GET("/persona/getRemoteTest", controller.GetRemoteTest)
	router.GET("/persona/checkAlarm", controller.CheckAlarm)
	return router
}
