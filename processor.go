package main

import (
	"BM/Models"
	"BM/Models/balance"
	"BM/Models/card"
	"BM/Models/journal"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

func StartProcess(src Models.Sources, ch chan response) {

	jr := journal.NewJournal()
	err := jr.LoadFromFile(src.Journal)
	if err != nil {
		fmt.Println(err)
		ch <- response{
			Progress: 0,
			Step:     RiJournal,
			Status:   false,
		}
		return
	}
	ch <- response{
		Progress: 100,
		Step:     RiJournal,
		Status:   true,
	}
	bl := balance.NewBalance()
	err = bl.LoadFromFile(src.Balance)
	if err != nil {
		fmt.Println(err)
		ch <- response{
			Progress: 0,
			Step:     RiBalance,
			Status:   false,
		}
		return
	}
	ch <- response{
		Progress: 100,
		Step:     RiBalance,
		Status:   true,
	}
	if src.Card == "" {
		fmt.Println("Process old method")
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
			ch <- response{
				Progress: 0,
				Step:     RiProcess,
				Status:   false,
			}
			return
		}
	} else {
		fmt.Println("Process new method")
		cr := card.NewCard()
		ans := make(chan int)

		go cr.LoadFromFile(src.Card, ans)
		for {
			time.Sleep(10 * time.Microsecond)
			val, ok := <-ans
			if ok == false {
				ch <- response{
					Progress: 0,
					Step:     RiCard,
					Status:   false,
				}
				return
			} else {
				if val < 0 {
					fmt.Println("file read error")
					ch <- response{
						Progress: 0,
						Step:     RiCard,
						Status:   false,
					}
					return
				}
				ch <- response{
					Progress: float64(val),
					Step:     RiCard,
					Status:   true,
				}
				if val == 100 {
					break
				}
			}
		}
		err := mergeV1(bl, jr)
		if err != nil {
			fmt.Println(err)
			ch <- response{
				Progress: 0,
				Step:     RiProcess,
				Status:   false,
			}
			return
		}
		ch <- response{
			Progress: 10,
			Step:     RiProcess,
			Status:   true,
		}

		go mergeV2(bl, cr, ans)
		for {
			time.Sleep(10 * time.Microsecond)
			val, ok := <-ans
			if ok == false {
				ch <- response{
					Progress: 0,
					Step:     RiProcess,
					Status:   false,
				}
				return
			} else {
				if val < 0 {
					fmt.Println("file read error")
					ch <- response{
						Progress: 0,
						Step:     RiProcess,
						Status:   false,
					}
					return
				}
				ch <- response{
					Progress: float64(val),
					Step:     RiProcess,
					Status:   true,
				}
				if val == 100 {
					break
				}
			}
		}
	}
	output = bl
	ch <- response{
		Progress: 100,
		Step:     RiProcess,
		Status:   true,
	}
}

func SaveMergedFile(fileName string) {
	if filepath.Ext(fileName) == "" {
		os.Remove(fileName)
		fileName += ".xlsx"
	}
	if output.Save(fileName) != nil {
		lbMessage.SetText("Ошибка при сохранении результатов обработки!")
	} else {
		lbMessage.SetText(fmt.Sprintf("Результаты обработки успешно сохранены в файл \"%s\"!", fileName))
	}
}

func dateFromDocument(val string) (time.Time, error) {
	if val == "" {
		return time.Time{}, errors.New("no date")
	}
	dateStr := strings.TrimSpace(strings.Split(strings.Split(val, ";")[0], ",")[2])
	return time.Parse("02.01.2006", dateStr)
}

func mergeV1(bl *balance.BalanceMem, jr *journal.JournalMem) error {
	for i := 0; i <= bl.GetItemsCount()-1; i++ {
		bi := bl.GetItem(i)
		st, indexes := jr.HasItem(bi.GetDescription(), bi.GetRest())
		bi.SetState(st)
		switch st {
		case Models.IsFound:
			bi.SetDocument(jr.GetItem(indexes[0]).GetDocument())
			date, err := dateFromDocument(jr.GetItem(indexes[0]).GetDocument())
			if err == nil {
				bi.SetDate(date)
			}
			break
		case Models.IsCollect:
			for idx := range indexes {
				date, err := dateFromDocument(jr.GetItem(indexes[0]).GetDocument())
				if bi.GetDocument() != "" {
					bi.SetComment(bi.GetDocument() + " ")
					if err == nil {
						if date.Before(bi.GetDate()) {
							bi.SetDate(date)
						}
					}
				} else {
					if err == nil {
						bi.SetDate(date)
					}
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
	//sort balance
	bl.SortByDate()
	bl.SetState(balance.IsMergeV1)
	return nil
}

func mergeV2(bl *balance.BalanceMem, crd *card.CardMem, ch chan int) {
	emptyBl := balance.NewBalance()
	coeff := 100.0 / float64(bl.GetItemsCount())
	for i := 0; i <= bl.GetItemsCount()-1; i++ {
		bi := bl.GetItem(i)
		idx := crd.HasItemOut(bi.GetDocument(), bi.GetDescription())
		if idx == -1 {
			if bi.GetDocument() != "" {
				emptyItem := balance.NewItem()
				emptyItem.SetParent(bi)
				//упростим название документа
				tmp := strings.ToLower(bi.GetDocument())
				tmp = strings.ReplaceAll(tmp, " ", "")
				tmp = strings.ReplaceAll(tmp, ".", "")
				tmp = strings.ReplaceAll(tmp, ",", "")
				tmp = strings.ReplaceAll(tmp, "-", "")
				emptyItem.SetDocument(tmp)
				//упростим насвание позиции
				emptyItem.SetDescription(card.RemeoveElements(bi.GetDescription()))
				emptyBl.AddItem(emptyItem)
			}
			continue
		}
		ch <- int(coeff * float64(i-emptyBl.GetItemsCount()))
		ci := crd.GetItemOut(idx)
		if ci.GetOut() <= bi.GetCount() {
			bi.SetSpent(ci.GetOut())
			ci.SetOut(0)
		} else {
			bi.SetSpent(bi.GetCount())
			ci.SetOut(ci.GetOut() - bi.GetCount())
		}
	}
	emptyCr := crd.GetMissing()

	weights := make([][]float64, emptyBl.GetItemsCount())
	for i := range weights {
		weights[i] = make([]float64, len(emptyCr))
	}

	for i := 0; i <= emptyBl.GetItemsCount()-1; i++ {
		b := emptyBl.GetItem(i)
		for j := 0; j <= emptyCr.Len()-1; j++ {
			c := emptyCr[j]
			if c.GetOut() == 0 || c.GetDocument() == "" {
				continue
			}
			if strings.Contains(b.GetDocument(), c.GetDocument()) {
				w := getWeight(b.GetDescription(), c.GetPosition())
				weights[i][j] = w
			}
		}
		ch <- int(coeff * float64(i+bl.GetItemsCount()-emptyBl.GetItemsCount()))
	}

	max := 0.0
	maxI := -1
	maxJ := -1

	for {

		fmt.Println("step")
		for i := 0; i <= emptyBl.GetItemsCount()-1; i++ {
			fmt.Println("")
			for j := 0; j <= emptyCr.Len()-1; j++ {
				if weights[i][j] == 0 {
					fmt.Print("--- ")
				} else {
					fmt.Printf("%.1f ", weights[i][j])
				}
			}
		}

		for i := 0; i <= emptyBl.GetItemsCount()-1; i++ {
			for j := 0; j <= emptyCr.Len()-1; j++ {
				w := weights[i][j]
				if w > max {
					max = w
					maxI = i
					maxJ = j
				}
			}
		}
		if max == 0 {
			break
		}
		b := emptyBl.GetItem(maxI)
		ci := emptyCr[maxJ]
		//if ci.GetOut() <= b.GetParent().GetCount() {
		b.GetParent().SetSpent(ci.GetOut())
		b.GetParent().SetAccuracy(1.0 - max)
		ci.SetOut(0)

		//for i := 0; i <= emptyBl.GetItemsCount()-1; i++ {
		//	weights[i][maxJ] = 0.0
		//}

		//} else {
		//	b.GetParent().SetSpent(b.GetParent().GetCount())
		//	b.GetParent().SetAccuracy(max)
		//	ci.SetOut(ci.GetOut() - b.GetParent().GetCount())

		//for j := 0; j <= emptyCr.Len()-1; j++ {
		//	weights[maxI][j] = 0.0
		//}

		//}
		//if b.GetParent().GetPosition() != ci.GetComment() {
		//	b.GetParent().SetPosition(b.GetParent().GetPosition() + ", " + ci.GetComment())
		//}
		b.GetParent().SetPosition(ci.GetComment())

		for i := 0; i <= emptyBl.GetItemsCount()-1; i++ {
			for j := 0; j <= emptyCr.Len()-1; j++ {
				if i == maxI || j == maxJ {
					weights[i][j] = 0.0
				}
			}
		}
		max = 0.0
		maxI = -1
		maxJ = -1
	}
	bl.SetState(balance.IsMergeV2)
	ch <- 100
}

func getWeight(val string, dict string) float64 {
	var res int
	str := dict
	for _, v := range val {
		pos := strings.Index(str, string(v))
		if pos != -1 {
			res++
			str = strings.Replace(str, string(v), "", 1)
		}
	}
	w1 := float64(res) / float64(len(val))

	res = 0
	str = val
	for _, v := range dict {
		pos := strings.Index(str, string(v))
		if pos != -1 {
			res++
			str = strings.Replace(str, string(v), "", 1)
		}
	}
	w2 := float64(res) / float64(len(dict))
	return math.Max(w1, w2)

}
