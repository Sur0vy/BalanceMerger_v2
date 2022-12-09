package balance

import (
	"time"

	"BM/Models"
)

type ItemMem struct {
	rest        float64
	description string
	document    string
	name        string
	bill        string
	count       int
	comment     string
	spent       int
	parent      *ItemMem
	date        time.Time
	state       Models.ItemState
	position    string
	accuracy    float64
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
	GetSpent() int
	SetSpent(val int)
	GetDate() time.Time
	SetDate(val time.Time)
	GetState() Models.ItemState
	SetState(val Models.ItemState)
	Equal(val *ItemMem) bool
	GetPosition() string
	SetPosition(val string)
	GetAccuracy() float64
	SetAccuracy(val float64)
}

func NewItem() *ItemMem {
	return &ItemMem{}
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
	return ColorFromStatus(i.state)
}

func (i *ItemMem) GetDescription() string {
	return i.description
}

func (i *ItemMem) SetDescription(val string) {
	i.description = val
}

func (i *ItemMem) GetName() string {
	return i.name
}

func (i *ItemMem) SetName(val string) {
	i.name = val
}

func (i *ItemMem) GetCount() int {
	return i.count
}

func (i *ItemMem) SetCount(val int) {
	if val >= 0 {
		i.count = val
	} else {
		i.count = 0
	}
}

func (i *ItemMem) GetRest() float64 {
	return i.rest
}

func (i *ItemMem) SetRest(val float64) {
	if val >= 0 {
		i.rest = val
	} else {
		i.rest = 0
	}
}

func (i *ItemMem) GetSpent() int {
	return i.spent
}

func (i *ItemMem) SetSpent(val int) {
	if val >= 0 {
		i.spent = val
	} else {
		i.spent = 0
	}
}

func (i *ItemMem) GetAccuracy() float64 {
	return i.accuracy
}

func (i *ItemMem) SetAccuracy(val float64) {
	val *= 100
	if val >= 0 {
		i.accuracy = val
	} else {
		i.accuracy = 0
	}
}

func (i *ItemMem) GetBill() string {
	return i.bill
}

func (i *ItemMem) SetBill(val string) {
	i.bill = val
}

func (i *ItemMem) GetState() Models.ItemState {
	return i.state
}

func (i *ItemMem) SetState(val Models.ItemState) {
	i.state = val
}

func (i *ItemMem) GetDocument() string {
	return i.document
}

func (i *ItemMem) SetDocument(val string) {
	i.document = val
}

func (i *ItemMem) GetComment() string {
	return i.comment
}

func (i *ItemMem) SetComment(val string) {
	i.comment = val
}

func (i *ItemMem) GetDate() time.Time {
	return i.date
}

func (i *ItemMem) SetDate(val time.Time) {
	i.date = val
}

func (i *ItemMem) GetParent() *ItemMem {
	return i.parent
}

func (i *ItemMem) SetParent(val *ItemMem) {
	i.parent = val
}

func (i *ItemMem) GetPosition() string {
	return i.position
}

func (i *ItemMem) SetPosition(val string) {
	i.position = val
}

func ColorFromStatus(status Models.ItemState) string {
	switch status {
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

func ColorFromMatchPers(val float64) string {
	if val < 15 {
		return string(accMatch)
	} else if val < 30 {
		return string(accAlmost)
	} else if val < 50 {
		return string(accMayBe)
	} else if val < 80 {
		return string(accHardly)
	} else {
		return string(accDifferent)
	}
}
