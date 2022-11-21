package balance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type BalanceMem struct {
	fileName string
	items    map[int]*ItemMem
}

type Balance interface {
	GetItemsCount() int
	LoadFromFile(fileName string) error
	findField(field string, row int, f *excelize.File) int
	findRow(f *excelize.File) int
	GetItem(idx int) *ItemMem

	Save(fileName string) error
	makeHeader(sheet string, f *excelize.File)
	saveData(sheet string, f *excelize.File)
}

func NewBalance() *BalanceMem {
	return &BalanceMem{
		fileName: "",
		items:    make(map[int]*ItemMem),
	}
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

	f.SetCellValue(sheet, "A1", fldBill)
	f.SetCellValue(sheet, "B1", fldName)
	f.SetCellValue(sheet, "C1", fldRest+fldPerEnd)
	f.SetCellValue(sheet, "D1", fldCount+fldPerEnd)
	f.SetCellValue(sheet, "E1", "Списание")
	f.SetCellValue(sheet, "F1", "Остаток после списания")
	f.SetCellValue(sheet, "G1", fldDesc)
	f.SetCellValue(sheet, "H1", "Документ")
	f.SetCellValue(sheet, "I1", "Статус")
	f.SetCellValue(sheet, "J1", "Примечание")

	style, _ := f.NewStyle(`{"alignment":{"horizontal":"center"},"font":{"bold":true}}`)
	f.SetRowStyle(sheet, 1, 1, style)
	f.SetColWidth(sheet, "A", "J", 16)
	f.SetColWidth(sheet, "C", "D", 32)
	f.SetColWidth(sheet, "G", "H", 60)

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
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), val.spent)
		formula := "=D" + strconv.Itoa(row) + "-E" + strconv.Itoa(row)
		f.SetCellFormula(sheet, "F"+strconv.Itoa(row), formula)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), val.description)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), val.document)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), val.statusToStr())
		style, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{val.statusToColor()}, Pattern: 1},
		})
		f.SetCellStyle(sheet, "I"+strconv.Itoa(row), "I"+strconv.Itoa(row), style)
		f.SetCellValue(sheet, "J"+strconv.Itoa(row), val.comment)
	}
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
			b.items[b.GetItemsCount()] = item
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
	item, ok := b.items[idx]
	if ok {
		return item
	} else {
		return nil
	}
}
