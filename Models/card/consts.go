package card

import "strings"

const (
	fldDoc  string = "Документ"
	fldDocD string = "Аналитика Дт"
	fldDocC string = "Аналитика Кт"
	fldCntD string = "Дебет"
	fldCntC string = "Кредит"

	repDocB   string = "вх.д."
	repDocE   string = "от"
	repCredit string = "<...>"

	tryCnt = 10
)

var Elements = [...]string{
	"провод монтажный",
	"стойка для печатных плат",
	"стойка",
	"планка",
	"панель",
	"пластина",
	"прокладка",
	"шпилька",
	"корпус",
	"блок питания",
	"печатная плата",
	"штекер",
	"разъём",
	"разъем",
	"гнездо разъём",
	"транзистор",
	"кварцевый резонатор",
	"резонатор",
	"чип резистор",
	"резистор",
	"микросхема",
	"танталовый конденсатор",
	"конденсатор",
	"реле",
	"феррит",
	"предохранитель",
	"светодиод желтый",
	"светодиод зелёный",
	"светодиод",
	"диод",
	"кнопка тактовая",
	"кнопка",
	"сердечник ферритовый",
	"хомут",
	"прижим",
	"уголок",
	"винт",
	"основание",
	"кронштейн",
	"диск",
}

func RemeoveElements(val string) string {
	tmp := strings.ToLower(val)
	for _, word := range Elements {
		if strings.Contains(tmp, word) {
			tmp = strings.ReplaceAll(tmp, word, "")
			break
		}
	}
	tmp = strings.ReplaceAll(tmp, " ", "")
	tmp = strings.ReplaceAll(tmp, ".", "")
	tmp = strings.ReplaceAll(tmp, ",", "")
	tmp = strings.ReplaceAll(tmp, "-", "")
	return tmp
}
