package billing

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"os"
	"strings"
)

const (
	timeFormat = "2006-01-02"
)

// 获取账单
func getBillExcel() error {
	monthDateData := GetMonthDate()

	// 账单地址
	shareBillUrl := fmt.Sprintf(apollo.Config.ShareBillUrl, monthDateData["lastMonthFirstDate"])
	sourceBillUrl := fmt.Sprintf(apollo.Config.SourceBillUrl, monthDateData["lastMonthFirstDate"], monthDateData["lastMonthLastDate"])

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

	fileFullName := fmt.Sprintf(apollo.Config.DownloadPath, filename)
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
	monthDateData := GetMonthDate()
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
		logrus.Println("账单文件开始写入: ", filename)
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
		logrus.Println("账单文件写入完成: ", filename)
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



// 创建数据表
// @Tags Billing API
// @Summary Create Table
// @Description 创建损益和资金口径账单数据表，对应账单状态表
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/create_table [get]
func createTable(c *gin.Context) {
	err := mysql.Engine.Sync2(new(config.ShareBill), new(config.SourceBill), new(config.BillStatus), new(BillData), new(SourceBillTex))
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

// 初始化折扣率数据
// @Tags Billing API
// @Summary Create Table
// @Description 初始化折扣率数据
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/init_tex_data [get]
func initTexData(c *gin.Context) {
	// source账单折扣率
	var sourceBillTex = map[string]float64{
		"Elasticsearch":   0.5,
		"DDoS高防IP ADAS":   0.8,
		"负载均衡 BLB":        0.8,
		"对等连接 PEERCONN":   0.8,
		"NAT网关":           0.8,
		"内容分发网络 CDN":      1,
		"海外CDN CDN_ABOAD": 1,
		"视频创作分发平台":        1,
		"音视频直播 LSS":       1,
		"对象存储 BOS":        1,
	}

	var data []SourceBillTex
	for name, tex := range sourceBillTex {
		var dataTemp SourceBillTex
		temp := make(map[string]interface{})
		temp["name"] = name
		temp["tex"] = tex

		if err := mapstructure.Decode(temp, &dataTemp); err != nil {
			logrus.Error("bill数据转换异常: ", err)
		}
		data = append(data, dataTemp)
	}
	_, err := mysql.Engine.Insert(&data)
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


// Update资源折扣率
// @Tags Billing API
// @Summary Get Data
// @Description 更新资源折扣率
// @Accept  application/json
// @Produce  application/json
// @Param tex body SourceBillTex true "new SourceBillTex"
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/tex [put]
func updateTexData(c *gin.Context) {
	var json SourceBillTex
	c.BindJSON(&json)

	is := InsertOrUpdateTexData(json.Name, json.Tex, "PUT")
	if is {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed",
		})
	}
}

// 新增资源折扣率
// @Tags Billing API
// @Summary Get Data
// @Description 新增资源折扣率
// @Accept  application/json
// @Produce  application/json
// @Param tex body SourceBillTex true "new SourceBillTex"
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/tex [post]
func insertTexData(c *gin.Context) {
	/*
	json 数据类型
	//json := make(map[string]interface{})
	{"name": "xxxx", "tex": xx}
	 */

	var json SourceBillTex
	c.BindJSON(&json)

	is := InsertOrUpdateTexData(json.Name, json.Tex, "POST")
	if is {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed",
		})
	}
}


// 删除资源折扣率
// @Tags Billing API
// @Summary Get Data
// @Description 删除资源折扣率
// @Accept  application/json
// @Produce  application/json
// @Param name query string false "delete name of tex"
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/tex [delete]
func deleteTexData(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	if name == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":   "failed",
			"error": "name error, name must not none",
		})
	} else {
		_, err := mysql.Engine.Where("name = ?", name).Delete(&SourceBillTex{Name: name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "failed",
				"error": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":     "success",
			})
		}
	}
}

// 获取资源折扣率
// @Tags Billing API
// @Summary Get Data
// @Description 获取资源折扣率
// @Accept  application/json
// @Produce  application/json
// @Param name query string false "select name of tex"
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/tex [get]
func getTexData(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	if name != "" {
		data, err := GetTexData(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "failed",
				"error": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":     "success",
				"data":    data,
				"columns": SourceBillTexColumns,
			})
		}
	} else {
		var data []SourceBillTex
		err := mysql.Engine.Find(&data)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg":   "failed",
				"error": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg":     "success",
				"data":    data,
				"columns": SourceBillTexColumns,
			})
		}
	}

}

// 账单数据录入
// @Tags Billing API
// @Summary Insert Data
// @Description 插入账单数据 资金和损益口径
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/insert_bill_data [get]
func insertData(c *gin.Context) {
	err := getBillExcel()
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

// 查看月首尾日期
// @Tags Billing API
// @Summary Select Month Data
// @Description 插入账单数据 资金和损益口径
// @Accept  application/json
// @Produce  application/json
// @Success 200 {object} config.ResponseData
// @Header 200 {object}  config.ResponseData
// @Failure 400,404 {object} string "Bad Request"
// @Router /billing/v1/get_month_data [get]
func getMonthData(c *gin.Context) {
	dateData := GetMonthDate()
	c.JSON(http.StatusOK, gin.H{
		"data": dateData,
	})
}

// 自定义获取账单周期
// TODO 1、后续可由apollo配置周期
//func JobService() error {
//	c := cron.New()
//	_, err := c.AddFunc("@monthly", getBillExcel)
//	if err != nil {
//		return err
//	}
//	c.Start()
//	return nil
//}
