package main

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func StartProcess(src Sources) {
	xlsx, err := excelize.OpenFile(src.journal)
	if err != nil {
		fmt.Println(err)
		return
	}

	sheetMap := xlsx.GetSheetMap()
	for _, sheetName := range sheetMap {
		// Get value from cell by given worksheet name and axis.
		cell := xlsx.GetCellValue(sheetName, "A5")
		fmt.Println(cell)
		// Get all the rows in the Sheet1.
		rows := xlsx.GetRows(sheetName)
		for _, row := range rows {
			for _, colCell := range row {
				fmt.Print(colCell, "\t")
			}
			fmt.Println()
		}
	}
}
