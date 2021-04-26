package compute

import "time"

// 应用内data数据结构
type AppData struct {
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Owner    string `json:"Owner"`
	Users    string `json:"Users"`
}

// 应用数据结构
type Apps struct {
	Code    int       `json:"code"`
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    []AppData `json:"data"`
}

// ecs data数据结构
type EcsData struct {
	Env          string `json:"env"`
	Ip           string `json:"privateIp"`
	InstanceId   string `json:"instanceId"`
	InstanceName string `json:"instanceName"`
	Status       string `json:"status"`
}

// ecs数据
type Ecs struct {
	Code int       `json:"code"`
	Data []EcsData `json:"data"`
}

//// machine结构
//type Machine struct {
//	Code    int       `json:"code"`
//	Success bool      `json:"success"`
//	Data    []EcsData `json:"data"`
//}

// instance data
type InstanceData struct {
	Env     string `json:"env"`
	Ip      string `json:"ip"`
	AppName string `json:"appName"`
}

// instance
type Instance struct {
	Code    int            `json:"code"`
	Success bool           `json:"success"`
	Data    []InstanceData `json:"data"`
}

// Volume data 结构
type VolumeAttachments struct {
	VolumeId   string `json:"volumeId"`
	Serial     string `json:"serial"`
	InstanceId string `json:"instanceId"`
	Device     string `json:"device"`
}

type VolumeData struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	DiskSizeInGB   int                 `json:"diskSizeInGB"`
	CreateTime     time.Duration       `json:"createTime"`
	ExpireTime     time.Duration       `json:"expireTime"`
	Status         string              `json:"status"`
	Type           string              `json:"type"`
	StorageType    string              `json:"storageType"`
	IsSystemVolume bool                `json:"isSystemVolume"`
	Description    string              `json:"description"`
	PaymentTiming  string              `json:"paymentTiming"`
	ZoneName       string              `json:"zoneName"`
	Encrypted      bool                `json:"encrypted"`
	Attachments    []VolumeAttachments `json"attachments"`
}

type Volumes struct {
	Volumes []VolumeData `json:"volumes"`
}

type Volume struct {
	Code int     `json:"code"`
	Data Volumes `json:"data"`
}

// appInfo 应用相关性数据
type AppInfoData struct {
	DependentName string `json:"dependent_name"`
	AppName       string `json:"app_name"`
}

type AppInfo struct {
	Code int           `json:"code"`
	Data []AppInfoData `json:"data"`
}
