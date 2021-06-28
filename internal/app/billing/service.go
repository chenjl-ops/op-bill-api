package billing

import (
	"github.com/Sirupsen/logrus"
	"op-bill-api/internal/pkg/mysql"
	"time"
)

// GetMonthFirstDate 获取月第一天和最后一天
// 获取月第一天 shift为偏移量 正为下月 负为上月
func GetMonthFirstDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day()+1)
	return t.Format(timeFormat)

}

// GetMonthLastDate 获取月最后一天 shift为偏移量 正为下月 负为上月
func GetMonthLastDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day())
	return t.Format(timeFormat)
}

// GetMonthDate 获取月第一天和最后一天
func GetMonthDate() map[string]string {
	data := make(map[string]string)

	t := time.Now()
	data["thisMonthFirstDate"] = GetMonthFirstDate(t, 0)  // 本月第一天
	data["thisMonthLastDate"] = GetMonthLastDate(t, 1)    // 本月最后一天
	data["lastMonthFirstDate"] = GetMonthFirstDate(t, -1) //上月第一天
	data["lastMonthLastDate"] = GetMonthLastDate(t, 0)    //上月最后一天

	return data
}

// GetTexData 获取折扣率数据
func GetTexData(name string) (SourceBillTex, error){
	var data SourceBillTex
	_, err := mysql.Engine.Where("name = ?", name).Get(&data)
	return data, err
}

// 更新折扣率，通过name和tex action标识动作
func InsertOrUpdateTexData(name string, tex float64, action string) bool {
	data := SourceBillTex{Name: name, Tex: tex}
	var err error
	if action == "POST" { // 支持新增和update减少客户端逻辑
		has, _ := mysql.Engine.Exist(&SourceBillTex{Name:name})
		if has {
			_, err = mysql.Engine.Update(&data, &SourceBillTex{Name: name})
		} else {
			_, err = mysql.Engine.Insert(&data)
		}
	} else if action == "PUT" {
		_, err = mysql.Engine.Update(&data, &SourceBillTex{Name: name})

	} else {
		return false
	}

	// 判断操作结果
	if err != nil {
		return false
	} else {
		return true
	}
}

// CheckBillData 校验账单数据
func CheckBillData(month string, isShare bool) bool {
	has, _ := mysql.Engine.Exist(&BillData{
		Month:   month,
		IsShare: isShare,
	})
	return has
}

// GetBillData 查询账单数据
func GetBillData(month string, isShare bool) (BillData, error) {
	var data BillData
	_, err := mysql.Engine.Where("month = ?", month).And("isShare = ?", isShare).Get(&data)
	return data, err
}

// InsertBillData 数据写入
func InsertBillData(month string, isShare bool, data map[string]float64) error {
	x := BillData{
		Month: month,
		IsShare: isShare,
		Data: data,
	}

	_, err := mysql.Engine.Insert(&x)
	return err
}

// GetAllBillData 获取账单全量数据 资金/损益 口径
// TODO 分页功能
func GetAllBillData(isShare bool) ([]BillData, error) {
	var data []BillData
	err := mysql.Engine.Where("isShare = ?", isShare).Find(&data)
	if err != nil {
		logrus.Println("获取账单全量数据异常: ", err)
	}
	return data, err
}
