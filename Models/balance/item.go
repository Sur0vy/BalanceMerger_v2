package balance

import "BM/Models"

type ItemMem struct {
	rest        float64
	description string
	document    string
	name        string
	bill        string
	count       int
	comment     string
	state       Models.ItemState
}

type Item interface {
	GetRest() float64
	SetRest(val float64)
	GetDescription() string
	SetDescription(val string)
	GetDocument() string
	SetDocument(val string)
	GetName() string
	SetName(val string)
	GetBill() string
	SetBill(val string)
	GetComment() string
	SetComment(val string)
	GetCount() int
	SetCount(val int)
	GetState() Models.ItemState
	SetState(val Models.ItemState)
	Equal(val *ItemMem) bool
}

func (i *ItemMem) Equal(val *ItemMem) bool {
	if val == nil {
		return false
	}
	return i.name == val.name
}

func (i *ItemMem) statusToStr() string {
	switch i.state {
	case Models.IsFound:
		return "Успешно"
	case Models.IsCollect:
		return "Объединённая строка"
	case Models.IsMissing:
		return "Нет данных"
	case Models.IsCollectMissing:
		return "Несколько совпадений"
	default:
		return "Не совпадает остаток"
	}
}

func (i *ItemMem) statusToColor() string {
	switch i.state {
	case Models.IsFound:
		return "#90EE90"
	case Models.IsCollect:
		return "#32CD32"
	case Models.IsMissing:
		return "#FF0033"
	case Models.IsCollectMissing:
		return "#FFB6C1"
	default:
		return "#FFFF00"
	}
}
