package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type AccessRecord struct {
	Sequence int
	IsHit    bool
	SetIndex int
}

var records []AccessRecord

func runTestWithAnimation(name string, sequence []int, table *widget.Table, progress *widget.ProgressBar, updateUI func()) map[string]float64 {
	var (
		cache = NewCache()
		total = len(sequence)
	)
	records = make([]AccessRecord, 0)

	for i, seqNum := range sequence {
		isHit, setIndex := cache.Access(seqNum)
		records = append(records, AccessRecord{
			Sequence: seqNum,
			IsHit:    isHit,
			SetIndex: setIndex,
		})
		progress.SetValue(float64(i+1) / float64(total))
		updateUI()
		time.Sleep(10 * time.Millisecond)
	}
	return cache.GetStats()
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Cache Simulator")

	// Title
	title := canvas.NewText("Cache Access Simulator", theme.PrimaryColor())
	title.TextSize = 24
	title.TextStyle.Bold = true
	titleContainer := container.NewCenter(title)

	// Settings input
	cacheSize := widget.NewEntry()
	cacheSize.SetText("32")
	settingsCard := widget.NewCard(
		"Parameters",
		"Configure cache settings",
		container.NewGridWithColumns(2,
			widget.NewLabel("Cache Blocks:"),
			cacheSize,
		),
	)

	// Results display
	results := widget.NewTextGrid()
	resultsScroll := container.NewScroll(results)
	resultsScroll.SetMinSize(fyne.NewSize(300, 300))

	resultsCard := widget.NewCard(
		"Results",
		"Test Case:",
		resultsScroll,
	)

	// Table settings
	colWidth := float32(120.0)

	// Access sequence table
	accessTable := widget.NewTable(
		func() (int, int) {
			return len(records), 4
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("000")
			label.Alignment = fyne.TextAlignCenter
			label.TextStyle = fyne.TextStyle{Monospace: true}
			return container.NewCenter(label)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			label := container.Objects[0].(*widget.Label)
			label.Alignment = fyne.TextAlignCenter

			if i.Row == -1 {
				headers := []string{"Seq", "Hit", "Miss", "Set"}
				label.SetText(headers[i.Col])
				label.TextStyle = fyne.TextStyle{Bold: true}
				return
			}

			if i.Row < len(records) {
				record := records[i.Row]
				switch i.Col {
				case 0:
					label.SetText(fmt.Sprintf("%d", record.Sequence))
				case 1:
					if record.IsHit {
						label.SetText("✓")
						label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
					} else {
						label.SetText("")
					}
				case 2:
					if !record.IsHit {
						label.SetText("✗")
						label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
					} else {
						label.SetText("")
					}
				case 3:
					label.SetText(fmt.Sprintf("%d", record.SetIndex))
					label.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
				}
			} else {
				label.SetText("")
			}
		})

	// Set table column widths
	for i := 0; i < 4; i++ {
		accessTable.SetColumnWidth(i, colWidth)
	}

	// Table header
	headerLabels := []string{"Seq", "Hit", "Miss", "Set"}
	headerContainers := make([]fyne.CanvasObject, 4)
	for i, text := range headerLabels {
		label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		headerContainers[i] = container.NewHBox(
			label,
			layout.NewSpacer(),
		)
	}
	// tableHeader := container.NewGridWithColumns(4, headerContainers...)

	// Table scroll container
	tableScroll := container.NewScroll(accessTable)
	tableScroll.SetMinSize(fyne.NewSize(colWidth*4+20, 400))

	// 添加进度条
	progress := widget.NewProgressBar()
	progress.Min = 0
	progress.Max = 1

	visualizationCard := widget.NewCard(
		"Visualization",
		"Process table (Seq , Hit , Miss , Set) ",
		container.NewVBox(
			progress,
			// container.NewPadded(tableHeader),
			container.NewMax(tableScroll),
		),
	)

	// Buttons
	seqButton := widget.NewButtonWithIcon("Sequential Test", theme.MediaPlayIcon(), func() {
		results.SetText("")
		records = make([]AccessRecord, 0)
		myWindow.Canvas().Refresh(accessTable)
		progress.SetValue(0) // 重置进度条
		n, _ := strconv.Atoi(cacheSize.Text)
		//seq := GenerateSameSetSequence(1)
		//seq2 := GenerateSameSetSequence(1)
		seq := GenerateSequential(n) //GenerateSequential
		stats := runTestWithAnimation("Sequential", seq, accessTable, progress, func() {
			myWindow.Canvas().Refresh(accessTable)
		})
		displayResults(results, "Sequential", stats)
	})

	randButton := widget.NewButtonWithIcon("Random", theme.MediaPlayIcon(), func() {
		results.SetText("")
		records = make([]AccessRecord, 0)
		myWindow.Canvas().Refresh(accessTable)
		progress.SetValue(0) // 重置进度条
		n, _ := strconv.Atoi(cacheSize.Text)
		seq := GenerateRandom(n)
		stats := runTestWithAnimation("Random", seq, accessTable, progress, func() {
			myWindow.Canvas().Refresh(accessTable)
		})
		displayResults(results, "Random", stats)
	})

	midRepeatButton := widget.NewButtonWithIcon("Mid-Repeat", theme.MediaPlayIcon(), func() {
		results.SetText("")
		records = make([]AccessRecord, 0)
		myWindow.Canvas().Refresh(accessTable)
		progress.SetValue(0) // 重置进度条
		n, _ := strconv.Atoi(cacheSize.Text)
		seq := GenerateMidRepeat(n)
		stats := runTestWithAnimation("Mid-Repeat", seq, accessTable, progress, func() {
			myWindow.Canvas().Refresh(accessTable)
		})
		displayResults(results, "Mid-Repeat", stats)
	})

	buttonsCard := widget.NewCard(
		"Case Options",
		"Select sequence type",
		container.NewGridWithColumns(3,
			seqButton,
			randButton,
			midRepeatButton,
		),
	)

	// Layout
	leftPanel := container.NewVBox(
		settingsCard,
		buttonsCard,
		resultsCard,
	)

	rightPanel := container.NewVBox(
		visualizationCard,
	)

	mainContent := container.NewHSplit(
		container.NewPadded(leftPanel),
		container.NewPadded(rightPanel),
	)
	mainContent.SetOffset(0.25)

	content := container.NewVBox(
		titleContainer,
		widget.NewSeparator(),
		container.NewPadded(mainContent),
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.ShowAndRun()
}

func displayResults(grid *widget.TextGrid, name string, stats map[string]float64) {
	result := fmt.Sprintf("\n statistics:  [%s]\n\n"+
		"📊 Basic Statistics:\n"+
		"  • Access Count: %.0f\n"+
		"  • Hit Count: %.0f\n"+
		"  • Miss Count: %.0f\n\n"+
		"📈 Performance Metrics:\n"+
		"  • Hit Rate: %.2f%%\n"+
		"  • Miss Rate: %.2f%%\n"+
		"  • Average Access Time: %.2f ns\n"+
		"  • Total Access Time: %.2f ns\n",
		name,
		stats["accessCount"],
		stats["hitCount"],
		stats["missCount"],
		stats["hitRate"]*100,
		stats["missRate"]*100,
		stats["avgAccessTime"],
		stats["totalAccessTime"])

	grid.SetText(result)
}
