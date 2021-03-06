package billing

var billDict = map[string]string{
	"分摊月":      "Month",
	"账号id":     "AccountId",
	"登录名":      "LoginName",
	"所属单元":     "Unit",
	"资源/量包id":  "AssetId",
	"产品名称":     "ProductName",
	"区域":       "Zone",
	"配置":       "Configuration",
	"订单号":      "OrderId",
	"下单时间":     "OrderTime",
	"订单类型":     "OrderType",
	"服务开始时间":   "StartTime",
	"服务结束时间":   "EndTime",
	"服务时长":     "ServiceTime",
	"原始账单":     "SourceBill",
	"原始应付":     "SourceCope",
	"原始代金券":    "SourceVoucher",
	"原始折扣":     "SourceSale",
	"本月分摊时长":   "ShareTime",
	"本月分摊账单":   "ShareBill",
	"本月分摊应付":   "ShareCope",
	"本月分摊代金劵":  "ShareVoucher",
	"本月分摊折扣金额": "ShareSale",
}

var sourceDict = map[string]string{
	"分摊月":   "Month",
	"账号ID":  "AccountId",
	"账号登录名": "LoginName",
	"单元":    "Unit",
	"实例ID":  "AssetId",
	"产品名":   "ProductName",
	"区域":    "Zone",
	"账单时间":  "OrderTime",
	"售卖方式":  "SellType",
	"账单金额":  "OrderCost",
}

var SourceBillTexColumns = []map[string]string{
	{"name": "name", "nickName": "资源名称"},
	{"name": "tex", "nickName": "折扣率"},
}
