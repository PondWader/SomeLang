package profiler

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func OpenProfilerResultsViewer() {
	application := app.New()
	window := application.NewWindow("Profiler Results Viewer")

	textWidget := widget.NewLabel("No file selected.")

	openResultsBtn := widget.NewButton("Open profiler results file", func() {
		filePath, err := dialog.File().Filter("CSV results file", "csv").Load()
		if err != nil {
			return
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			dialog.Message("Failed to read file content: %s", err).Error()
			return
		}
		result, err := ParseCsv(string(content))
		if err != nil {
			dialog.Message("Failed to parse results: %s", err).Error()
			return
		}
		textWidget.SetText("Viewing " + filepath.Base(filePath) + "\n" + result.ToSortedStringFormat(0))
	})

	content := container.New(
		layout.NewVBoxLayout(),
		openResultsBtn,
		textWidget,
	)

	window.Resize(fyne.NewSize(1200, 800))
	window.SetContent(content)
	window.ShowAndRun()
}
