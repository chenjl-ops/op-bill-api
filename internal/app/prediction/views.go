package prediction

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/baiducloud"
	"op-bill-api/internal/pkg/mysql"
)

// 获取本月每天 自定义page和pageSize数据
func getBaiduBillData(startTime string, endTime string, productType string, pageNum int, pageSize int, ch chan<- BaiduBill) {
	url := "https://billing.baidubce.com/v1/bill/resource/month?beginTime=%s&endTime=%s&productType=%s&pageNo=%d&pageSize=%d"
	headers := map[string]string{
		"Host":         "billing.baidubce.com",
		"Content-Type": "application/json",
	}

	// 计费类型：prepay/ postpay，分别表示预付费/后付费
	bc := baiducloud.NewBaiduCloud(fmt.Sprintf(url, startTime, endTime, productType, pageNum, pageSize), headers, "GET")
	data := map[string]interface{}{}
	var resultData BaiduBill
	err := bc.Request(data, &resultData)
	if err == nil {
		ch <- resultData
	}
}

func insertBillData(billData []BaiduBillData) error {
	// 清空表数据
	truncateSql := "truncate table baidu_bill_data"
	_, err := mysql.Engine.QueryString(truncateSql)
	if err != nil {
		return err
	}

	// 写入账单数据，由于账单数据无法表示唯一记录，所以需要删除数据再录入
	for i := 0; i <= len(billData); i = i + apollo.Config.InsertMysqlSum {
		insertData := make([]BaiduBillData, 0)
		if i+apollo.Config.InsertMysqlSum >= len(billData) {
			insertData = billData[i:]
		} else {
			insertData = billData[i : i+apollo.Config.InsertMysqlSum]
		}
		_, err := mysql.Engine.Insert(&insertData)
		if err != nil {
			logrus.Println("bill数据写入异常: ", err)
			return err
		}
	}
	return nil
}

// @Tags Prediction API
// @Summary Create BillData Table
// @Description create BillData table
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /prediction/v1/create_table [get]
func createTable(c *gin.Context) {
	err := mysql.Engine.Sync2(new(BaiduBillData), new(PredData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed",
			"error": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	}
}

// @Tags Prediction API
// @Summary Insert Bill Data
// @Description Insert Bill Data
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /prediction/v1/insert_baidu_bill_data [get]
func insertBaiduBillData(c *gin.Context) {
	err := GetBaiduBillEveryDayData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed",
			"error": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	}
}
