package budget

type Budget struct {
	Name string
	Type string
}

type BaiduBillData struct {
	AccountId           string  `json:"accountId" xorm:"VARCHAR(32) 'accountId' comment('账单所属账号ID')"`
	Amount              string  `json:"amount" xorm:"VARCHAR(32) 'amount' comment('用量')"`
	AmountUnit          string  `json:"amountUnit" xorm:"VARCHAR(32) 'amountUnit' comment('用量单位')"`
	Cash                float64 `json:"cash" xorm:"Float 'cash' comment('现金支付金额')"`
	ChargeItem          string  `json:"chargeItem" xorm:"VARCHAR(32) 'chargeItem' comment('后付费计费项英文名')"`
	ChargeItemDesc      string  `json:"chargeItemDesc" xorm:"VARCHAR(32) 'chargeItemDesc' comment('后付费计费项中文名')"`
	ConfigurationCH     string  `json:"configurationCH" xorm:"VARCHAR(1024) 'configurationCH' comment('产品中文配置')"`
	CouponPrice         float64 `json:"couponPrice" xorm:"Float 'couponPrice' comment('代金券金额')"`
	CreditCost          float64 `json:"creditCost" xorm:"Float 'creditCost' comment('账期未还金额')"`
	CreditRefund        float64 `json:"creditRefund" xorm:"Float 'creditRefund' comment('账期退款金额')"`
	Debt                float64 `json:"debt" xorm:"Float 'debt' comment('欠款金额')"`
	DiscountCouponPrice float64 `json:"discountCouponPrice" xorm:"Float 'discountCouponPrice' comment('折扣券金额')"`
	DiscountPrice       float64 `json:"discountPrice" xorm:"Float 'discountPrice' comment('折扣金额')"`
	DiscountUnit        string  `json:"discountUnit" xorm:"VARCHAR(32) 'discountUnit' comment('折扣率')"`
	Duration            string  `json:"duration" xorm:"VARCHAR(32) 'duration' comment('预付费服务时长')"`
	FinancePrice        float64 `json:"financePrice" xorm:"Float 'financePrice' comment('应付金额')"`
	InstanceId          string  `json:"instanceId" xorm:"VARCHAR(32) 'instanceId' comment('资源ID')"`
	NoPaidPrice         float64 `json:"noPaidPrice" xorm:"Float 'noPaidPrice' comment('优惠金额')"`
	OrderId             string  `json:"orderId" xorm:"VARCHAR(32) 'orderId' comment('资源购买的订单id')"`
	OrderPurchaseTime   string  `json:"orderPurchaseTime" xorm:"VARCHAR(32) 'orderPurchaseTime' comment('订单的支付时间，utc时间')"`
	OrderType           string  `json:"orderType" xorm:"VARCHAR(32) 'orderType' comment('订单类型')"`
	OrderTypeDesc       string  `json:"orderTypeDesc" xorm:"VARCHAR(32) 'orderTypeDesc' comment('订单类型中文')"`
	OriginPrice         float64 `json:"originPrice" xorm:"Float 'originPrice' comment('账单金额')"`
	PricingUnit         string  `json:"pricingUnit" xorm:"VARCHAR(32) 'pricingUnit' comment('价格单位')"`
	ProductType         string  `json:"productType" xorm:"VARCHAR(32) 'productType' comment('计费类型')"`
	Rebate              float64 `json:"rebate" xorm:"Float 'rebate' comment('返点支付金额')"`
	Region              string  `json:"region" xorm:"VARCHAR(32) 'region' comment('区域')"`
	ServiceType         string  `json:"serviceType" xorm:"VARCHAR(32) 'serviceType' comment('产品名')"`
	ServiceTypeName     string  `json:"serviceTypeName" xorm:"VARCHAR(32) 'serviceTypeName' comment('产品名中文')"`
	StartTime           string  `json:"startTime" xorm:"VARCHAR(32) 'startTime' comment('开始时间')"`
	SysGold             float64 `json:"sysGold" xorm:"Float 'sysGold' comment('消账金额')"`
	Tag                 string  `json:"tag" xorm:"VARCHAR(1024) 'tag' comment('tag')"`
	Tex                 float64 `json:"tex" xorm:"Float 'tex' comment('折扣率')"`
	UnitPrice           string  `json:"unitPrice" xorm:"VARCHAR(1024) 'unitPrice' comment('单价')"`
	Vendor              string  `json:"vendor" xorm:"VARCHAR(32) 'vendor' comment('供应商')"`
}

type BaiduBill struct {
	AccountId    string          `json:"accountId"`
	BillMonth    string          `json:"billMonth"`
	LoginName    string          `json:"loginName"`
	OuName       string          `json:"ouName"`
	SubAccountId string          `json:"subAccountId"`
	SubLoginName string          `json:"subLoginName"`
	PageNo       int             `json:"pageNo"`
	PageSize     int             `json:"pageSize"`
	TotalCount   int             `json:"totalCount"`
	Bills        []BaiduBillData `json:"bills"`
}
