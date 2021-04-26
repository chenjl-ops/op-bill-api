package compute

import (
	"github.com/gin-gonic/gin"
	"op-bill-api/internal/app/middleware/logger"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/config"
	"op-bill-api/internal/pkg/mysql"
)

func GetBilling(c *gin.Context) {
	month := c.DefaultQuery("month", "")
	if month == "" {
		c.JSON(500, gin.H{
			"msg": "parameter month not none",
		})
	} else {
		data, err := ComputerBilling(month)
		if err != nil {
			c.JSON(500, gin.H{
				"msg": err,
			})
		} else {
			c.JSON(200, gin.H{
				"msg": "success",
				//"data": data,
				"length": len(data),
			})
		}
	}

}

// 计算决算数据
func ComputerBilling(month string) ([]config.ShareBill, error) {
	// 获取数据库所有配置
	billData := make([]config.ShareBill, 0)
	if err := mysql.Engine.Where("month = ?", month).Find(&billData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
		return nil, err
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
		go GetAppInstanceData(app.Name, instanceCh)
	}

	var instances []InstanceData
	for range allApps.Data {
		// 全局instances instances = append(instances, appInstances.Data...)
		// 循环instance prod环境加入列表
		appInstances := <-instanceCh
		for _, instance := range appInstances.Data {
			if instance.Env == "prod" {
				instances = append(instances, instance)
			}
		}
	}

	var infoData map[string][]string
	var appInfoData map[string]map[string][]string

	for _, instance := range instances {
		// 循环ecs主机
		for _, ecs := range ecss.Data {
			if instance.Ip == ecs.Ip {
				infoData["ecs"] = append(infoData["ecs"], ecs.InstanceId)
			}
			// 循环磁盘
			for _, volume := range allVolumes.Data.Volumes {
				for _, v := range volume.Attachments {
					if v.InstanceId == ecs.InstanceId {
						infoData["volume"] = append(infoData["volume"], v.VolumeId)
					}
				}
			}
		}
		appInfoData[instance.AppName] = infoData
	}

	// 获取所有前相关应用数据
	var dependentApp []string
	for _, dependent := range dependents.Data {
		if dependent.DependentName == apollo.Config.DependentName {
			dependentApp = append(dependentApp, dependent.AppName)
		}
	}

	// 获取所有相关性数据

	return billData, nil
}

// 计算预测数据
func ComputerBudget() {
	sourceData := make([]config.SourceBill, 0)

	if err := mysql.Engine.Find(&sourceData); err != nil {
		logger.Log.Error("查询数据异常: ", err)
	}
}
