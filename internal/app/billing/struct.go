package billing

type BillData struct {
	Month     string  `json:"month" xorm:"VARCHAR(32) 'month' comment('决算月份')"`
	Cost      float64 `json:"cost" xorm:"Float 'accountId' comment('用车成本花费')"`
	NonCost   float64 `json:"nonCost" xorm:"Float 'nonCost' comment('非用车成本花费')"`
	OtherCost float64 `json:"otherCost" xorm:"Float 'otherCost' comment('其他花费')"`
	AllCost   float64 `json:"allCost" xorm:"Float 'allCost' comment('所有花费')"`
	IsShare   bool    `json:"isShare" xorm:"INT(2) 'isShare' comment('是否是分摊花费')"`
}
