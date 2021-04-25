package compute

type Calculate interface {
	CalculateBilling(month int, year int)
	CalculateBudget(month int, year int)
}

type Setter interface {
	Set()
}

type Deleter interface {
	Delete(id string)
}

type Inserter interface {
	InsertCalculate()
}
