package billing

// BillData 存放决算数据
type BillData struct {
	Month   string             `json:"month" xorm:"VARCHAR(32) 'month' comment('决算月份')"`
	IsShare bool               `json:"isShare" xorm:"INT(2) 'isShare' comment('是否是分摊花费')"`
	Data    map[string]map[string]float64 `json:"data" xorm:"TEXT 'data' comment('数据')"`
}

type SourceBillTex struct {
	Name string  `json:"name" xorm:"VARCHAR(128) unique 'name' comment('资源名称')"`
	Tex  float64 `json:"tex" xorm:"Float 'tex' comment('折扣率')"`
}

type BillDataResponse struct {
	Month     string  `json:"month"`
	IsShare   bool    `json:"isShare"`
	AllCost   float64 `json:"allCost"`
	Cost      float64 `json:"cost"`
	NoneCost  float64 `json:"noneCost"`
	OtherCost float64 `json:"otherCost"`
}
