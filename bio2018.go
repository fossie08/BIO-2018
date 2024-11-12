package main

import (
	"fmt"
	"os"
	"strconv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var fixedRepaymentValue int = 50
var debt float64 = 100.00

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

// Allow only numeric input
func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' {
		e.Entry.TypedRune(r)
	}
}

// Prevent copy-paste for non-numeric input
func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func main() {
	// Set up the app and main window
	debtCalcApp := app.New()
	mainWindow := debtCalcApp.NewWindow("Debt Calculator")

	// Title label
	titleLabel := widget.NewLabel("Debt Calculator")
	titleLabel.Alignment = fyne.TextAlignCenter

	// Input fields for Interest and Repayment percentages
	interestPercentageEntry := newNumericalEntry()
	interestPercentageEntry.SetPlaceHolder("Interest %")

	repaymentPercentageEntry := newNumericalEntry()
	repaymentPercentageEntry.SetPlaceHolder("Repayment %")

	// Output label for total repaid amount
	totalRepaidLabel := widget.NewLabel("Total amount repaid:")
	totalRepaidOutput := widget.NewLabel("")

	// Variable to store monthly repayment data
	var monthlyReport []string

	// Save report button
	saveReportButton := widget.NewButtonWithIcon("Save Report", theme.DocumentSaveIcon(), func() {
		if len(monthlyReport) == 0 {
			dialog.NewInformation("No Data", "Please calculate repayment first.", mainWindow).Show()
			return
		}

		// Write report to file
		file, err := os.Create("repayment_report.txt")
		if err != nil {
			dialog.NewInformation("Error", "Could not save report.", mainWindow).Show()
			return
		}
		defer file.Close()

		for _, line := range monthlyReport {
			file.WriteString(line + "\n")
		}

		dialog.NewInformation("Report Saved", "Report has been saved successfully as repayment_report.txt.", mainWindow).Show()
	})

	// Calculate button
	calculateButton := widget.NewButtonWithIcon("Calculate Repayment", theme.ConfirmIcon(), func() {
		intInterestPercentage, _ := strconv.ParseInt(interestPercentageEntry.Text, 10, 64)
		intRepaymentPercentage, _ := strconv.ParseInt(repaymentPercentageEntry.Text, 10, 64)

		// Validate input values
		if intInterestPercentage < 0 || intInterestPercentage > 100 || intRepaymentPercentage < 0 || intRepaymentPercentage > 100 {
			dialog.NewInformation("Invalid Input", "Enter valid values between 0 and 100 for Interest and Repayment percentages.", mainWindow).Show()
		} else {
			repayed, report := calcLoop(debt, int(intInterestPercentage), int(intRepaymentPercentage))
			totalRepaidOutput.SetText("£" + strconv.FormatFloat(repayed, 'f', 2, 64))
			monthlyReport = report // Store report for saving
		}
	})

	// Layout
	layout := container.NewVBox(
		titleLabel,
		container.NewGridWithColumns(2,
			widget.NewLabel("Interest %:"), interestPercentageEntry,
			widget.NewLabel("Repayment %:"), repaymentPercentageEntry,
			totalRepaidLabel, totalRepaidOutput,
			calculateButton, saveReportButton,
		),
	)

	// Set main window content
	mainWindow.SetContent(layout)
	mainWindow.Resize(fyne.NewSize(300, 300))
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}

// Calculation function with monthly tracking
func calcLoop(debt float64, interest int, repayment int) (totalRepayed float64, report []string) {
	interestPerc := (float32(interest) / 100.00) + 1.00
	repaymentPerc := float32(repayment) / 100.00

	var repayed float64 = 0
	var oldDebt float64 = 0
	month := 1

	for {
		debt = debt * float64(interestPerc)
		var monthlyRepayment float64

		if (debt * float64(repaymentPerc)) > float64(fixedRepaymentValue) {
			oldDebt = debt
			debt = debt * float64(repaymentPerc)
			monthlyRepayment = oldDebt - debt
		} else {
			monthlyRepayment = float64(fixedRepaymentValue)
			debt = debt - monthlyRepayment
		}

		repayed += monthlyRepayment

		// Append monthly report
		reportLine := fmt.Sprintf("Month %d: Repaid £%.2f, Remaining Debt: £%.2f", month, monthlyRepayment, debt)
		report = append(report, reportLine)
		month++

		// Stop if debt is less than the fixed repayment value
		if debt < float64(fixedRepaymentValue) {
			debt = debt * float64(interestPerc)
			repayed += debt
			report = append(report, fmt.Sprintf("Month %d: Final Payment £%.2f, Remaining Debt: £0.00", month, debt))
			break
		}
	}

	return repayed, report
}
