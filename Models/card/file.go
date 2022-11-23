package card

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type CardMem struct {
	fileName string
	itemsIn  map[int]*ItemMem
	itemsOut map[int]*ItemMem
}

type Card interface {
	GetItemsCount() int
	LoadFromFile(fileName string) error
	HasItemOut(doc string, position string) int
	GetItemIn(idx int) *ItemMem
	GetItemOut(idx int) *ItemMem

	simplifyDocument(val string) string
	simplifyPosition(val string) string
	findField(field string, row int, f *excelize.File) int
	findRow(field string, f *excelize.File) int
	fillOutItems()
}

func NewCard() *CardMem {
	return &CardMem{
		fileName: "",
		itemsIn:  make(map[int]*ItemMem),
		itemsOut: make(map[int]*ItemMem),
	}
}

func (c *CardMem) LoadFromFile(fileName string) error {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return errors.New("file is corrupted")
	}
	row := c.findRow(fldDoc, xlsx)
	if row == -1 {
		return errors.New("file read error")
	}
	iDoc := c.findField(fldDoc, row, xlsx)
	if iDoc == -1 {
		return errors.New("file read error")
	}
	iDocD := c.findField(fldDocD, row, xlsx)
	if iDocD == -1 {
		return errors.New("file read error")
	}
	iDocC := c.findField(fldDocC, row, xlsx)
	if iDocC == -1 {
		return errors.New("file read error")
	}
	iCntC := c.findField(fldCntC, row, xlsx)
	if iCntC == -1 {
		return errors.New("file read error")
	} else {
		iCntC++
	}
	iCntD := c.findField(fldCntD, row, xlsx)
	if iCntD == -1 {
		return errors.New("file read error")
	} else {
		iCntD++
	}
	i := 0
	row += 2
	for i < tryCnt {
		row += 2

		cell, _ := excelize.CoordinatesToCellName(iDoc, row)
		doc, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
		if doc == "" {
			i++
		} else {
			i = 0
			item := NewItem()
			item.SetDocument(c.simplifyDocument(doc))
			cell, _ = excelize.CoordinatesToCellName(iDocD, row)
			name, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)

			if err != nil {
				i++
				continue
			}
			var cntStr string
			var in bool
			if name == "" {
				i++
				continue
			} else if strings.Contains(name, repCredit) {
				//out
				in = false
				cell, _ := excelize.CoordinatesToCellName(iDocC, row)
				name, err = xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				if err != nil {
					i++
					continue
				}
				cell, _ = excelize.CoordinatesToCellName(iCntC, row+1)
				cntStr, err = xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				item.SetPosition(c.simplifyPosition(name))
				cnt, err := strconv.ParseFloat(cntStr, 64)
				if err != nil {
					i++
					continue
				}
				item.SetOut(int(cnt))
			} else {
				//in
				in = true
				cell, _ := excelize.CoordinatesToCellName(iCntD, row+1)
				cntStr, err = xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				item.SetPosition(c.simplifyPosition(name))
				cnt, err := strconv.ParseFloat(cntStr, 64)
				if err != nil {
					i++
					continue
				}
				item.SetIn(int(cnt))
			}
			if err != nil {
				i++
				continue
			}
			if in {
				c.itemsIn[c.GetItemsCount()] = item
			} else {
				c.itemsOut[c.GetItemsCount()] = item
			}
		}
	}
	if c.GetItemsCount() > 0 {
		c.fillOutItems()
		return nil
	} else {
		return errors.New("no items in file")
	}
}

func (c *CardMem) findField(field string, row int, f *excelize.File) int {
	i := 1
	col := 1
	value := ""
	for i < tryCnt {
		cell, _ := excelize.CoordinatesToCellName(col, row)
		value, _ = f.GetCellValue(f.GetSheetName(0), cell)
		value := strings.ReplaceAll(value, "\n", " ")
		if value == "" {
			i++
		} else if value == field {
			return col
		} else {
			i = 1
		}
		col++
	}
	return -1
}

func (c *CardMem) findRow(field string, f *excelize.File) int {
	iRow := 1
	row := 1
	value := ""
	for iRow < tryCnt {
		iCol := 1
		col := 1
		for iCol < tryCnt {
			cell, _ := excelize.CoordinatesToCellName(col, row)
			value, _ = f.GetCellValue(f.GetSheetName(0), cell)
			if value == "" {
				iCol++
			} else if value == field {
				return row
			} else {
				iCol = 1
			}
			col++
		}
		row++
		cell, _ := excelize.CoordinatesToCellName(1, row)
		value, _ = f.GetCellValue(f.GetSheetName(0), cell)
		if value == "" {
			iRow++
		} else {
			iRow = 1
		}
	}
	return -1
}

func (c *CardMem) GetItemsCount() int {
	return len(c.itemsIn) + len(c.itemsOut)
}

func (c *CardMem) fillOutItems() {
	for i, _ := range c.itemsOut {
		out := c.GetItemOut(i)
		if out.GetNeed() == 0 {
			continue
		}
		for j, _ := range c.itemsIn {
			in := c.GetItemIn(j)
			if in.GetIn() == 0 {
				continue
			}
			if out.position == in.position {
				if out.GetNeed() <= in.GetIn() {
					val := out.GetNeed()
					out.SetIn(out.GetOut() + val)
					in.SetIn(in.GetIn() - val)
					out.SetDocument(in.GetDocument())
					break
				} else {
					out.SetIn(out.GetIn() + in.GetIn())
					in.SetIn(0)
					if out.GetDocument() != "" {
						out.SetDocument(out.GetDocument() + ", ")
					}
					out.SetDocument(out.GetDocument() + in.GetDocument())
				}
			}
		}
	}
}

func (c *CardMem) simplifyDocument(val string) string {
	start := strings.Index(val, repDocB)
	if start == -1 {
		return ""
	}
	start += len(repDocB) + 1
	ret := val[start:]

	ret = ret[:strings.Index(ret, " ")]
	return ret
}

func (c *CardMem) simplifyPosition(val string) string {
	//нижний регистр
	ret := strings.ToLower(val)

	//удалим тему
	if strings.Contains(ret, "\n") {
		ret = ret[:strings.Index(ret, "\n")]
	}
	//проверим резистор, особенная обработка
	if strings.Contains(ret, "резистор") {
		start := strings.Index(ret, "(")
		stop := strings.Index(ret, ")")
		if start < stop {
			ret = ret[start+1 : stop]
		} else {
			ret = strings.ReplaceAll(ret, "резистор", "")
		}
	} else {
		//удалим реперные слова
		for _, word := range elements {
			if strings.Contains(ret, word) {
				ret = strings.ReplaceAll(ret, word, "")
			}
		}
	}
	ret = strings.TrimSpace(ret)
	return ret
}

func (c *CardMem) GetItemIn(idx int) *ItemMem {
	item, ok := c.itemsIn[idx]
	if ok {
		return item
	} else {
		return nil
	}
}

func (c *CardMem) GetItemOut(idx int) *ItemMem {
	item, ok := c.itemsOut[idx]
	if ok {
		return item
	} else {
		return nil
	}
}

func (c *CardMem) HasItemOut(doc string, position string) int {
	var docS string
	attempt := 0
	for attempt < 2 {
		if attempt == 0 {
			docS = doc
		} else if attempt == 1 {
			docS = strings.ReplaceAll(strings.ReplaceAll(docS, " ", ""), "-", "")
		}

		for i, val := range c.itemsOut {
			if val.GetOut() == 0 {
				continue
			}
			var docD string
			if attempt == 0 {
				docD = val.document
			} else if attempt == 1 {
				docD = strings.ReplaceAll(strings.ReplaceAll(docD, " ", ""), "-", "")
			}
			if strings.Contains(docS, docD) {
				find := strings.ReplaceAll(strings.ToLower(position), " ", "")
				//fmt.Println(find + "\t" + val.position)
				if strings.Contains(find, val.position) {
					return i
				}
			}
		}
		attempt++
	}
	return -1
}
