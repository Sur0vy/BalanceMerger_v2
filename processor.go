package main

import (
	"BM/Models"
	"BM/Models/balance"
	"BM/Models/card"
	"BM/Models/journal"
	"fmt"
	"strconv"
)

type ResponseStep int

const (
	RiJournal ResponseStep = iota
	RiBalance
	RiCard
	RiProcess
)

type response struct {
	Status   bool
	Step     ResponseStep
	Progress float64
}

var output *balance.BalanceMem

func StartProcess(src Models.Sources, c chan response) {

	jr := journal.NewJournal()
	err := jr.LoadFromFile(src.Journal)
	if err != nil {
		fmt.Println(err)
		c <- response{
			Progress: 0,
			Step:     RiJournal,
			Status:   false,
		}
		return
	}
	c <- response{
		Progress: 100,
		Step:     RiJournal,
		Status:   true,
	}
	bl := balance.NewBalance()
	err = bl.LoadFromFile(src.Balance)
	if err != nil {
		fmt.Println(err)
		c <- response{
			Progress: 0,
			Step:     RiBalance,
			Status:   false,
		}
		return
	}
	c <- response{
		Progress: 100,
		Step:     RiBalance,
		Status:   true,
	}
	if src.Card == "" {
		fmt.Println("Process old method")
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
			c <- response{
				Progress: 0,
				Step:     RiProcess,
				Status:   false,
			}
			return
		}
	} else {
		fmt.Println("Process new method")
		cr := card.NewCard()
		err = cr.LoadFromFile(src.Card)
		if err != nil {
			fmt.Println(err)
			c <- response{
				Progress: 0,
				Step:     RiCard,
				Status:   false,
			}
			return
		}
		c <- response{
			Progress: 100,
			Step:     RiCard,
			Status:   true,
		}
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
			c <- response{
				Progress: 0,
				Step:     RiProcess,
				Status:   false,
			}
			return
		}
		c <- response{
			Progress: 10,
			Step:     RiProcess,
			Status:   true,
		}
		//todo start in gorutine
		err = mergeV2(bl, cr)
		if err != nil {
			fmt.Println(err)
			c <- response{
				Progress: 0,
				Step:     RiProcess,
				Status:   false,
			}
			return
		}
		//end todo
	}
	output = bl
	c <- response{
		Progress: 100,
		Step:     RiProcess,
		Status:   true,
	}
}

func SaveMergedFile(fileName string) {
	if output.Save(fileName) != nil {
		lbMessage.SetText("Ошибка при сохранении результатов обработки!")
	} else {
		lbMessage.SetText("Результаты обработки успешно сохранены!")
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
