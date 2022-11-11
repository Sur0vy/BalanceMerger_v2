package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Списание V2.0")
	w.Resize(fyne.NewSize(500, 250))

	//журнал
	lblJournal := widget.Label{Text: "Журнал:"}
	entJournal := widget.Entry{PlaceHolder: "Файл журнала (*.xlsx)"}
	btnJournal := widget.Button{Icon: theme.FolderOpenIcon()}

	//баланс
	lblBalance := widget.Label{Text: "Баланс:"}
	entBalance := widget.Entry{PlaceHolder: "Файл журнала (*.xlsx)"}
	btnBalance := widget.Button{Icon: theme.FolderOpenIcon()}

	//карточка
	lblCard := widget.Label{Text: "Карточка счета:"}
	entCard := widget.Entry{PlaceHolder: "Файл журнала (*.xlsx)"}
	btnCard := widget.Button{Icon: theme.FolderOpenIcon()}

	//Результат
	btnProcess := widget.Button{Text: "Обработать", Icon: theme.ConfirmIcon()}
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
					&btnJournal,
					&btnBalance,
					&btnCard,
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
	w.ShowAndRun()
}
