package billing

import (
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
)

// 创建数据表
func createTable(c *gin.Context) {
	err := mysql.Engine.Sync2(new(config.ShareBill), new(config.SourceBill), new(config.BillStatus))
	if err != nil {
		c.JSON(500, gin.H{
			"msg":   "failed",
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"msg": "success",
		})
	}
}

// 账单数据录入
func insertData(c *gin.Context) {
	err := getBillExcel()
	if err != nil {
		c.JSON(500, gin.H{
			"msg":   "failed",
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"msg": "success",
		})
	}
}

// 查看月首尾日期
func getMonthData(c *gin.Context) {
	data := getMonthDate()
	c.JSON(200, gin.H{
		"data": data,
	})
}
