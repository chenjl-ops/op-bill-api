package billing

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Sirupsen/logrus"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"os"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

// 获取月第一天和最后一天
// 获取月第一天 shift为偏移量 正为下月 负为上月
func getMonthFirstDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day()+1)
	return t.Format(timeFormat)

}

// 获取月最后一天 shift为偏移量 正为下月 负为上月
func getMonthLastDate(t time.Time, shift int) string {
	t = t.AddDate(0, shift, -t.Day())
	return t.Format(timeFormat)
}

// 获取月第一天和最后一天
func getMonthDate() map[string]string {
	data := make(map[string]string)

	t := time.Now()
	data["thisMonthFirstDate"] = getMonthFirstDate(t, 0)  // 本月第一天
	data["thisMonthLastDate"] = getMonthLastDate(t, 1)    // 本月最后一天
	data["lastMonthFirstDate"] = getMonthFirstDate(t, -1) //上月第一天
	data["lastMonthLastDate"] = getMonthLastDate(t, 0)    //上月最后一天

	return data
}

// 获取账单
func getBillExcel() error {
	monthDateData := getMonthDate()

	// 账单地址
	// TODO 未来可以apollo管理账单地址，非硬编码
	shareBillUrl := fmt.Sprintf("https://test-public-cn.bj.bcebos.com/bill/%s_share_bill.xlsx", monthDateData["lastMonthFirstDate"])
	sourceBillUrl := fmt.Sprintf("https://test-public-cn.bj.bcebos.com/bill/%s_%s_resource_bill.xlsx", monthDateData["lastMonthFirstDate"], monthDateData["lastMonthLastDate"])

	// 获取月 相关文件名称数据
	shareFileArray := strings.Split(shareBillUrl, "/")
	sourceFileArray := strings.Split(sourceBillUrl, "/")
	shareFileName := shareFileArray[len(shareFileArray)-1]
	sourceFileName := sourceFileArray[len(sourceFileArray)-1]

	// 下载费用和资金账单数据
	var filenames []string

	filenames = append(filenames, downloadFile(shareBillUrl, shareFileName))
	filenames = append(filenames, downloadFile(sourceBillUrl, sourceFileName))

	for _, fileName := range filenames {
		if strings.Contains(fileName, "share") {
			excelData, err1 := readExcel(fileName, true)
			if err1 != nil {
				logrus.Error("获取账单数据失败: ", fileName, err1)
				return err1
			}
			err := insertBillData(excelData, true, fileName)
			if err != nil {
				logrus.Error(err)
				return err
			}
		} else {
			excelData, err1 := readExcel(fileName, false)
			if err1 != nil {
				logrus.Error("获取账单数据失败: ", fileName, err1)
				return err1
			}
			err := insertBillData(excelData, false, fileName)
			if err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	return nil
	// compute.ComputerBilling()
}

// 下载文件
func downloadFile(url string, filename string) string {
	resp, err := http.Get(url)

	if err != nil {
		logrus.Error("获取boss账单地址excel失败: ", err)
	}
	defer resp.Body.Close()

	fileFullName := fmt.Sprintf("/tmp/%s", filename)
	out, err := os.Create(fileFullName)
	if err != nil {
		logrus.Error("创建本地目录文件失败: ", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logrus.Error("写入excel内容失败: ", err)
	}
	return fileFullName
}

// 读取excel文件获取相关数据
func readExcel(filename string, isShare bool) (excelData []map[string]string, err error) {
	monthDateData := getMonthDate()
	// fmt.Println(filename)
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}

	var sheetname string
	if isShare {
		sheetname = fmt.Sprintf("工作表 1 - %s_share_bill", monthDateData["lastMonthFirstDate"])
	} else {
		sheetname = fmt.Sprintf("工作表 1 - %s_%s_r", monthDateData["lastMonthFirstDate"], monthDateData["lastMonthLastDate"])
	}

	// 获取所有单元格
	rows := f.GetRows(sheetname)

	// 获取第一行
	index := rows[1]
	var Data []map[string]string

	for _, row := range rows[1:] {
		data := make(map[string]string)
		for i, colCell := range row {
			//fmt.Print(colCell, "\t")
			if isShare {
				data[billDict[index[i]]] = colCell
			} else {
				data[sourceDict[index[i]]] = colCell
			}
		}
		// 补充资金口径数据月份信息
		if !isShare {
			data["Month"] = fmt.Sprintf("%s_%s", monthDateData["lastMonthFirstDate"], monthDateData["lastMonthLastDate"])
		}
		Data = append(Data, data)
		// fmt.Println()
	}
	return Data, nil
}

// 录入账单数据
func insertBillData(data []map[string]string, isShare bool, filename string) error {
	// 判断账单记录是否已经录入
	if checkInsertMonthData(filename) {
		logrus.Println("账单文件已经写入: ", filename)
	} else {
		// 存储账单录入记录
		billStatus := config.BillStatus{FileName: filename, Status: true}
		_, err := mysql.Engine.Insert(&billStatus)
		if err != nil {
			logrus.Println("bill记录状态异常: ", err)
			return err
		}
		if isShare {
			var billData []config.ShareBill
			for _, x := range data {
				var billDataTmp config.ShareBill
				if err := mapstructure.Decode(x, &billDataTmp); err != nil {
					logrus.Error("bill数据转换异常: ", err)
				}
				billData = append(billData, billDataTmp)
			}
			// 每100条数据批量写入一次数据库操作
			for i := 0; i <= len(billData); i = i + apollo.Config.InsertMysqlSum {
				insertData := make([]config.ShareBill, 0)
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
		} else {
			var sourceData []config.SourceBill
			for _, x := range data {
				var sourceDataTmp config.SourceBill
				if err := mapstructure.Decode(x, &sourceDataTmp); err != nil {
					logrus.Error("source数据转换异常: ", err)
				}
				sourceData = append(sourceData, sourceDataTmp)
			}
			for i := 0; i <= len(sourceData); i = i + apollo.Config.InsertMysqlSum {
				insertData := make([]config.SourceBill, 0)
				if i+apollo.Config.InsertMysqlSum >= len(sourceData) {
					insertData = sourceData[i:]
				} else {
					insertData = sourceData[i : i+apollo.Config.InsertMysqlSum]
				}
				_, err := mysql.Engine.Insert(&insertData)
				if err != nil {
					logrus.Println("source数据写入异常: ", err)
					return err
				}
			}
			//_, err := mysql.Engine.Insert(&sourceData)
			//if err != nil {
			//	return err
			//}
		}
	}
	return nil
}

// 判断账单数据是否已经录入
func checkInsertMonthData(filename string) bool {
	has, _ := mysql.Engine.Exist(&config.BillStatus{
		FileName: filename,
	})
	return has
}

// 自定义获取账单周期
// TODO 1、后续可有apollo配置周期
//func JobService() error {
//	c := cron.New()
//	_, err := c.AddFunc("@monthly", getBillExcel)
//	if err != nil {
//		return err
//	}
//	c.Start()
//	return nil
//}
