package main

import (
	"BM/Models"
	"BM/Models/balance"
	"BM/Models/journal"
	"fmt"
)

func StartProcess(src Models.Sources) {

	jr := journal.NewJournal()
	err := jr.LoadFromFile(src.Journal)
	if err != nil {
		fmt.Println(err)
		return
	}

	bl := balance.NewBalance()
	err = bl.LoadFromFile(src.Balance)
	if err != nil {
		fmt.Println(err)
		return
	}

	if src.Card == "" {
		fmt.Println("Process old method")

	} else {
		fmt.Println("Process new method")
		//cr := card.NewCard()
		//err = cr.LoadFromFile(src.Card)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
	}
}
