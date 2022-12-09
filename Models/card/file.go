package card

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

type ItemsArray []*ItemMem

func (a ItemsArray) Len() int           { return len(a) }
func (a ItemsArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ItemsArray) Less(i, j int) bool { return a[i].GetDate().Before(a[j].GetDate()) }

type CardMem struct {
	fileName string
	itemsIn  ItemsArray
	itemsOut ItemsArray
}

type Card interface {
	GetItemsCount() int
	LoadFromFile(fileName string, ch chan int)
	HasItemOut(doc string, position string) int
	GetItemIn(idx int) *ItemMem
	GetItemOut(idx int) *ItemMem
	GetMissing() ItemsArray

	simplifyDocument(val string) string
	simplifyPosition(val string) string
	findField(field string, row int, f *excelize.File) int
	findRow(field string, f *excelize.File) int
	fillOutItems()
}

func NewCard() *CardMem {
	return &CardMem{}
}

func (c *CardMem) LoadFromFile(fileName string, ch chan int) {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		ch <- -1
		return
	}
	row := c.findRow(fldDoc, xlsx)
	if row == -1 {
		ch <- -2
		return
	}

	iDoc := c.findField(fldDoc, row, xlsx)
	if iDoc == -1 {
		ch <- -2
		return
	}

	//get row count
	var count int
	step := 100
	for step > 0 {
		cell, _ := excelize.CoordinatesToCellName(iDoc, count)
		doc, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
		if doc == "" {
			if count > 0 {
				count -= step
				step = step / 2
				continue
			}
		}
		count += step
	}
	coeff := 100.0 / float64(count)

	iDocD := c.findField(fldDocD, row, xlsx)
	if iDocD == -1 {
		ch <- -2
		return
	}
	iDocC := c.findField(fldDocC, row, xlsx)
	if iDocC == -1 {
		ch <- -2
		return
	}
	iCntC := c.findField(fldCntC, row, xlsx)
	if iCntC == -1 {
		ch <- -2
		return
	} else {
		iCntC++
	}
	iCntD := c.findField(fldCntD, row, xlsx)
	if iCntD == -1 {
		ch <- -2
		return
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
			docDate, err := c.getDocumentDate(doc)
			if err != nil {
				i++
				continue
			}
			item.SetDate(docDate)
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
				c.itemsIn = append(c.itemsIn, item)
			} else {
				find := false
				for i, v := range c.itemsOut {
					if v.document == item.document && v.position == item.position {
						c.itemsOut[i].SetOut(v.GetOut() + item.GetOut())
						find = true
						item = nil
						break
					}
				}
				if !find {
					c.itemsOut = append(c.itemsOut, item)
				}
			}
			ch <- int(coeff * float64(row))
		}
	}
	if c.GetItemsCount() > 0 {
		c.fillOutItems()
		ch <- 100
		return
	} else {
		ch <- -3
		return
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
	//отсортируем элементы по дате поступления на склад
	sort.Sort(c.itemsIn)
	for i := 0; i < len(c.itemsOut); i++ {
		out := c.GetItemOut(i)
		for j, _ := range c.itemsIn {
			if out.GetNeed() == 0 {
				break
			}
			in := c.GetItemIn(j)
			if in.GetIn() == 0 {
				continue
			}
			if out.position == in.position {
				if out.GetNeed() <= in.GetIn() {
					out.SetIn(out.GetOut())
					in.SetIn(in.GetIn() - out.GetOut())
					out.SetDocument(in.GetDocument())
					break
				} else {
					item := NewItem()
					item.SetPosition(out.position)
					item.SetOut(out.GetOut() - in.GetIn())
					c.itemsOut = append(c.itemsOut, item)

					out.SetIn(in.GetIn())
					out.SetOut(in.GetIn())
					in.SetIn(0)
					out.SetDocument(in.GetDocument())
				}
			}
		}
	}
}

func (c *CardMem) getDocumentDate(val string) (time.Time, error) {
	pos := strings.LastIndex(val, "от")
	if pos == -1 {
		return time.Time{}, errors.New("no date")
	}
	dateStr := val[pos:]
	dateStr = strings.TrimSpace(strings.Split(dateStr, " ")[1])
	return time.Parse("02.01.2006", dateStr)
}

func (c *CardMem) simplifyDocument(val string) string {
	val = strings.ToLower(val)
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

	ret = strings.TrimSpace(ret)
	return ret
}

func (c *CardMem) GetItemIn(idx int) *ItemMem {
	if idx > -1 && idx < len(c.itemsIn) {
		return c.itemsIn[idx]
	} else {
		return nil
	}
}

func (c *CardMem) GetItemOut(idx int) *ItemMem {
	if idx > -1 && idx < len(c.itemsOut) {
		return c.itemsOut[idx]
	} else {
		return nil
	}
}

func (c *CardMem) HasItemOut(doc string, position string) int {
	//упростим поисковую строку до нескольких вариантов
	var val []string
	val = append(val, strings.ToLower(doc))
	val = append(val, strings.ReplaceAll(val[0], " ", ""))
	val = append(val, strings.ReplaceAll(val[1], "-", ""))

	for i, item := range c.itemsOut {
		if item.GetOut() == 0 || item.document == "" {
			continue
		}
		var dict []string
		dict = append(dict, item.document)
		dict = append(dict, strings.ReplaceAll(dict[0], " ", ""))
		dict = append(dict, strings.ReplaceAll(dict[1], "-", ""))

		for _, v := range val {
			for _, d := range dict {
				if strings.Contains(v, d) {
					if c.checkPosition(position, item.position) {
						return i
					}
				}
			}
		}
	}
	return -1
}

func (c *CardMem) GetMissing() ItemsArray {
	var items ItemsArray
	for _, item := range c.itemsOut {
		if item.GetOut() != 0 &&
			item.document != "" {

			tmp := item.GetDocument()
			tmp = strings.ReplaceAll(tmp, " ", "")
			tmp = strings.ReplaceAll(tmp, "-", "")
			item.SetDocument(tmp)
			item.SetComment(item.GetPosition())
			item.SetPosition(RemeoveElements(item.GetPosition()))
			items = append(items, item)
		}
	}
	return items
}

func (c *CardMem) getDict(position string) []string {
	var dict []string
	position = strings.ToLower(position)
	dict = append(dict, position)
	tmp := dict[0]
	if strings.Contains(position, " ") {
		tmp = strings.ReplaceAll(tmp, " ", "")
		dict = append(dict, tmp)
	}
	if strings.Contains(position, "-") {
		tmp = strings.ReplaceAll(dict[0], " ", strings.ReplaceAll(tmp, "-", ""))
		dict = append(dict, tmp)
	}
	if strings.Contains(position, ".") {
		tmp = strings.ReplaceAll(tmp, ".", "")
		dict = append(dict, tmp)
	}
	if strings.Contains(position, ",") {
		tmp = strings.ReplaceAll(tmp, ",", "")
		dict = append(dict, tmp)
	}

	if strings.Contains(position, "резистор") || strings.Contains(position, "конденсатор") {
		tmp := position
		start := strings.Index(tmp, "(")
		stop := strings.Index(tmp, ")")
		if start < stop {
			tmp = tmp[start+1 : stop]
		} else {
			tmp = strings.ReplaceAll(tmp, "резистор", "")
		}
		dict = append(dict, strings.TrimSpace(tmp))
	} else {
		//удалим реперные слова
		tmp := position
		for _, word := range Elements {
			if strings.Contains(tmp, word) {
				tmp = strings.ReplaceAll(tmp, word, "")
				tmp = strings.ReplaceAll(tmp, " ", "")
				tmp = strings.ReplaceAll(tmp, ".", "")
				tmp = strings.ReplaceAll(tmp, ",", "")
				tmp = strings.ReplaceAll(tmp, "-", "")
				dict = append(dict, tmp)
				break
			}
		}
	}
	return dict
}

func (c *CardMem) checkPosition(position string, str string) bool {
	val := c.getDict(position)
	dict := c.getDict(str)

	for _, v := range val {
		for _, d := range dict {
			if strings.Contains(v, d) || strings.Contains(d, v) {
				return true
			}
		}
	}
	return false
}
