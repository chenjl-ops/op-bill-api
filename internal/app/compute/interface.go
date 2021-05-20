package compute

type CMDBService interface {
	GetAllApplications() Apps
	GetEcsData() Ecs
	GetVolumeData() Volume
	GetAppInstanceData(appName string, ch chan<- Instance)
	GetAppInfoData() AppInfo
}

type Calculate interface {
	Bill
	Prediction
}

type Bill interface {
	CalculateBilling(month string, isShare bool) (c float64, nc float64, oc float64, ac float64, err error)
}

type Prediction interface {
	CalculatePrediction() (map[string]map[string]map[string]float64, error)
}
