package journal

type ItemMem struct {
	amount      float64
	description string
	document    string
}

type Item interface {
	GetAmount() float64
	SetAmount(val float64)
	GetDescription() string
	SetDescription(val string)
	GetDocument() string
	SetDocument(val string)
	Equal(val *ItemMem) bool
}

func NewItem() *ItemMem {
	return &ItemMem{}
}

func (i *ItemMem) GetAmount() float64 {
	return i.amount
}

func (i *ItemMem) SetAmount(val float64) {
	if val >= 0 {
		i.amount = val
	} else {
		i.amount = 0
	}
}

func (i *ItemMem) GetDescription() string {
	return i.description
}

func (i *ItemMem) SetDescription(val string) {
	i.description = val
}

func (i *ItemMem) GetDocument() string {
	return i.document
}

func (i *ItemMem) SetDocument(val string) {
	i.document = val
}

func (i *ItemMem) Equal(val *ItemMem) bool {
	return i.description == val.description
}
