package config

import (
	"fmt"

	"github.com/Cameliuu/zenith/input"
)

type Config struct {
	ESP *ESP
}

type ESP struct {
	Hotkey  int
	Toggled bool
}

func InitConfig() *Config {
	return &Config{
		ESP: &ESP{
			Hotkey:  input.VK_F2,
			Toggled: false},
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
