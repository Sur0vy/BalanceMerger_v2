package balance

type BalanceState int

const (
	IsEmpty BalanceState = iota
	IsMergeV1
	IsMergeV2
)

const (
	fldBill   string = "Счет"
	fldName   string = "Артикул"
	fldRest   string = "Остаток"
	fldCount  string = "Количество"
	fldPerEnd string = " на окончание периода"
	fldDesc   string = "Номенклатура"

	tryCnt = 10
)
