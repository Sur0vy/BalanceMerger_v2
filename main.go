package main

import (
	"BM/Models"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"path/filepath"
	"time"
)

var src Models.Sources

var ap = app.New()

var mainWindow fyne.Window

var currentDir string

var pbJournal = widget.ProgressBar{
	Min: 0,
	Max: 100,
}
var pbBalance = widget.ProgressBar{
	Min: 0,
	Max: 100,
}
var pbCard = widget.ProgressBar{
	Min: 0,
	Max: 100,
}
var pbProcess = widget.ProgressBar{
	Min: 0,
	Max: 100,
}

var lbMessage = widget.Label{
	Alignment: fyne.TextAlignCenter,
	Wrapping:  fyne.TextWrapWord,
}

var cbAutoSave = widget.Check{
	Text: "Автосохранение",
}

var btnSave = widget.Button{
	Text:          "Сохранить...",
	Icon:          theme.DocumentSaveIcon(),
	Importance:    0,
	Alignment:     0,
	IconPlacement: 0,
}

var onSaveTapped = func() {
	dlgSave := dialog.NewFileSave(func(r fyne.URIWriteCloser, _ error) {
		if r != nil {
			absMediaFile, _ := filepath.Abs(r.URI().Path())
			SaveMergedFile(absMediaFile)
			currentDir = filepath.Dir(absMediaFile)
			btnSave.Disable()
		}
	}, mainWindow)

	if currentDir != "" {
		mfileURI := storage.NewFileURI(currentDir)
		mfileLister, _ := storage.ListerForURI(mfileURI)
		dlgSave.SetLocation(mfileLister)
	}

	dlgSave.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
	dlgSave.SetFileName(src.GetOutFileName(true))
	dlgSave.Show()
}

func main() {
	//src.Card = "/Users/Ses/Sur0vy/краснополь списание.xlsx"
	run()
}

func run() {
	mainWindow = ap.NewWindow("Списание V2.0")
	initGUI(mainWindow)
	mainWindow.ShowAndRun()
}

func processDocuments() {
	pbProcess.SetValue(0)
	pbCard.SetValue(0)
	pbBalance.SetValue(0)
	pbJournal.SetValue(0)

	ans := make(chan response)
	go StartProcess(src, ans)
	processDone := false

	for {
		if processDone {
			break
		}
		time.Sleep(10 * time.Microsecond)
		val, ok := <-ans
		if ok == false {
			lbMessage.SetText("Непредвиденная ошибка обработки!")
			break // exit break loop
		} else {
			switch val.Step {
			case RiJournal:
				pbJournal.SetValue(val.Progress)
			case RiBalance:
				pbBalance.SetValue(val.Progress)
			case RiCard:
				pbCard.SetValue(val.Progress)
			default:
				pbProcess.SetValue(val.Progress)
				if val.Progress == 100 {
					processDone = true
					break
				}
			}
			if val.Status == false {
				switch val.Step {
				case RiJournal:
					lbMessage.SetText("Ошибка при открытии файла журнала, возможно файл поврежден " +
						"или данные в некорректном формате.")
				case RiBalance:
					lbMessage.SetText("Ошибка при открытии файла баланса, возможно файл поврежден " +
						"или данные в некорректном формате.")
				case RiCard:
					lbMessage.SetText("Ошибка при открытии файла карточки счета, возможно файл поврежден " +
						"или данные в некорректном формате.")
				default:
					lbMessage.SetText("Ошибка при обработке данных.")
				}
				break
			}
		}
	}

	if processDone {
		lbMessage.SetText("Обработка данных выполнена успешно!")
		if cbAutoSave.Checked {
			SaveMergedFile(src.GetOutFileName(false))
		} else {
			btnSave.Enable()
		}
	}
}

func initGUI(w fyne.Window) {
	w.Resize(fyne.NewSize(600, 500))

	//журнал
	lblJournal := widget.Label{Text: "Журнал:"}
	entJournal := widget.Entry{
		PlaceHolder: "Файл журнала (*.xlsx)",
		Wrapping:    fyne.TextTruncate,
	}
	btnJournal := widget.NewButton("", func() {
		dlgJournal := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				if r != nil {
					src.Journal, _ = filepath.Abs(r.URI().Path())
					currentDir = filepath.Dir(src.Journal)
					entJournal.SetText(src.Journal)
				}
			}, w)
		if currentDir != "" {
			mfileURI := storage.NewFileURI(currentDir)
			mfileLister, _ := storage.ListerForURI(mfileURI)
			dlgJournal.SetLocation(mfileLister)
		}
		dlgJournal.SetFilter(
			storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlgJournal.Show()
	})
	btnJournal.SetIcon(theme.FolderOpenIcon())

	//баланс
	lblBalance := widget.Label{Text: "Баланс:"}
	entBalance := widget.Entry{
		PlaceHolder: "Файл баланса (*.xlsx)",
		Wrapping:    fyne.TextTruncate,
	}
	btnBalance := widget.NewButton("", func() {
		dlgBalance := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				if r != nil {
					src.Balance, _ = filepath.Abs(r.URI().Path())
					currentDir = filepath.Dir(src.Balance)
					entBalance.SetText(src.Balance)
				}
			}, w)
		if currentDir != "" {
			mfileURI := storage.NewFileURI(currentDir)
			mfileLister, _ := storage.ListerForURI(mfileURI)
			dlgBalance.SetLocation(mfileLister)
		}
		dlgBalance.SetFilter(
			storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlgBalance.Show()
	})
	btnBalance.SetIcon(theme.FolderOpenIcon())

	//карточка
	lblCard := widget.Label{Text: "Карточка счета:"}
	entCard := widget.Entry{
		PlaceHolder: "Файл карточки счета (*.xlsx)",
		Wrapping:    fyne.TextTruncate,
	}
	btnCard := widget.NewButton("", func() {
		dlgCard := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				if r != nil {
					src.Card, _ = filepath.Abs(r.URI().Path())
					currentDir = filepath.Dir(src.Card)
					entBalance.SetText(src.Card)
				}
			}, w)
		if currentDir != "" {
			mfileURI := storage.NewFileURI(currentDir)
			mfileLister, _ := storage.ListerForURI(mfileURI)
			dlgCard.SetLocation(mfileLister)
		}

		dlgCard.SetFilter(
			storage.NewExtensionFileFilter([]string{".xlsx"}))

		dlgCard.Show()
	})
	btnCard.SetIcon(theme.FolderOpenIcon())

	//Результат
	btnProcess := widget.Button{Text: "Обработать", Icon: theme.ConfirmIcon()}
	btnProcess.OnTapped = processDocuments
	btnExt := widget.Button{Text: "Закрыть"}
	btnExt.OnTapped = func() {
		w.Close()
	}

	w.SetContent(container.NewBorder(
		nil,
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewVBox(
				&btnSave,
				&btnExt,
			),
			container.NewVBox(),
		),
		nil,
		nil,
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				container.NewVBox(
					&lblJournal,
					&lblBalance,
					&lblCard,
				),
				container.NewVBox(
					btnJournal,
					btnBalance,
					btnCard,
				),
				container.NewVBox(
					&entJournal,
					&entBalance,
					&entCard,
				),
			),
			container.NewBorder(
				nil,
				nil,
				nil,
				&btnProcess,
				&cbAutoSave,
			),
			container.NewBorder(
				nil,
				nil,
				container.NewVBox(
					widget.NewLabel("Журнал:"),
					widget.NewLabel("Баланс:"),
					widget.NewLabel("Карточка счета:"),
					widget.NewLabel("Обработка:"),
				),
				nil,
				container.NewVBox(
					&pbJournal,
					&pbBalance,
					&pbCard,
					&pbProcess,
					&lbMessage,
				),
			),
		),
	))
	btnSave.OnTapped = onSaveTapped
	btnSave.Disable()
}
