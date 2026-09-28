package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed asset/icon.jpeg
var aByteIcon []byte

func main() {
	oIcon := fyne.NewStaticResource("icon.jpeg", aByteIcon)

	a := app.New()
	a.SetIcon(oIcon)

	w := a.NewWindow("Desktop Fyne")
	w.SetIcon(oIcon)
	w.ShowAndRun()
}
