package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Cameliuu/zenith/engine"
	"github.com/Cameliuu/zenith/engine/config"
	"github.com/Cameliuu/zenith/input"
)

func main() {
	a := app.New()
	cfg := config.InitConfig()

	go engine.Run(cfg)

	w := a.NewWindow("zenith")
	espSelect := widget.NewSelect(input.KeyNames, cfg.ESP.SetHotkey)
	espSelect.SetSelected(input.KeyNameFromCode[cfg.ESP.Hotkey])
	content := container.NewVBox(
		widget.NewLabel("ESP"),
		espSelect,
	)
	w.SetContent(content)
	w.ShowAndRun()
}
