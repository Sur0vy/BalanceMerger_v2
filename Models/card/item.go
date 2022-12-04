package card

import "time"

type ItemMem struct {
	document string
	date     time.Time
	position string
	in       int
	out      int
	comment  string
}

type Item interface {
	GetIn() int
	SetIn(val int)
	GetOut() int
	SetOut(val int)
	GetNeed() int
	GetPosition() string
	SetPosition(val string)
	GetDocument() string
	SetDocument(val string)
	GetDate() time.Time
	SetDate(val time.Time)
	GetComment() string
	SetComment(val string)
}

func NewItem() *ItemMem {
	return &ItemMem{}
}

func (i *ItemMem) GetIn() int {
	return i.in
}

func (i *ItemMem) SetIn(val int) {
	if val >= 0 {
		i.in = val
	} else {
		i.in = 0
	}
}

func (i *ItemMem) GetOut() int {
	return i.out
}

func (i *ItemMem) SetOut(val int) {
	if val >= 0 {
		i.out = val
	} else {
		i.out = 0
	}
}

func (i *ItemMem) GetPosition() string {
	return i.position
}

func (i *ItemMem) SetPosition(val string) {
	i.position = val
}

func (i *ItemMem) GetDocument() string {
	return i.document
}

func (i *ItemMem) SetDocument(val string) {
	i.document = val
}

func (i *ItemMem) GetNeed() int {
	return i.GetOut() - i.GetIn()
}

func (i *ItemMem) GetDate() time.Time {
	return i.date
}

func (i *ItemMem) SetDate(val time.Time) {
	i.date = val
}

func (i *ItemMem) GetComment() string {
	return i.comment
}

func (i *ItemMem) SetComment(val string) {
	i.comment = val
}
