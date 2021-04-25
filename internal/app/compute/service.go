package compute

import (
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
)


func GetBilling(c *gin.Context) {
	month := c.DefaultQuery("month", "")
	if month == "" {
		c.JSON(500, gin.H{
			"msg": "parameter month not none",
		})
	} else {
		data, err := ComputerBilling(month)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": err,
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "success",
				//"data": data,
				"length": len(data),
			})
		}
	}

}

// 计算决算数据
func ComputerBilling(month string) ([]config.ShareBill ,error) {
	billData := make([]config.ShareBill, 0)

	if err := mysql.Engine.Where("month = ?", month).Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return nil, err
	}
	return billData, nil
}

// 计算预测数据

func ComputerBudget() {
	sourceData := make([]config.SourceBill, 0)

	if err := mysql.Engine.Find(&sourceData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
	}
}
