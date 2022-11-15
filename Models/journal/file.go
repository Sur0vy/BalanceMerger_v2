package journal

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"

	"BM/Models"
)

type JournalMem struct {
	fileName string
	items    map[int]*ItemMem
}

type Journal interface {
	GetItemsCount() int
	LoadFromFile(fileName string) error
	HasItem(name string, rest float64) (Models.ItemState, *[]int)
	GetItem(idx int) *Item
	findField(field string, row int, f *excelize.File) int
	findRow(field string, f *excelize.File) int
}

func NewJournal() *JournalMem {
	return &JournalMem{
		fileName: "",
		items:    make(map[int]*ItemMem),
	}
}

func (j *JournalMem) LoadFromFile(fileName string) error {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return errors.New("file is corrupted")
	}
	row := j.findRow(fldContent, xlsx)
	if row == -1 {
		return errors.New("file read error")
	}
	iCont := j.findField(fldContent, row, xlsx)
	if iCont == -1 {
		return errors.New("file read error")
	}
	iDoc := j.findField(fldDoc, row, xlsx)
	if iDoc == -1 {
		return errors.New("file read error")
	}
	iAmount := j.findField(fldAmount, row, xlsx)
	if iAmount == -1 {
		return errors.New("file read error")
	}
	i := 0
	for i < tryCnt {
		row++

		cell, _ := excelize.CoordinatesToCellName(iDoc, row)
		doc, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
		if doc != "" {
			i = 0
			for {
				row++

				item := NewItem()
				item.SetDocument(doc)

				cell, _ := excelize.CoordinatesToCellName(iAmount, row)
				amountStr, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				if err != nil {
					row--
					break
				}
				amount, err := strconv.ParseFloat(amountStr, 64)
				if err != nil {
					row--
					break
				}

				item.SetAmount(amount)

				cell, _ = excelize.CoordinatesToCellName(iCont, row)
				desc, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				if err != nil {
					row--
					break
				}
				item.SetDescription(desc)

				if item.GetDescription() != "" {
					fmt.Println(item)
					j.items[j.GetItemsCount()] = item
				}
			}
		} else {
			cell, _ := excelize.CoordinatesToCellName(iCont, row)
			content, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
			if content == "" {
				i++
				continue
			}
		}
	}
	if j.GetItemsCount() > 0 {
		return nil
	} else {
		return errors.New("no items in file")
	}
}

func (j *JournalMem) findField(field string, row int, f *excelize.File) int {
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

func (j *JournalMem) findRow(field string, f *excelize.File) int {
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

func (j *JournalMem) GetItemsCount() int {
	return len(j.items)
}

func (j *JournalMem) GetItem(idx int) *ItemMem {
	item, ok := j.items[idx]
	if ok {
		return item
	} else {
		return nil
	}
}

func (j *JournalMem) HasItem(name string, rest float64) (Models.ItemState, *[]int) {
	indexes := []int{}
	for i, val := range j.items {
		if strings.ContainsAny(val.GetDescription(), name) {
			indexes = append(indexes, i)
			if val.GetAmount() == rest {
				return Models.IsFound, &indexes
			}
		}
	}
	if len(indexes) == 0 {
		return Models.IsMissing, &indexes
	} else if len(indexes) == 1 {
		return Models.IsDifBalance, &indexes
	} else {
		b := 0.0
		for i := range indexes {
			b += j.GetItem(i).GetAmount()
		}
		if b == rest {
			return Models.IsCollect, &indexes
		} else {
			return Models.IsCollectMissing, &indexes
		}
	}
}
