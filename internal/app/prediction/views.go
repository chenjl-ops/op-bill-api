package prediction

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/app/billing"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/baiducloud"
	"op-bill-api/internal/pkg/mysql"
	"time"
)

func getBaiduBillEveryDayData() error {
	dateData := billing.GetMonthDate()
	// 默认偏移量为1天
	shift := 1
	now := time.Now()
	// 为了保证数据计算，百度建议10点以后 10点之前取两天前数据
	if now.Hour() < 10 {
		shift = 2
	}

	startTime := dateData["thisMonthFirstDate"]
	tShift, _ := time.ParseDuration(fmt.Sprintf("%dh", -shift*24))
	EndTime := now.Add(tShift)
	endTime := EndTime.Format("2006-01-02")

	url := "https://billing.baidubce.com/v1/bill/resource/month?beginTime=%s&endTime=%s&productType=%s&pageNo=%d&pageSize=%d"
	headers := map[string]string{
		"Host":         "billing.baidubce.com",
		"Content-Type": "application/json",
	}
	// 计费类型：prepay/ postpay，分别表示预付费/后付费
	productType := [2]string{"prepay", "postpay"}
	// productType := [1]string{"prepay"}
	pageSize := 100

	billData := make([]BaiduBillData, 0)

	for _, v := range productType {
		logrus.Println("数据处理开始: ", fmt.Sprintf(url, startTime, endTime, v, 1, pageSize))
		bc := baiducloud.NewBaiduCloud(fmt.Sprintf(url, startTime, endTime, v, 1, pageSize), headers, "GET")
		data := map[string]interface{}{}
		var resultData BaiduBill
		err := bc.Request(data, &resultData)

		logrus.Println("数据处理中... ")
		if err == nil {
			nu := resultData.TotalCount/pageSize + 1
			billCh := make(chan BaiduBill)
			for i := 0; i <= nu; i++ {
				go getBaiduBillData(startTime, endTime, v, i, pageSize, billCh)
			}
			for i := 0; i <= nu; i++ {
				responseData := <-billCh
				billData = append(billData, responseData.Bills...)
			}
		}
	}

	// 入库
	err := insertBillData(billData)
	if err != nil {
		logrus.Println("Bill Data 入库异常: ", err)
		return err
	}
	return nil
}

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

func createTable(c *gin.Context) {
	err := mysql.Engine.Sync2(new(BaiduBillData))
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

func insertBaiduBillData(c *gin.Context) {
	err := getBaiduBillEveryDayData()
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
