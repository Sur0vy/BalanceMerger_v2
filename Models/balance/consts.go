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

type MatchStatusColor string

const (
	accMatch     MatchStatusColor = "ffffff" //совпадение 85%
	accAlmost    MatchStatusColor = "c3f7c3" //совпадение 70%-85%
	accMayBe     MatchStatusColor = "dfeb65" //совпадение 50%-70%
	accHardly    MatchStatusColor = "e57a15" //совпадение 20%-50%
	accDifferent MatchStatusColor = "e51515" //совпадение <20%
)
