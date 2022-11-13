package journal

type ItemMem struct {
	rest        float64
	description string
	document    string
}

type Item interface {
	GetRest() float64
	SetRest(val float64)
	GetDescription() string
	SetDescription(val string)
	GetDocument() string
	SetDocument(val string)
	Equal(val *ItemMem) bool
}

func NewItem() *ItemMem {
	return &ItemMem{}
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
