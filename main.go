package main

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/zenith/engine"
	"github.com/Cameliuu/zenith/engine/config"
	"github.com/Cameliuu/zenith/input"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	cfg := config.InitConfig()

	go engine.Run(cfg)

	w := a.NewWindow("zenith")

	// ── ESP ──────────────────────────────────────────────────────────────
	espHotkeySelect := widget.NewSelect(input.KeyNames, cfg.ESP.SetHotkey)
	espHotkeySelect.SetSelected(input.KeyNameFromCode[cfg.ESP.Hotkey])

	colorButton := widget.NewButton("Enemy Color", func() {
		picker := dialog.NewColorPicker("Enemy Box Color", "Pick a color", func(c color.Color) {
			r, g, b, _ := c.RGBA()
			cfg.ESP.EnemyColor = draw.Color{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
			}
		}, w)
		picker.Show()
	})

	espSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## ESP"),
		container.NewGridWithColumns(2,
			widget.NewLabel("Toggle Key"),
			espHotkeySelect,
		),
		colorButton,
	)

	// ── Aimbot ───────────────────────────────────────────────────────────
	aimbotHotkeySelect := widget.NewSelect(input.KeyNames, cfg.Aimbot.SetHotkey)
	aimbotHotkeySelect.SetSelected(input.KeyNameFromCode[cfg.Aimbot.Hotkey])

	radiusLabel := widget.NewLabel("FOV Radius: " + strconv.Itoa(int(cfg.Aimbot.Radius)))
	radiusSlider := widget.NewSlider(50, 500)
	radiusSlider.SetValue(float64(cfg.Aimbot.Radius))
	radiusSlider.OnChanged = func(v float64) {
		cfg.Aimbot.Radius = float32(v)
		radiusLabel.SetText("FOV Radius: " + strconv.Itoa(int(v)))
	}

	smoothLabel := widget.NewLabel("Smoothness: " + strconv.FormatFloat(float64(cfg.Aimbot.Smoothness), 'f', 1, 32))
	smoothSlider := widget.NewSlider(1, 20)
	smoothSlider.SetValue(float64(cfg.Aimbot.Smoothness))
	smoothSlider.OnChanged = func(v float64) {
		cfg.Aimbot.Smoothness = float32(v)
		smoothLabel.SetText("Smoothness: " + strconv.FormatFloat(v, 'f', 1, 64))
	}

	aimbotSection := container.NewVBox(
		widget.NewRichTextFromMarkdown("## Aimbot"),
		container.NewGridWithColumns(2,
			widget.NewLabel("Toggle Key"),
			aimbotHotkeySelect,
		),
		radiusLabel,
		radiusSlider,
		smoothLabel,
		smoothSlider,
	)

	// ── Tabs ─────────────────────────────────────────────────────────────
	tabs := container.NewAppTabs(
		container.NewTabItem("ESP", espSection),
		container.NewTabItem("Aimbot", aimbotSection),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(450, 350))
	w.ShowAndRun()
}
