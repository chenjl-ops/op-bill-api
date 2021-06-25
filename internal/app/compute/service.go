package compute

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"op-bill-api/internal/app/billing"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/app/prediction"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"op-bill-api/internal/pkg/requests"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

// 请求各种接口获取需要接口数据

// GetAllApplications 请求CMDB获取所有APP
func GetAllApplications() Apps {
	var data Apps
	err := requests.Request(apollo.Config.CmdbAppUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// GetEcsData 获取CMDB ecs 数据集合 参数可控
func GetEcsData() Ecs {
	var data Ecs
	err := requests.Request(apollo.Config.CmdbAppUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// GetVolumeData 获取volume数据结合接口
func GetVolumeData() Volume {
	var data Volume
	err := requests.Request(apollo.Config.CmdbVolumeUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// GetAppInstanceData 获取应用对应instance数据接口
func GetAppInstanceData(appName string, ch chan<- Instance) {
	var data Instance
	url := fmt.Sprintf(apollo.Config.CmdbAppInstanceUrl, appName)

	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
	}
	ch <- data
}

// GetAppInfoData 获取所有appInfo数据  应用强相关数据
func GetAppInfoData() AppInfo {
	var data AppInfo
	err := requests.Request(apollo.Config.AppInfoUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// CalculateBilling 计算决算数据
func CalculateBilling(month string, isShare bool) (map[string]float64, error) {
	// 获取数据库所有配置
	logrus.Println("isShare", isShare)
	otherBillData := make(map[string]float64, 0)

	billData := make([]config.ShareBill, 0)
	if err := mysql.Engine.Where("month = ?", month).Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return otherBillData, err
	}

	sourceData := make([]config.SourceBill, 0)
	if err := mysql.Engine.Where("month = ?", month).Find(&sourceData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return otherBillData, err
	}

	// 获取所有应用
	allApps := GetAllApplications()
	// 获取所有磁盘
	allVolumes := GetVolumeData()
	// 获取所有ECS数据
	ecss := GetEcsData()
	// 获取相关性数据
	dependents := GetAppInfoData()

	// 获取所有应用对应instance prod环境

	// 并发请求所有应用实例信息
	instanceCh := make(chan Instance)

	// 控制并发数
	//var wg sync.WaitGroup
	//limitChan := make(chan int, 100)
	//for i := 0; i< len(allApps.Data); i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		for d := range limitChan {
	//			app := allApps.Data[d]
	//			GetAppInstanceData(app.Name, instanceCh)
	//		}
	//	}()
	//}
	//
	//for i := 0; i< len(allApps.Data); i++ {
	//	limitChan <- 1
	//	limitChan <- 2
	//}
	//
	//close(limitChan)
	//wg.Wait()

	for _, app := range allApps.Data {
		// TODO 并发请求数据造成后端接口无法响应，后续增加限流措施
		go GetAppInstanceData(app.Name, instanceCh)
	}

	logrus.Println("数据处理开始: ")

	// 处理bcc ip和instanceId对应关系
	bccIdIpData := make(map[string]string)
	for _, ecs := range ecss.Data {
		if ecs.Status == "Running" && ecs.Env == "Prod" {
			bccIdIpData[ecs.Ip] = ecs.InstanceId
		}

	}

	// 定义bcc id 对应appName
	bccIdData := make(map[string]string)
	for range allApps.Data {
		// 全局instances instances = append(instances, appInstances.Data...)
		// 循环instance prod环境加入列表
		appInstances := <-instanceCh
		if len(appInstances.Data) > 0 {
			for _, instance := range appInstances.Data {
				if instance.Env == "prod" {
					bccIdData[bccIdData[instance.Ip]] = instance.AppName
				}
			}
		}
	}

	// 处理磁盘
	volumeIdData := make(map[string]string)

	for _, volume := range allVolumes.Data.Volumes {
		for _, v := range volume.Attachments {
			volumeIdData[v.VolumeId] = bccIdData[v.InstanceId]
		}
	}

	logrus.Println("数据处理完成: ")

	// 获取所有前相关应用数据
	var dependentApp []string
	for _, dependent := range dependents.Data {
		if dependent.DependentName == apollo.Config.DependentName {
			dependentApp = append(dependentApp, dependent.AppName)
		}
	}

	// 获取所有相关性数据
	cost := 0.00
	nonCost := 0.00
	otherCost := 0.00
	allCost := 0.00

	if isShare {
		// TODO share账单需要新增source账单后付费账单信息
		for _, v := range shareBillProductNames {
			otherBillData[v] = 0.00
		}
		for _, v := range billData {
			// 计算prod 强相关数据
			if v.ProductName == "云服务器 BCC" {
				x, err := strconv.ParseFloat(v.ShareCope, 64)
				if err == nil {
					allCost = allCost + x
				}
				if funk.Contains(dependentApp, bccIdData[v.AssetId]) {
					if err == nil {
						cost = cost + x
					}
				} else {
					if err == nil {
						nonCost = nonCost + x
					}
				}
			} else if v.ProductName == "云磁盘 CDS" {
				x, err := strconv.ParseFloat(v.ShareCope, 64)
				if err == nil {
					allCost = allCost + x
				}
				if funk.Contains(dependentApp, volumeIdData[v.AssetId]) {
					if err == nil {
						cost = cost + x
					}
				} else {
					if err == nil {
						nonCost = nonCost + x
					}
				}
			} else {
				x, err := strconv.ParseFloat(v.ShareCope, 64)
				if err == nil {
					allCost = allCost + x
					otherCost = otherCost + x
				}

				if _, ok := otherBillData[v.ProductName]; ok {
					otherBillData[v.ProductName] = otherBillData[v.ProductName] + x
				}
			}

		}
	} else {
		logrus.Println("source 计算: ")
		for _, v := range sourceBillProductNames {
			otherBillData[v] = 0.00
		}
		for _, v := range sourceData {
			// 计算prod 强相关数据
			if v.ProductName == "云服务器 BCC" {
				x, err := strconv.ParseFloat(v.OrderCost, 64)
				if err == nil {
					allCost = allCost + x
				}
				if funk.Contains(dependentApp, bccIdData[v.AssetId]) {
					if err == nil {
						cost = cost + x
					}
				} else {
					if err == nil {
						nonCost = nonCost + x
					}
				}
			} else if v.ProductName == "云磁盘 CDS" {
				x, err := strconv.ParseFloat(v.OrderCost, 64)
				if err == nil {
					allCost = allCost + x
				}
				if funk.Contains(dependentApp, volumeIdData[v.AssetId]) {
					if err == nil {
						cost = cost + x
					}
				} else {
					if err == nil {
						nonCost = nonCost + x
					}
				}
			} else {
				x, err := strconv.ParseFloat(v.OrderCost, 64)
				if err == nil {
					allCost = allCost + x
					otherCost = otherCost + x
				}
				if _, ok := otherBillData[v.ProductName]; ok {
					otherBillData[v.ProductName] = otherBillData[v.ProductName] + x
				}
			}
		}

	}

	resultData := make(map[string]float64)
	resultData = otherBillData
	//logrus.Println("COST: ", cost, nonCost, otherCost, allCost)
	resultData["cost"] = cost
	resultData["nonCost"] = nonCost
	resultData["otherCost"] = otherCost
	resultData["allCost"] = allCost

	// 计算数据 * 折扣率
	if !isShare {
		for k, v := range resultData {
			// 获取数据库折扣率
			texData, err := getTexData(k)
			if err != nil {
				vv, _ := decimal.NewFromFloat(v * texData.Tex).Round(2).Float64()
				resultData[k] = vv
			} else {
				// 默认3.9折
				vv, _ := decimal.NewFromFloat(v * 0.39).Round(2).Float64()
				resultData[k] = vv
			}
		}
	}

	// 入库操作
	if !billing.CheckBillData(month, isShare) {
		err := billing.InsertBillData(month, isShare, resultData)
		if err != nil {
			logrus.Println("决算数据入库异常: ", err)
		}
	}

	return resultData, nil
}

// CalculatePrediction 计算预测数据 后付费相关数据
func CalculatePrediction() (map[string]map[string]map[string]float64, error) {
	logrus.Println("开始...")
	dateData := billing.GetMonthDate()

	shift := 1
	now := time.Now()
	// 为了保证数据计算，百度建议10点以后 10点之前取两天前数据
	if now.Hour() < 10 {
		shift = 2
	}

	thisMonthL := strings.Split(dateData["thisMonthLastDate"], "-")
	thisMonthTotal, err := strconv.Atoi(thisMonthL[len(thisMonthL)-1])
	if err != nil {
		logrus.Println("获取月最后一天数据转换失败: ", err)
	}
	tShift, _ := time.ParseDuration(fmt.Sprintf("%dh", -shift*24))
	EndTime := now.Add(tShift)

	logrus.Println("日期: ", thisMonthTotal, EndTime.Day())

	billDataMap := map[string]map[string]string{
		"视频创作分发平台":      {"name": "VOD", "type": "source"},
		"内容分发网络 CDN":    {"name": "CDN", "type": "source"},
		"函数计算 CFC":      {"name": "CFC", "type": "source"},
		"云数据库RDS":       {"name": "RDS", "type": "share"},
		"对象存储 BOS":      {"name": "BOS", "type": "source"},
		"Elasticsearch": {"name": "BES", "type": "source"},
	}

	thisMonthBudgetData := make(map[string]map[string]map[string]float64)

	sellTypes := [2]string{"prepay", "postpay"}

	for _, sellType := range sellTypes {
		thisMonthBudgetData[sellType] = make(map[string]map[string]float64)
		for k, v := range billDataMap {
			logrus.Println("处理中...", k, v)

			// 获取目前消费总金额
			billData, err := prediction.GetQueryBaiduBillData(v["name"], sellType)
			if err != nil {
				return nil, err
			}

			// 当月已产生金额总和
			financePriceTotal := 0.00
			for _, bill := range billData {
				financePriceTotal = financePriceTotal + bill.FinancePrice
			}
			//logrus.Println("当月产生金额总和: ", financePriceTotal)

			// 查询上个月消费总和
			lastMonthCostSum := 0.00
			sourceMonth := fmt.Sprintf("%s_%s", dateData["lastMonthFirstDate"], dateData["lastMonthLastDate"])

			// 当月应付金额总和 和 新增金额
			thisMonthFinancePriceTotal := 0.00
			thisMonthLastMonthAdd := 0.00

			if v["type"] == "source" {
				thisMonthFinancePriceTotal = financePriceTotal / (float64(EndTime.Day()) / float64(thisMonthTotal))

				lastMonthSourceData, err := getLastMonthSourceCost(sourceMonth, k, sellType)
				if err != nil {
					return nil, err
				}
				for _, v := range lastMonthSourceData {
					v1, err := strconv.ParseFloat(v.OrderCost, 64)
					if err != nil {
						return nil, err
					}
					lastMonthCostSum = lastMonthCostSum + v1
				}
				//logrus.Println("上个月消费总和: ", lastMonthCostSum)

				thisMonthLastMonthAdd = thisMonthFinancePriceTotal - lastMonthCostSum
			} else {
				thisMonthFinancePriceTotal = financePriceTotal / (float64(EndTime.Hour()) / float64(thisMonthTotal))

				lastShareData, err := getLastMonthShareCost(strings.Replace(dateData["lastMonthFirstDate"], "-", "/", -1), k)
				if err != nil {
					return nil, err
				}
				for _, v := range lastShareData {
					v1, err := strconv.ParseFloat(v.ShareCope, 64)
					if err != nil {
						return nil, err
					}
					lastMonthCostSum = lastMonthCostSum + v1
				}
				thisMonthLastMonthAdd = thisMonthFinancePriceTotal - lastMonthCostSum
			}
			fmt.Println(thisMonthFinancePriceTotal, thisMonthLastMonthAdd)
			thisMonthBudgetData[sellType][k] = make(map[string]float64)
			thisMonthBudgetData[sellType][k]["total"] = thisMonthFinancePriceTotal
			thisMonthBudgetData[sellType][k]["add"] = thisMonthLastMonthAdd

		}
	}

	logrus.Println("处理完成...", thisMonthBudgetData)
	return thisMonthBudgetData, nil
}

// 预测数据计算V2
func CalculatePredictionV2() (map[string]map[string]map[string]float64, error) {
	//logrus.Println("开始...")
	sellTypes := [2]string{"prepay", "postpay"} // 初始化账单纬度 预付费 后付费

	dateData := billing.GetMonthDate()

	shift := 1
	now := time.Now()
	// 为了保证数据计算，百度建议10点以后 10点之前取两天前数据
	if now.Hour() < 10 {
		shift = 2
	}

	thisMonthL := strings.Split(dateData["thisMonthLastDate"], "-")    // 转换月最后一天日期为数组 例如: 2021-06-30 -> [2021,06,30]
	thisMonthTotal, err := strconv.Atoi(thisMonthL[len(thisMonthL)-1]) // 获取最后一天日期 30 并转换成int
	if err != nil {
		logrus.Println("获取月最后一天数据转换失败: ", err)
	}
	tShift, _ := time.ParseDuration(fmt.Sprintf("%dh", -shift*24)) // 日期按偏移量 shift * 24小时做时间回归
	EndTime := now.Add(tShift)                                     // 获取账单数据的时间日期，例如: 今天是 2021-06-23 10:00 以后，获取的2021-06-22之前账单
	//logrus.Println("日期: ", thisMonthTotal, EndTime.Day())

	thisMonthPredictionData := make(map[string]map[string]map[string]float64) // 定义预测数据存放结果map

	// 获取账单所有名称和名称类别
	namesData, err := prediction.GetAllServiceAndName()

	for _, sellType := range sellTypes {
		thisMonthPredictionData[sellType] = make(map[string]map[string]float64)

		for _, v := range namesData[sellType] {
			// 获取baidu_bill_data内所有账单数据
			billData, err := prediction.GetQueryBaiduBillData(v, sellType)
			if err != nil {
				return nil, err
			}

			// 当月已产生金额总和
			financePriceTotal := 0.00

			// 后付费直接按花费总和计算
			if sellType == "postpay" {
				for _, bill := range billData {
					tempFinancePriceTotal := financePriceTotal + bill.FinancePrice/(float64(EndTime.Day())/float64(thisMonthTotal))
					tempFinancePriceTotal, _ = decimal.NewFromFloat(tempFinancePriceTotal).Round(2).Float64()
					financePriceTotal = tempFinancePriceTotal
				}
			} else {
				// 预付费需要按分摊数据计算月花费
				for _, bill := range billData {
					// 获取非续费数据
					if bill.OrderType != "RENEW" {
						// 获取应付金额
						tempFinancePrice := bill.FinancePrice
						// 获取分摊时间 xx 天|xx 月|xx 年|/
						duration := bill.Duration
						// 设置偏移天数init
						initShiftDay := 0
						if duration != "/" {
							tempDuration := strings.Split(duration, " ")
							initShiftDay, err = strconv.Atoi(tempDuration[0])
							if err != nil {
								logrus.Println("获取月最后一天数据转换失败: ", err)
							}

							if tempDuration[1] == "年" {
								// 年计费周期按 365天一年计算
								initShiftDay = initShiftDay * 365
							} else if tempDuration[1] == "月" {
								// 计费周期为月 判断 是否为12个月 12个月为 365天周期，其他一律按 月*30 计算天周期
								if tempDuration[0] == "12" {
									initShiftDay = 365
								} else {
									initShiftDay = initShiftDay * 30
								}
							}
						}
						// 计算每天单价
						tempDayUnitPrice := tempFinancePrice / float64(initShiftDay)

						// 获取订单开始时间  2021-06-21T08:39:42Z
						startDay := strings.Split(strings.Split(bill.StartTime, "T")[0], "-")
						// 计算开始时间到月底一共多少天
						tempStartDay, err := strconv.Atoi(startDay[len(startDay)-1])
						if err != nil {
							logrus.Println("获取账单开启时间异常: ", err)
						}
						// 计算资源开始时间到月底总天数
						orderTimeRange := thisMonthTotal - tempStartDay
						// 计算当月资源花费
						tempFinancePriceTotal := financePriceTotal + (tempDayUnitPrice * float64(orderTimeRange))
						tempFinancePriceTotal, _ = decimal.NewFromFloat(tempFinancePriceTotal).Round(2).Float64()
						financePriceTotal = tempFinancePriceTotal
					}
				}
			}
			//logrus.Println("当月产生金额总和: ", v, financePriceTotal)

			// 查询上个月消费总和
			lastMonthCostSum := 0.00
			sourceMonth := fmt.Sprintf("%s_%s", dateData["lastMonthFirstDate"], dateData["lastMonthLastDate"])

			thisMonthPredictionData[sellType][v] = make(map[string]float64)
			if sellType == "postpay" {
				lastMonthSourceData, err := getLastMonthSourceCost(sourceMonth, apiNameBillNameMap[v], sellType)
				if err != nil {
					return nil, err
				}
				for _, v := range lastMonthSourceData {
					v1, err := strconv.ParseFloat(v.OrderCost, 64)
					if err != nil {
						return nil, err
					}
					v1, _ = decimal.NewFromFloat(v1).Round(2).Float64()
					lastMonthCostSum = lastMonthCostSum + v1
				}

				// 计算上个月花费总和折扣点
				texData, err := getTexData(apiNameBillNameMap[v])
				if err != nil {
					lastMonthCostSum = lastMonthCostSum * texData.Tex
				} else {
					// 默认3.9折
					lastMonthCostSum = lastMonthCostSum * 0.39
				}

				financePriceTotal, _ = decimal.NewFromFloat(financePriceTotal).Round(2).Float64()
				lastMonthCostSum, _ = decimal.NewFromFloat(lastMonthCostSum).Round(2).Float64()

				thisMonthPredictionData[sellType][v]["Total"] = financePriceTotal
				thisMonthPredictionData[sellType][v]["LastMonthCost"] = lastMonthCostSum
				thisMonthPredictionData[sellType][v]["Add"] = financePriceTotal - lastMonthCostSum
			}
			if sellType == "prepay" {
				//lastShareData, err := getLastMonthShareCost(strings.Replace(dateData["lastMonthFirstDate"], "-", "/", -1), apiNameBillNameMap[v])
				//if err != nil {
				//	return nil, err
				//}
				//for _, v := range lastShareData {
				//	v1, err := strconv.ParseFloat(v.ShareCope, 64)
				//	if err != nil {
				//		return nil, err
				//	}
				//	v1, _ = decimal.NewFromFloat(v1).Round(2).Float64()
				//	lastMonthCostSum = lastMonthCostSum + v1
				//}
				financePriceTotal, _ = decimal.NewFromFloat(financePriceTotal).Round(2).Float64()
				thisMonthPredictionData[sellType][v]["Total"] = financePriceTotal

				// message.NewPrinter(language.English).Sprintln(financePriceTotal) // 千位数打印
			}
		}
	}

	// 预测数据入库
	err = prediction.InsertPrediction(thisMonthPredictionData)
	if err != nil {
		logrus.Println("预测数据入库失败: ", err)
	}
	return thisMonthPredictionData, nil
}


