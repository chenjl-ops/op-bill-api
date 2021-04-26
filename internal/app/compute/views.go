package compute

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"op-bill-api/internal/pkg/apollo"
	"op-bill-api/internal/pkg/requests"
)

// 请求各种接口获取需要接口数据

// 请求CMDB获取所有APP
func GetAllApplications() Apps {
	var data Apps
	err := requests.Request(apollo.Config.CmdbAppUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// 获取CMDB ecs 数据集合 参数可控
func GetEcsData() Ecs {
	var data Ecs
	err := requests.Request(apollo.Config.CmdbAppUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// 获取volume数据结合接口
func GetVolumeData() Volume {
	var data Volume
	err := requests.Request(apollo.Config.CmdbVolumeUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}

// 获取应用对应instance数据接口
func GetAppInstanceData(appName string, ch chan<- Instance) {
	var data Instance
	url := fmt.Sprintf(apollo.Config.CmdbAppInstanceUrl, appName)

	err := requests.Request(url, &data)
	if err != nil {
		logrus.Println(err)
	}
	ch <- data
}

// 获取所有appInfo数据  应用强相关数据
func GetAppInfoData() AppInfo {
	var data AppInfo
	err := requests.Request(apollo.Config.AppInfoUrl, &data)
	if err != nil {
		logrus.Println(err)
	}
	return data
}
