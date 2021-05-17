package compute

import (
	"fmt"
	"github.com/Sirupsen/logrus"
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

// ComputerBilling 计算决算数据
func ComputerBilling(month string, isShare bool) (c float64, nc float64, oc float64, ac float64, err error) {
	// 获取数据库所有配置
	logrus.Println("isShare", isShare)

	billData := make([]config.ShareBill, 0)
	if err := mysql.Engine.Where("month = ?", month).Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return 0.00, 0.00, 0.00, 0.00, err
	}

	sourceData := make([]config.SourceBill, 0)
	if err := mysql.Engine.Where("month = ?", month).Find(&sourceData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return 0.00, 0.00, 0.00, 0.00, err
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
		for _, instance := range appInstances.Data {
			if instance.Env == "prod" {
				bccIdData[bccIdData[instance.Ip]] = instance.AppName
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
			}

		}
	} else {
		logrus.Println("source 计算: ")
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
			}
		}

	}

	return cost, nonCost, otherCost, allCost, nil
}

// ComputerPrediction 计算预测数据 后付费相关数据
func ComputerPrediction() (map[string]map[string]map[string]float64, error) {
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

