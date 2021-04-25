package billing

type Getter interface {
	GetMonthBilling(month int, year int)
	GetDepartmentMonthBilling(name string, month int, year int)
	GetUserMonthBilling(name string, month int, year int)
}


