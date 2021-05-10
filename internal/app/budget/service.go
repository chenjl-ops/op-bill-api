package budget

import (
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/pkg/mysql"
)

// 获取账单接口数据
func GetBaiduBillData() ([]BaiduBillData, error){
	billData := make([]BaiduBillData, 0)
	if err := mysql.Engine.Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return nil, err
	}
	return billData, nil
}

// 获取特定数据集合
func GetQueryBaiduBillData(query string) ([]BaiduBillData, error) {
	billData := make([]BaiduBillData, 0)
	if err := mysql.Engine.Where("serviceType = ?", query).Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return nil, err
	}
	return billData, nil
}