package card

type ItemMem struct {
	document string
	position string
	in       int
	out      int
	//isAdded  bool
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
	//GetIsAdded() bool
	//SetIsAdded(val bool)
}

func NewItem() *ItemMem {
	return &ItemMem{
		//isAdded: false,
	}
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

//func (i *ItemMem) SetIsAdded(val bool) {
//	i.isAdded = val
//}
//
//func (i *ItemMem) GetIsAdded() bool {
//	return i.isAdded
//}
