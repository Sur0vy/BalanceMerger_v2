package main

import (
	"BM/Models"
	"BM/Models/balance"
	"BM/Models/card"
	"BM/Models/journal"
	"fmt"
	"strconv"
)

var output *balance.BalanceMem

func StartProcess(src Models.Sources) bool {

	jr := journal.NewJournal()
	err := jr.LoadFromFile(src.Journal)
	if err != nil {
		fmt.Println(err)
		return false
	}

	bl := balance.NewBalance()
	err = bl.LoadFromFile(src.Balance)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if src.Card == "" {
		fmt.Println("Process old method")
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
			return false
		}
	} else {
		fmt.Println("Process new method")
		cr := card.NewCard()
		err = cr.LoadFromFile(src.Card)
		if err != nil {
			fmt.Println(err)
			return false
		}
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
			return false
		}
		err = mergeV2(bl, cr)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	output = bl
	return true
}

func SaveMergedFile(fileName string) {
	output.Save(fileName)
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

func mergeV2(bl *balance.BalanceMem, crd *card.CardMem) error {
	for i := 0; i <= bl.GetItemsCount()-1; i++ {
		bi := bl.GetItem(i)
		idx := crd.HasItemOut(bi.GetDocument(), bi.GetDescription())
		if idx == -1 {
			continue
		}
		ci := crd.GetItemOut(idx)
		if ci.GetOut() <= bi.GetCount() {
			bi.SetSpent(ci.GetOut())
			ci.SetOut(0)
		} else {
			bi.SetSpent(bi.GetCount())
			ci.SetOut(ci.GetOut() - bi.GetCount())
		}
	}
	return nil
}
