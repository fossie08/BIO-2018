package main
//import libraries
import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	//setup the app and main window
	deptCalcApp := app.New()
	mainWindow := deptCalcApp.NewWindow("Dept Calculator")

	
	//set the content and layout of the main window
	content := container.NewBorder(nil, nil, nil, nil, nil)
	//apply content to main window
	mainWindow.SetContent(content)
	//center on the screen and launch the application
	mainWindow.CenterOnScreen()
	mainWindow.ShowAndRun()
}