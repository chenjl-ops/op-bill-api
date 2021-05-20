package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"op-bill-api/internal/app/billing"
	"op-bill-api/internal/app/compute"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/app/prediction"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/mysql"
	"op-bill-api/internal/pkg/rdssentinels"
	"op-bill-api/internal/pkg/snowflake"
)

/*
TODO
1、Event跨实例间通讯
*/

// 初始化 apollo config
func initApolloConfig() {
	var err error
	err = apollo.ReadRemoteConfig()
	if nil != err {
		logrus.Fatal(err)
	}
}

// 初始化mysql
func initMysql() {
	err := mysql.NewMysql()
	if nil != err {
		logrus.Fatal(err)
	}
}

// 初始化redis
func initRedis() {
	rdssentinels.NewRedis(nil)
}

// 初始化log配置
func (s *server) initLog() *gin.Engine {
	logs := logger.LogMiddleware()
	s.App.Use(logs)
	return s.App
}

// 初始化雪花算法
func initSnowFlake() {
	err := snowflake.InitSnowWorker(1, 1)
	if err != nil {
		logrus.Fatal(err)
	}
}

//  初始化账单job
//func initBill() error {
//	err := billing.JobService()
//	return err
//}

// InitRouter 加载gin 路由配置
func (s *server) InitRouter() *gin.Engine {
	billing.Url(s.App)    // 决算
	prediction.Url(s.App) // 预算
	compute.Url(s.App)
	return s.App
}

// InitSwagger init swagger
func (s *server) InitSwagger() *gin.Engine {
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	s.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return s.App
}
