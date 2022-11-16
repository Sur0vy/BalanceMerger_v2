package main

import (
	"BM/Models"
	"BM/Models/balance"
	"BM/Models/journal"
	"fmt"
	"strconv"
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
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
		} else {
			//bl.Save(src.Balance + "111")
			bl.Save("/Users/Sur0vy/out.xlsx")
		}
	} else {
		fmt.Println("Process new method")
		//cr := card.NewCard()
		//err = cr.LoadFromFile(src.Card)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		err := mergeV2()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func mergeV1(bl *balance.BalanceMem, jr *journal.JournalMem) error {
	for i := 0; i <= bl.GetItemsCount()-1; i++ {
		bi := bl.GetItem(i)
		st, indexes := jr.HasItem(bi.GetDescription(), bi.GetRest())
		bi.SetState(st)
		switch st {
		case Models.IsFound:
			bi.SetDocument(jr.GetItem(indexes[0]).GetDocument())
			break
		case Models.IsCollect:
			for idx := range indexes {
				if bi.GetDocument() != "" {
					bi.SetComment(bi.GetDocument() + " ")
				}
				bi.SetDocument(bi.GetDocument() + jr.GetItem(indexes[idx]).GetDocument())
				if bi.GetComment() != "" {
					bi.SetComment(bi.GetComment() + " ")
				}
				bi.SetComment(bi.GetComment() + jr.GetItem(indexes[idx]).GetDescription())

			}
			break
		case Models.IsDifBalance:
			bi.SetComment("Остаток в журнале: " + strconv.FormatFloat(jr.GetItem(indexes[0]).GetAmount(), 'f', 6, 64))
			break
		case Models.IsCollectMissing:
			for idx := range indexes {
				if bi.GetComment() != "" {
					bi.SetComment(bi.GetComment() + " ")
				}
				newComment := bi.GetComment() +
					jr.GetItem(indexes[0]).GetDescription() + " (" +
					strconv.FormatFloat(jr.GetItem(indexes[idx]).GetAmount(), 'f', 6, 64) + ")"
				bi.SetComment(newComment)
			}
			break
		default:
			//
			break
		}
	}
	return nil
}

func mergeV2() error {
	return nil
}
