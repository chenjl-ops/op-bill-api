package config

/*
ShareBill billing package中结构体数据，由于模块化处理，避免循环调用问题，抽离为上层结构
*/
type ShareBill struct {
	Month         string  `json:"month" xorm:"VARCHAR(32) 'month' comment('分摊月')"`
	AccountId     string  `json:"accountId" xorm:"VARCHAR(64) 'account_id' comment('账号ID')"`
	LoginName     string  `json:"loginName" xorm:"VARCHAR(32) 'login_name' comment('登录名')"`
	Unit          string  `json:"unit" xorm:"VARCHAR(32) 'unit' comment('单元')"`
	AssetId       string  `json:"assetId" xorm:"VARCHAR(64) 'asset_id' comment('资源/量包id')"`
	ProductName   string  `json:"productName" xorm:"VARCHAR(32) 'product_name' comment('产品名')"`
	Zone          string  `json:"zone" xorm:"VARCHAR(32) 'zone' comment('区域')"`
	Configuration string  `json:"configuration" xorm:"VARCHAR(1024) 'configuration' comment('配置')"`
	OrderId       string  `json:"orderId" xorm:"VARCHAR(64) 'order_id' comment('账单ID')"`
	OrderTime     string  `json:"orderTime" xorm:"VARCHAR(64) 'order_time' comment('下单时间')"`
	OrderType     string  `json:"orderType" xorm:"VARCHAR(16) 'order_type' comment('账单类型')"`
	StartTime     string  `json:"startTime" xorm:"VARCHAR(64) 'start_time' comment('服务开始时间')"`
	EndTime       string  `json:"endTime" xorm:"VARCHAR(64) 'end_time' comment('服务结束时间')"`
	ServiceTime   string  `json:"serviceTime" xorm:"VARCHAR(16) 'service_time' comment('服务时长')"`
	SourceBill    string  `json:"sourceBill" xorm:"VARCHAR(32) 'source_bill' comment('原始账单')"`
	SourceCope    string  `json:"sourceCope" xorm:"VARCHAR(32) 'source_cope' comment('原始应付')"`
	SourceVoucher string  `json:"sourceVoucher" xorm:"VARCHAR(16) 'source_voucher' comment('原始代金劵')"`
	SourceSale    string  `json:"sourceSale" xorm:"VARCHAR(32) 'source_sale' comment('原始折扣')"`
	ShareTime     string  `json:"shareTime" xorm:"VARCHAR(16) 'share_time' comment('分摊时长')"`
	ShareBill     string  `json:"shareBill" xorm:"VARCHAR(32) 'share_bill' comment('分摊账单')"`
	ShareCope     string  `json:"shareCope" xorm:"VARCHAR(32) 'share_cope' comment('分摊应付')"`
	ShareVoucher  float32 `json:"shareVoucher" xorm:"FLOAT 'share_voucher' comment('分摊代金劵')"`
	ShareSale     string  `json:"shareSale" xorm:"VARCHAR(32) 'share_sale' comment('分摊折扣')"`
}

type SourceBill struct {
	Month       string `json:"month" xorm:"VARCHAR(32) 'month' comment('分摊月')"`
	AccountId   string `json:"accountId" xorm:"VARCHAR(64) 'account_id' comment('账号ID')"`
	LoginName   string `json:"loginName" xorm:"VARCHAR(32) 'login_name' comment('登录名')"`
	Unit        string `json:"unit" xorm:"VARCHAR(32) 'unit' comment('单元')"`
	AssetId     string `json:"assetId" xorm:"VARCHAR(64) 'asset_id' comment('资源/量包id')"`
	ProductName string `json:"productName" xorm:"VARCHAR(32) 'product_name' comment('产品名')"`
	Zone        string `json:"zone" xorm:"VARCHAR(32) 'zone' comment('区域')"`
	OrderTime   string `json:"orderTime" xorm:"VARCHAR(64) 'order_time' comment('下单时间')"`
	SellType    string `json:"orderType" xorm:"VARCHAR(16) 'sell_type' comment('售卖方式')"`
	OrderCost   string `json:"orderCost" xorm:"VARCHAR(64) 'order_cost' comment('账单金额')"`
}

type BillStatus struct {
	FileName string `json:"filename" xorm:"VARCHAR(128) 'filename' comment('账单文件名称')"`
	Status   bool   `json:"status" xorm:"TINYINT 'status' comment('状态')"`
}

type ResponseData struct {
	Msg   string                 `json:"msg"`
	Error string                 `json:"error"`
	Data  map[string]interface{} `json:"data"`
}
