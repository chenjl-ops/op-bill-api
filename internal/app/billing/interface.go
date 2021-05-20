package billing

// BillService 主要实现决算相关细分接口
type BillService interface {
	GetMonthBilling(month int, year int)
	GetDepartmentMonthBilling(name string, month int, year int)
	GetUserMonthBilling(name string, month int, year int)
}
