package main

import (
	log "github.com/Sirupsen/logrus"
	_ "op-bill-api/api"
	"op-bill-api/internal/app/server"
)

/*
swagger 规范

// @Summary 摘要
// @Description 描述
// @Description 接口的详细描述
// @Id 全局标识符
// @Version 接口版本号
// @Tags 接口分组，相当于归类
// @Accept  json 浏览器可处理数据类型
// @Produce  json 设置返回数据的类型和编码
// @Param   参数格式 从左到右：参数名、入参类型、数据类型、是否必填和注释 案例：id query int true "ID"
// @Success 响应成功 从左到右：状态码、参数类型、数据类型和注释  案例：200 {string} string    "ok"
// @Failure 响应失败 从左到右：状态码、参数类型、数据类型和注释  案例：400 {object} web.APIError "We need ID!!"
// @Router 路由： 地址和http方法   案例：/testapi/get-string-by-int/{some_id} [get]
// @contact.name 接口联系人
// @contact.url 联系人网址
// @contact.email 联系人邮箱


// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host apollo.Config.SelfDomain
// @BasePath

// @in header
// @name Authorization
*/

// @title Op-bill-api API
// @version 1.0
// @description This is op-bill-api api server.
func main() {
	log.Info("Start op-bill-api Service ....")
	err := server.StartServer()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Start Service Successful")
}
