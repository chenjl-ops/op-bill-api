package compute

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/thoas/go-funk"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
	"strconv"
)

func GetBilling(c *gin.Context) {
	month := c.DefaultQuery("month", "")
	isShare, _ := strconv.ParseBool(c.Query("isShare"))
	if month == "" {
		c.JSON(500, gin.H{
			"msg": "parameter month not none",
		})
	} else {
		cost, nonCost, otherCost, allCost, err := ComputerBilling(month, isShare)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": err,
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "success",
				//"data": data,
				"cost":      cost,
				"nonCost":   nonCost,
				"otherCost": otherCost,
				"allCost":   allCost,
			})
		}
	}

}

// 计算决算数据
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

	// 处理bcc ip和instanceid对应关系
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

// 计算预测数据
func ComputerBudget() {
	sourceData := make([]config.SourceBill, 0)

	if err := mysql.Engine.Find(&sourceData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
	}
}
