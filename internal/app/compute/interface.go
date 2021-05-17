package compute


type Getter interface {
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
	ComputerBilling(month string, isShare bool) (c float64, nc float64, oc float64, ac float64, err error)
}

type Prediction interface {
	ComputerPrediction() (map[string]map[string]map[string]float64, error)
}
