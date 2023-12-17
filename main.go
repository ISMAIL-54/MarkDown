package main
import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
    "fyne.io/fyne/v2/container"
)

type config struct {
    EditWidget *widget.Entry
    PreviewWidget *widget.RichText
}

var conf config

func main(){
    // create an app
    a := app.New()

    // create a window
    w := a.NewWindow("MarkDown")

    // get the user interface
    edit, preview := conf.makeUI()

    conf.createMenuItems(w)
    // set the content of the window
    w.SetContent(container.NewHSplit(edit, preview))

    w.Resize(fyne.Size{Width: 800, Height: 500})
    w.CenterOnScreen()
    w.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
    edit := widget.NewMultiLineEntry()
    preview := widget.NewRichTextFromMarkdown("")
    app.EditWidget = edit
    app.PreviewWidget = preview

    edit.OnChanged = preview.ParseMarkdown

    return edit, preview
}

func (app *config) createMenuItems(w fyne.Window) {
    openMenuItem := fyne.NewMenuItem("Open", func(){})
    saveMenuItem := fyne.NewMenuItem("Save", func(){})
    saveAsMenuItem := fyne.NewMenuItem("Save as", func(){})
    fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
    menu := fyne.NewMainMenu(fileMenu)
    w.SetMainMenu(menu)
}
