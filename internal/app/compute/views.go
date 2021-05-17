package compute

import (
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"strconv"
)

// 获取上个月总和消费数据 资金口径
func getLastMonthSourceCost(month string, query string, sellType string) ([]config.SourceBill, error) {
	sourceBillData := make([]config.SourceBill, 0)
	sellTypes := map[string]string{"prepay": "预付费", "postpay": "后付费"}

	if err := mysql.Engine.Where("product_name = ?", query).And("month = ?", month).And("sell_type = ?", sellTypes[sellType]).Find(&sourceBillData); err != nil {
		return nil, err
	}

	return sourceBillData, nil
}

// 获取上个月总和消费数据 分摊口径
func getLastMonthShareCost(month string, query string) ([]config.ShareBill, error) {
	shareBillData := make([]config.ShareBill, 0)

	if err := mysql.Engine.Where("product_name = ?", query).And("month = ?", month).Find(&shareBillData); err != nil {
		return nil, err
	}

	return shareBillData, nil
}

func getBilling(c *gin.Context) {
	month := c.DefaultQuery("month", "")
	isShare, _ := strconv.ParseBool(c.Query("isShare"))
	if month == "" {
		c.JSON(500, gin.H{
			"msg": "parameter month not none",
		})
	} else {
		cost, nonCost, otherCost, allCost, err := ComputerBilling(month, isShare)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": err,
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "success",
				"cost":      cost,
				"nonCost":   nonCost,
				"otherCost": otherCost,
				"allCost":   allCost,
			})
		}
	}
}

func getBudget(c *gin.Context) {
	data, err := ComputerPrediction()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": err,
		})
	} else {
		c.JSON(200, gin.H{
			"msg":  "success",
			"data": data,
		})
	}
}

// testBaiduBill 测试验证百度接口以及验签
//func testBaiduBill(c *gin.Context) {
//	url := "https://billing.baidubce.com/v1/bill/resource/month?month=%s&productType=%s&pageNo=%s&pageSize=%s"
//	headers := map[string]string{
//		"Host":         "billing.baidubce.com",
//		"Content-Type": "application/json",
//	}
//
//	bc := baiducloud.NewBaiduCloud(fmt.Sprintf(url, "2021-03", "prepay", "1", "100"), headers, "GET")
//
//	data := map[string]interface{}{}
//	result := new(interface{})
//	err := bc.Request(data, &result)
//	if err != nil {
//		c.JSON(500, gin.H{
//			"error": err,
//		})
//	} else {
//		c.JSON(200, gin.H{
//			"msg":       "success",
//			"data":      result,
//			"Signature": bc.GetAuthorization(),
//		})
//	}
//
//}
