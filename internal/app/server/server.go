package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator"
	"op-bill-api/internal/pkg/apollo"
)

type server struct {
	//Config   *apollo.Specification
	App *gin.Engine
	//Validate *validator.Validate
}

func NewServer() (*server, error) {
	return &server{
		//Config: globalConfig,
		App: gin.New(),
	}, nil
}

func StartServer() error {
	initApolloConfig()
	initMysql()
	initRedis()
	initSnowFlake()

	//billErr := initBill()
	//if billErr != nil {
	//	return billErr
	//}

	server, err1 := NewServer()
	// server.App.Use(Cors())
	if err1 != nil {
		return err1
	}
	// 初始化日志
	server.initLog()
	// 初始化swagger
	server.InitSwagger()
	// 初始化路由
	server.InitRouter()

	//启动服务
	err := server.Run()
	if err != nil {
		return err
	}
	return nil
}

// Run 启动服务
func (s *server) Run() error {
	return s.App.Run(fmt.Sprintf("0.0.0.0:%s", apollo.Config.ListenPort))
}

// 跨域设置
//func Cors() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.Header("Access-Control-Allow-Origin", "*")
//		c.Next()
//	}
//}
