package main

import (
	"io"
    "strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
    EditWidget *widget.Entry
    PreviewWidget *widget.RichText
    SaveMenuItem *fyne.MenuItem
    CurrentFile fyne.URI
}

var conf config
var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

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
    openMenuItem := fyne.NewMenuItem("Open", app.openFunc(w)) 
    saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(w))
    app.SaveMenuItem = saveMenuItem
    app.SaveMenuItem.Disabled = true
    saveAsMenuItem := fyne.NewMenuItem("Save as", app.saveAsFunc(w))
    fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
    menu := fyne.NewMainMenu(fileMenu)
    w.SetMainMenu(menu)
}

func (app *config) saveAsFunc(w fyne.Window) func() {
    return func() {
        saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
            if err != nil {
                dialog.ShowError(err, w)
                return
            }

            if write == nil {
                return
            }
            
            if !strings.HasSuffix(strings.ToLower(write.URI().Name()), ".md") {
                dialog.ShowInformation("Error", "The name of the file must end with .md extension!", w)
                return
            }
            write.Write([]byte(app.EditWidget.Text))
            app.CurrentFile = write.URI()
            defer write.Close()
            w.SetTitle(w.Title() + " - " + write.URI().Name())
            app.SaveMenuItem.Disabled = false
        }, w)
        saveDialog.SetFileName("untitled.md")
        saveDialog.SetFilter(filter)
        saveDialog.Show()
    }
}


func (app *config) openFunc(w fyne.Window) func() {
    return func(){
        openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
            if err != nil {
                dialog.ShowError(err, w)
                return
            }

            if read == nil {
                return
            }
            
            data, err := io.ReadAll(read)
            if err != nil {
                dialog.ShowError(err, w)
                return
            }

            app.EditWidget.SetText(string(data))
            app.CurrentFile = read.URI()
            w.SetTitle(w.Title() + " - " + read.URI().Name())
        }, w)
        openDialog.SetFilter(filter)
        openDialog.Show()
    }
}

func (app *config) saveFunc(w fyne.Window) func() {
    return func() {
        if app.CurrentFile != nil {
            write, err := storage.Writer(app.CurrentFile)
            if err != nil {
                dialog.ShowError(err, w)
                return
            }

            write.Write([]byte(app.EditWidget.Text))
            defer write.Close()
        }
    }
}
