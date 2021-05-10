package compute

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/pkg/baiducloud"
	"strconv"
)



func GetBilling(c *gin.Context) {
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
				//"data": data,
				"cost":      cost,
				"nonCost":   nonCost,
				"otherCost": otherCost,
				"allCost":   allCost,
			})
		}
	}
}

func GetBudget(c *gin.Context) {
	data, err := ComputerBudget()
	if err != nil {
		c.JSON(500, gin.H{
			"msg": err,
		})
	} else {
		c.JSON(200, gin.H{
			"msg": "success",
			"data": data,
		})
	}
}

func TestBaiduBill(c *gin.Context) {
	url := "https://billing.baidubce.com/v1/bill/resource/month?month=%s&productType=%s&pageNo=%s&pageSize=%s"
	headers := map[string]string{
		"Host":         "billing.baidubce.com",
		"Content-Type": "application/json",
	}

	bc := baiducloud.NewBaiduCloud(fmt.Sprintf(url, "2021-03", "prepay", "1", "100"), headers, "GET")

	data := map[string]interface{}{}
	result := new(interface{})
	err := bc.Request(data, &result)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"msg":       "success",
			"data":      result,
			"Signature": bc.GetAuthorization(),
		})
	}

}

