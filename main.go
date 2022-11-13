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
)

var src Models.Sources

func main() {
	run()
}

func run() {
	a := app.New()
	w := a.NewWindow("Списание V2.0")
	initGUI(w)
	w.ShowAndRun()
}

func processDocuments() {
	//todo: проверить существование файлов
	StartProcess(src)
}

func initGUI(w fyne.Window) {
	w.Resize(fyne.NewSize(600, 500))

	//журнал
	lblJournal := widget.Label{Text: "Журнал:"}
	entJournal := widget.Entry{PlaceHolder: "Файл журнала (*.xlsx)"}
	btnJournal := widget.NewButton("", func() {
		dlgJournal := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				src.Journal = r.URI().Path()
				entJournal.SetText(src.Journal)
			}, w)
		dlgJournal.SetFilter(
			storage.NewExtensionFileFilter([]string{".xlsx", ".csv"}))
		dlgJournal.Show()
	})
	btnJournal.SetIcon(theme.FolderOpenIcon())

	//баланс
	lblBalance := widget.Label{Text: "Баланс:"}
	entBalance := widget.Entry{PlaceHolder: "Файл журнала (*.xlsx)"}
	btnBalance := widget.NewButton("", func() {
		dlgBalance := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				src.Balance = r.URI().Path()
				entBalance.SetText(src.Balance)
			}, w)
		dlgBalance.SetFilter(
			storage.NewExtensionFileFilter([]string{".xlsx"}))
		dlgBalance.Show()
	})
	btnBalance.SetIcon(theme.FolderOpenIcon())

	//карточка
	lblCard := widget.Label{Text: "Карточка счета:"}
	entCard := widget.Entry{PlaceHolder: "Файл журнала (*.xlsx)"}
	btnCard := widget.NewButton("", func() {
		dlgCard := dialog.NewFileOpen(
			func(r fyne.URIReadCloser, _ error) {
				src.Card = r.URI().Path()
				entCard.SetText(src.Card)
			}, w)
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
	progress := widget.ProgressBar{}

	w.SetContent(container.NewBorder(
		nil,
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewVBox(
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
				container.NewVBox(
					&btnProcess,
				),
				container.NewVBox(
					&progress,
				),
			),
		),
	))
}
