package config

import (
	"fmt"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/zenith/input"
)

type Config struct {
	ESP    *ESP
	Aimbot *Aimbot
}

type ESP struct {
	Hotkey     int
	Toggled    bool
	EnemyColor draw.Color
	TeamColor  draw.Color
}
type Aimbot struct {
	Radius     float32
	Smoothness float32
	Hotkey     int
	toggled    bool
}

func (a *Aimbot) SetHotkey(value string) {
	a.Hotkey = input.KeyBinds[value]
}

func (a *Aimbot) GetHotkeyName() string {
	return input.KeyNameFromCode[a.Hotkey]
}

func (a *Aimbot) Toggle() {
	printValue := "OFF"
	if a.IsToggledOn() {
		printValue = "ON"
	}
	fmt.Println("zenith: Aimbot ", printValue)
	a.toggled = !a.toggled
}

func (a *Aimbot) IsToggledOn() bool {
	return a.toggled
}

func InitConfig() *Config {
	return &Config{
		ESP: &ESP{
			Hotkey:     input.VK_F2,
			Toggled:    false,
			EnemyColor: draw.Red,
			TeamColor:  draw.Green,
		},
		Aimbot: &Aimbot{
			Radius:     100,
			Smoothness: 1,
			Hotkey:     input.VK_F1,
		},
	}
}

func (esp *ESP) SetHotkey(value string) {
	esp.Hotkey = input.KeyBinds[value]
}
func (esp *ESP) Toggle() {
	printValue := "OFF"
	if esp.Toggled {
		printValue = "ON"
	}
	fmt.Println("zenith: ESP ", printValue)
	esp.Toggled = !esp.Toggled
}

func (esp *ESP) IsToggledOn() bool {
	return esp.Toggled == true
}
