package prediction

type Getter interface {
	GetMonthBudget(month int, year int)
	GetDepartmentMonthBudget(name string, month int, year int)
	GetUserMonthBudget(name string, month int, year int)
}

type Setter interface {
	Set(id string, name string, month int, year int)
}

type Deleter interface {
	Delete(id string)
}

type Inserter interface {
	InsertBudget(name string, month int, year int)
}
