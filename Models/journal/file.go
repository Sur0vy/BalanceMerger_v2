package journal

import (
	"BM/Models"
	"container/list"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type JournalMem struct {
	fileName string
	items    *list.List
}

type Journal interface {
	GetItemsCount() int
	LoadFromFile(fileName string) bool
	HasItem(name string, rest float64, indexes *list.List) Models.ItemState
	GetItem(idx int) *Item
	findField(field string, row int, f *excelize.File) int
	findRow(field string, f *excelize.File) int
}

func NewJournal() *JournalMem {
	return &JournalMem{
		items: list.New(),
	}
}

func (j *JournalMem) LoadFromFile(fileName string) bool {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return false
	}
	row := j.findRow(fldContent, xlsx)
	if row == -1 {
		return false
	}
	iCont := j.findField(fldContent, row, xlsx)
	if iCont == -1 {
		return false
	}
	iDoc := j.findField(fldDoc, row, xlsx)
	if iDoc == -1 {
		return false
	}
	iAmount := j.findField(fldAmount, row, xlsx)
	if iAmount == -1 {
		return false
	}
	i := 0
	for i < tryCnt {
		row++
		//document :=
		cell, _ := excelize.CoordinatesToCellName(iDoc, row)
		doc, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
		if doc != "" {
			i = 0
			for {
				row++

				var item Item = NewItem()
				item.SetDocument(doc)

				cell, _ := excelize.CoordinatesToCellName(iAmount, row)
				restStr, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				if err != nil {
					row--
					break
				}
				rest, err := strconv.ParseFloat(restStr, 64)
				if err != nil {
					row--
					break
				}

				item.SetRest(rest)

				cell, _ = excelize.CoordinatesToCellName(iCont, row)
				desc, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
				if err != nil {
					row--
					break
				}
				item.SetDescription(desc)

				if item.GetDescription() != "" {
					j.items.PushBack(item)
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
		return true
	} else {
		return false
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
	return j.items.Len()
}

func (j *JournalMem) GetItem(idx int) *Item {
	i := 0
	for e := j.items.Front(); e != nil; e = e.Next() {
		if i == idx {
			return e.Value.(*Item)
		}
		i++
	}
	return nil
}

//func (j *JournalMem) HasItem(name string, rest float64, indexes *list.List) Models.ItemState {
//	var index int
//	for e := j.items.Front(); e != nil; e = e.Next() {
//		itm := e.Value.(*Item)
//		//index = e.Value.(*Item)
//	}
//}

/*   public ItemState HasItem(string name, double rest, ref List<int> indexes)
     {
         int index;
         for (int i = 0; i < items.Count; i++)
         {
             index = items[i].Description.IndexOf(name);
             if (index > -1)
             {
                 indexes.Add(i);
                 if (items[indexes[0]].Rest == rest)
                     return ItemState.isFound;
             }
         }
         if (indexes.Count == 0)
         {
             return ItemState.isMissing;
         }
         else if (indexes.Count == 1)
         {
             return ItemState.isDifBalance;
         }
         else
         {
             double b = 0;
             for (int i = 0; i < indexes.Count; i++)
             {
                 b = b + GetItem(indexes[i]).Rest;
             }
             if (b == rest)
             {
                 return ItemState.isCollect;
             }
             else
             {
                 return ItemState.isCollectMissing;
             }
         }
     }
*/
