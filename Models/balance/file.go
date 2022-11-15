package balance

import (
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type BalanceMem struct {
	fileName string
	items    map[int]*ItemMem
}

type Journal interface {
	GetItemsCount() int
	LoadFromFile(fileName string) bool
	findField(field string, row int, f *excelize.File) int
	findRow(field string, f *excelize.File) int

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

func (b *BalanceMem) findRow(field string, f *excelize.File) int {
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
	f.SetCellValue(sheet, "E1", fldDesc)
	f.SetCellValue(sheet, "F1", "Документ")
	f.SetCellValue(sheet, "G1", "Статус")
	f.SetCellValue(sheet, "H1", "Примечание")

	style, _ := f.NewStyle(`{"alignment":{"horizontal":"center"},"font":{"bold":true}}`)
	f.SetRowStyle(sheet, 1, 1, style)
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
	for i, val := range b.items {
		if val.rest == 0 && val.count == 0 {
			skipCnt++
			continue
		}
		row = i + 2 - skipCnt
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), val.bill)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), val.name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), val.rest)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), val.count)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), val.description)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), val.document)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), val.statusToStr())
		style, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{Type: "pattern", Color: []string{val.statusToColor()}, Pattern: 1},
		})
		f.SetCellStyle(sheet, "G"+strconv.Itoa(row), "G"+strconv.Itoa(row), style)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), val.comment)
	}
	//f.SetColWidth()
	//style2, _ := f.NewStyle(&excelize.Style{
	//	Fill: excelize.Fill{Type: "pattern", Color: []string{val.statusToColor()}, Pattern: 1},
	//})
	//f.SetCellStyle(sheet, "G"+strconv.Itoa(row), "G"+strconv.Itoa(row), style)
	//sheet.Cells[row, Helper.A_FIELD].NumberFormat = "@";
	//sheet.Columns[Helper.A_FIELD + ":" + Helper.H_FIELD].AutoFit();
}
func (b *BalanceMem) LoadFromFile(fileName string) bool {
	//	xlsx, err := excelize.OpenFile(fileName)
	//	if err != nil {
	//		fmt.Println(err)
	//		return false
	//	}
	//	row := j.findRow(fldContent, xlsx)
	//	if row == -1 {
	//		return false
	//	}
	//	iCont := j.findField(fldContent, row, xlsx)
	//	if iCont == -1 {
	//		return false
	//	}
	//	iDoc := j.findField(fldDoc, row, xlsx)
	//	if iDoc == -1 {
	//		return false
	//	}
	//	iAmount := j.findField(fldAmount, row, xlsx)
	//	if iAmount == -1 {
	//		return false
	//	}
	//	i := 0
	//	for i < tryCnt {
	//		row++
	//
	//		cell, _ := excelize.CoordinatesToCellName(iDoc, row)
	//		doc, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
	//		if doc != "" {
	//			i = 0
	//			for {
	//				row++
	//
	//				item := NewItem()
	//				item.SetDocument(doc)
	//
	//				cell, _ := excelize.CoordinatesToCellName(iAmount, row)
	//				restStr, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
	//				if err != nil {
	//					row--
	//					break
	//				}
	//				rest, err := strconv.ParseFloat(restStr, 64)
	//				if err != nil {
	//					row--
	//					break
	//				}
	//
	//				item.SetRest(rest)
	//
	//				cell, _ = excelize.CoordinatesToCellName(iCont, row)
	//				desc, err := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
	//				if err != nil {
	//					row--
	//					break
	//				}
	//				item.SetDescription(desc)
	//
	//				if item.GetDescription() != "" {
	//					j.items[j.GetItemsCount()] = item
	//				}
	//			}
	//		} else {
	//			cell, _ := excelize.CoordinatesToCellName(iCont, row)
	//			content, _ := xlsx.GetCellValue(xlsx.GetSheetName(0), cell)
	//			if content == "" {
	//				i++
	//				continue
	//			}
	//		}
	//	}
	//	if j.GetItemsCount() > 0 {
	//		return true
	//	} else {
	//		return false
	//	}
	return true
}

/*   private bool LoadFromXLS()
{
    try
    {
        Excel.Worksheet objWorksheet;
        objWorksheet = GetActiveSheet(application, fileName);

        int row = FindRow(objWorksheet);
        if (row == -1)
            return false;
        int iBill = FindField(Helper.BILL, row, objWorksheet);
        if (iBill == -1)
            return false;
        int iName = FindField(Helper.NAME, row, objWorksheet);
        if (iName == -1)
            return false;
        int iCount = FindField(Helper.COUNT + " " + Helper.PER_END, row, objWorksheet);
        if (iCount == -1)
            return false;
        int iDesc = FindField(Helper.DESC, row, objWorksheet);
        if (iDesc == -1)
            return false;
        int iRest = FindField(Helper.REST + " " + Helper.PER_END, row, objWorksheet);
        if (iRest == -1)
            return false;

        int i = 0;

        while (i < Helper.TRY_COUNT)
        {
            row++;
            BalanceItem BI = new BalanceItem
            {
                Bill = objWorksheet.Cells[row, iBill].Text.ToString()
            };
            if (BI.Bill.Equals(""))
            {
                i++;
            }
            else
            {
                BI.Description = objWorksheet.Cells[row, iDesc].Text.ToString();
                if (BI.Description.Equals(""))
                {
                    i++;
                    continue;
                }
                BI.Name = objWorksheet.Cells[row, iName].Text.ToString();
                if (BI.Name.Equals(""))
                {
                    i++;
                    continue;
                }
                try
                {
                    BI.Count = int.Parse(objWorksheet.Cells[row, iCount].Text.ToString());
                }
                catch (FormatException)
                {
                    i++;
                    continue;
                }
                try
                {
                    BI.Rest = double.Parse(objWorksheet.Cells[row, iRest].Text.ToString());
                }
                catch (FormatException)
                {
                    i++;
                    continue;
                }
                items.Add(BI);
                i = 1;
            }
        }
        if (items.Count > 0)
        {
            return true;
        }
        else
        {
            return false;
        }
    }
    finally
    {
        application.Workbooks.Close();
    }
}
*/
