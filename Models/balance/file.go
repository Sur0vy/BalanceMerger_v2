package balance

import (
	"BM/Models"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ItemsArray []*ItemMem

func (a ItemsArray) Len() int           { return len(a) }
func (a ItemsArray) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ItemsArray) Less(i, j int) bool { return a[i].GetDate().Before(a[j].GetDate()) }

type BalanceMem struct {
	fileName string
	state    BalanceState
	items    ItemsArray
}

type Balance interface {
	GetItemsCount() int
	LoadFromFile(fileName string) error
	findField(field string, row int, f *excelize.File) int
	findRow(f *excelize.File) int
	GetItem(idx int) *ItemMem
	SortByDate()
	Save(fileName string) error
	SetState(state BalanceState)
	GetState() BalanceState

	makeHeader(sheet string, f *excelize.File)
	saveData(sheet string, f *excelize.File)
}

func NewBalance() *BalanceMem {
	return &BalanceMem{
		state: IsEmpty,
	}
}

func (b *BalanceMem) SetState(state BalanceState) {
	b.state = state
}

func (b *BalanceMem) GetState() BalanceState {
	return b.state
}

func (b *BalanceMem) AddItem(item *ItemMem) {
	b.items = append(b.items, item)
}

func (b *BalanceMem) findField(field string, row int, f *excelize.File) int {
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

func (b *BalanceMem) findRow(f *excelize.File) int {
	i := 1
	value := ""
	for i < tryCnt {
		cell, _ := excelize.CoordinatesToCellName(1, i)
		value, _ = f.GetCellValue(f.GetSheetName(0), cell)
		if value == "" {
			i++
		} else {
			return i
		}
	}
	return -1
}

func (b *BalanceMem) GetItemsCount() int {
	return len(b.items)
}

func (b *BalanceMem) makeHeader(sheet string, f *excelize.File) {

	row := 1
	col := 1
	cell, _ := excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, fldBill)
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, fldName)
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, fldRest+fldPerEnd)
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, fldCount+fldPerEnd)

	var beginCol string
	var endCol string
	if b.state == IsMergeV2 {
		col++
		cell, _ = excelize.CoordinatesToCellName(col, row)
		f.SetCellValue(sheet, cell, "Списание")
		col++
		cell, _ = excelize.CoordinatesToCellName(col, row)
		f.SetCellValue(sheet, cell, "Остаток после списания")
	}
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, fldDesc)
	if b.state == IsMergeV2 {
		col++
		cell, _ = excelize.CoordinatesToCellName(col, row)
		f.SetCellValue(sheet, cell, fldDesc+"(Карточка)")
		beginCol = "G"
		endCol = "I"
	} else {
		beginCol = "E"
		endCol = "F"
	}
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, "Документ")
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, "Статус")
	col++
	cell, _ = excelize.CoordinatesToCellName(col, row)
	f.SetCellValue(sheet, cell, "Примечание")

	style, _ := f.NewStyle(`{"alignment":{"horizontal":"center"},"font":{"bold":true}}`)
	f.SetRowStyle(sheet, 1, 1, style)

	f.SetColWidth(sheet, "A", "K", 16)
	f.SetColWidth(sheet, "C", "D", 32)
	f.SetColWidth(sheet, beginCol, endCol, 46)
}

func (b *BalanceMem) Save(fileName string) error {
	f := excelize.NewFile()
	sheetName := "Баланс"
	f.SetSheetName("Sheet1", sheetName)
	b.makeHeader(sheetName, f)
	b.saveData(sheetName, f)
	if err := f.SaveAs(fileName); err != nil {
		return err
	}
	return nil
}

func (b *BalanceMem) saveData(sheet string, f *excelize.File) {
	skipCnt := 0
	row := 0
	i := 1
	for _, val := range b.items {
		i++
		if val.rest == 0 && val.count == 0 {
			skipCnt++
			continue
		}
		row = i - skipCnt
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), val.bill)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), val.name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), val.rest)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), val.count)

		var cell string
		col := 5
		var colStat string
		var colComment string
		if b.state == IsMergeV2 {
			f.SetCellValue(sheet, "E"+strconv.Itoa(row), val.spent)
			if val.spent == 0 {
				style, _ := f.NewStyle(&excelize.Style{
					Fill: excelize.Fill{Type: "pattern", Color: []string{ColorFromStatus(Models.IsMissing)}, Pattern: 1},
				})
				f.SetCellStyle(sheet, "E"+strconv.Itoa(row), "F"+strconv.Itoa(row), style)
			}
			formula := "D" + strconv.Itoa(row) + "-E" + strconv.Itoa(row)
			f.SetCellFormula(sheet, "F"+strconv.Itoa(row), formula)
			col = 7
		}

		cell, _ = excelize.CoordinatesToCellName(col, row)
		f.SetCellValue(sheet, cell, val.description)
		col++
		if b.state == IsMergeV2 {
			cell, _ = excelize.CoordinatesToCellName(col, row)
			f.SetCellValue(sheet, cell, val.position)

			style, _ := f.NewStyle(&excelize.Style{
				Fill: excelize.Fill{Type: "pattern", Color: []string{ColorFromMatchPers(val.GetAccuracy())}, Pattern: 1},
			})
			f.SetCellStyle(sheet, cell, cell, style)

			col++
			colStat = "J"
			colComment = "K"
		} else {
			colStat = "G"
			colComment = "H"
		}
		cell, _ = excelize.CoordinatesToCellName(col, row)
		f.SetCellValue(sheet, cell, val.document)

		col++
		cell, _ = excelize.CoordinatesToCellName(col, row)
		f.SetCellValue(sheet, cell, val.statusToStr())
		style, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{val.statusToColor()}, Pattern: 1},
		})
		f.SetCellStyle(sheet, colStat+strconv.Itoa(row), colStat+strconv.Itoa(row), style)
		f.SetCellValue(sheet, colComment+strconv.Itoa(row), val.comment)
	}

	formula := "sum(D2:D" + strconv.Itoa(row) + ")"
	f.SetCellFormula(sheet, "D"+strconv.Itoa(row+1), formula)

	formula = "sum(E2:E" + strconv.Itoa(row) + ")"
	f.SetCellFormula(sheet, "E"+strconv.Itoa(row+1), formula)

	formula = "sum(F2:F" + strconv.Itoa(row) + ")"
	f.SetCellFormula(sheet, "F"+strconv.Itoa(row+1), formula)
}

func (b *BalanceMem) LoadFromFile(fileName string) error {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return errors.New("file is corrupted")
	}

	row := b.findRow(xlsx)
	if row == -1 {
		return errors.New("file read error")
	}
	iBill := b.findField(fldBill, row, xlsx)
	if iBill == -1 {
		return errors.New("file read error")
	}
	iName := b.findField(fldName, row, xlsx)
	if iName == -1 {
		return errors.New("file read error")
	}
	iDesc := b.findField(fldDesc, row, xlsx)
	if iDesc == -1 {
		return errors.New("file read error")
	}
	iCount := b.findField(fldCount+fldPerEnd, row, xlsx)
	if iCount == -1 {
		return errors.New("file read error")
	}
	iRest := b.findField(fldRest+fldPerEnd, row, xlsx)
	if iRest == -1 {
		return errors.New("file read error")
	}
	i := 0
	for i < tryCnt {
		row++
		cell, _ := excelize.CoordinatesToCellName(iBill, row)
		bill, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
		if bill == "" {
			i++
		} else {
			item := NewItem()
			item.SetBill(bill)

			cell, _ = excelize.CoordinatesToCellName(iDesc, row)
			desc, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
			if err != nil {
				i++
				continue
			}
			item.SetDescription(desc)

			cell, _ = excelize.CoordinatesToCellName(iName, row)
			name, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
			if err != nil {
				i++
				continue
			}
			item.SetName(name)

			cell, _ = excelize.CoordinatesToCellName(iCount, row)
			countStr, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
			if err != nil {
				i++
				continue
			}
			count, err := strconv.ParseInt(countStr, 10, 64)
			if err != nil {
				i++
				continue
			}
			item.SetCount(int(count))

			cell, _ = excelize.CoordinatesToCellName(iRest, row)
			restStr, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
			if err != nil {
				i++
				continue
			}
			rest, err := strconv.ParseFloat(restStr, 64)
			if err != nil {
				i++
				continue
			}
			item.SetRest(rest)
			b.items = append(b.items, item)
			i = 1
		}
	}
	if b.GetItemsCount() > 0 {
		return nil
	} else {
		return errors.New("no items in file")
	}
}

func (b *BalanceMem) GetItem(idx int) *ItemMem {
	if idx > -1 && idx < len(b.items) {
		return b.items[idx]
	} else {
		return nil
	}
}

func (b *BalanceMem) SortByDate() {
	sort.Sort(b.items)
}
