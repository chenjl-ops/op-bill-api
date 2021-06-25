package prediction

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/thoas/go-funk"
	"op-bill-api/internal/app/billing"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/pkg/baiducloud"
	"op-bill-api/internal/pkg/mysql"
	"time"
)

// 获取账单接口数据
//func GetBaiduBillData() ([]BaiduBillData, error){
//	billData := make([]BaiduBillData, 0)
//	if err := mysql.Engine.Find(&billData); err != nil {
//		logger.Log.Error("查询数据异常: ", err)
//		return nil, err
//	}
//	return billData, nil
//}

const (
	timeFormat = "2006-01-02"
)

// GetQueryBaiduBillData 获取特定数据集合
func GetQueryBaiduBillData(serviceType string, sellType string) ([]BaiduBillData, error) {
	billData := make([]BaiduBillData, 0)
	if err := mysql.Engine.Where("serviceTypeName = ?", serviceType).And("productType = ?", sellType).Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return nil, err
	}
	return billData, nil
}

// GetAllServiceAndName 获取账单数据所有资源类别和资源名称集合
func GetAllServiceAndName() (map[string][]string, error) {
	data := make(map[string][]string)

	sellTypes := [2]string{"prepay", "postpay"}
	for _, sellType := range sellTypes {
		data[sellType] = make([]string, 0)
		namesData := make([]string, 0)
		//typesData := make([]string, 0)
		billData := make([]BaiduBillData, 0)
		if err := mysql.Engine.Where("productType = ?", sellType).Find(&billData); err != nil {
			return data, err
		}

		for _, data := range billData {
			if !funk.Contains(namesData, data.ServiceTypeName) {
				namesData = append(namesData, data.ServiceTypeName)
			}
			//if !funk.Contains(typesData, data.ServiceType) {
			//	typesData = append(typesData, data.ServiceType)
			//}
		}
		data[sellType] = namesData
	}
	return data, nil
}

func GetBaiduBillEveryDayData() error {
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
		logrus.Println("bill Data: ", billData)
		logrus.Println("Bill Data 入库异常: ", err)
		return err
	}
	return nil
}

// InsertPrediction 入库预测数据
func InsertPrediction(data map[string]map[string]map[string]float64) error {
	x := PredData{
		Date: time.Now().Format(timeFormat),
		Data: data,
	}
	_, err := mysql.Engine.Insert(&x)
	return err
}

// GetPrediction 查询预测数据
func GetPrediction(date string) (PredData, bool) {
	var data PredData
	has, err := mysql.Engine.Exist(&PredData{Date: date})
	logrus.Println("校验预测数据是否存在: ", has, err)
	if has {
		_, err := mysql.Engine.Where("date = ?", date).Get(&data)
		if err != nil {
			return data, false
		} else {
			return data, true
		}
	} else {
		return data, false
	}
}