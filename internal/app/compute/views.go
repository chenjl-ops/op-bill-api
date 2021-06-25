package compute

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"op-bill-api/internal/app/billing"
	"op-bill-api/internal/app/prediction"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"strconv"
	"time"
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

// 获取折扣率数据
func getTexData(name string) (billing.SourceBillTex, error) {
	var texData billing.SourceBillTex
	_, err := mysql.Engine.Where("name = ?", name).Get(&texData)
	return texData, err

}

// @Tags Compute API
// @Summary Select billing data
// @Description 查询决算数据
// @Accept  application/json
// @Produce  application/json
// @Param month query string true "get bill of month"
// @Param isShare query boolean true "get bill of share or source"
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /bill/v1/get_bill_data [get]
func getBilling(c *gin.Context) {
	month := c.DefaultQuery("month", "")
	isShare, _ := strconv.ParseBool(c.Query("isShare"))
	if month == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "parameter month not none",
		})
	} else {
		if billing.CheckBillData(month, isShare) {
			data, err := billing.GetBillData(month, isShare)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "success",
					"data": data.Data,
				})
			}
		} else {
			data, err := CalculateBilling(month, isShare)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "success",
					"data": data,
				})
			}
		}
	}
}

// @Tags Compute API
// @Summary Select prediction data
// @Description 查询预测数据
// @Accept  application/json
// @Produce  application/json
// @Param date query string false "get prediction of date, default today"
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /bill/v1/get_prediction_data [get]
func getPrediction(c *gin.Context) {
	date := c.DefaultQuery("date", time.Now().Format(timeFormat))
	data, has := prediction.GetPrediction(date)
	if !has {
		err := prediction.GetBaiduBillEveryDayData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": err,
			})
		} else {
			// data, err := CalculatePrediction()
			data, err := CalculatePredictionV2()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": err,
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"msg":  "success",
					"data": data,
				})
			}
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"data": data.Data,
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
