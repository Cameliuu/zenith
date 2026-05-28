package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
	"github.com/Cameliuu/zenith/engine/config"
	entity "github.com/Cameliuu/zenith/engine/entitiy"
	"github.com/Cameliuu/zenith/engine/overlay"
	"github.com/Cameliuu/zenith/input"
	"github.com/Cameliuu/zenith/memory"
	"golang.org/x/sys/windows"
)

var modelToTeams = map[string]uint{
	"gign":     2,
	"gsg":      2,
	"sas":      2,
	"urban":    2,
	"terror":   1,
	"leet":     1,
	"arctic":   1,
	"guerilla": 1,
}

type GameData struct {
	Handle                  windows.Handle
	HWModuleBaseAddress     uintptr
	ClientModuleBaseAddress uintptr
	Entities                []entity.PlayerInfo
	LocalPlayerInfo         entity.LocalPlayerInfo
	animFrames              [64]float32
	oldAnimFrames           [64]float32
	windowDetails           *window.Window
}

const windowName = "Counter-Strike"

func Init() (*GameData, error) {
	pid, err := memory.FindProcessIdByName("cstrike.exe")

	if err != nil {
		return nil, err
	}

	hwModuleBase, err := memory.FindModuleBaseAddressByPid(pid, "hw.dll")

	if err != nil {
		return nil, err
	}
	clientModuleBase, err := memory.FindModuleBaseAddressByPid(pid, "client.dll")

	if err != nil {
		return nil, err
	}
	handle, err := memory.OpenHandleToProcessByPID(pid)

	if err != nil {
		return nil, err
	}

	return &GameData{
		Handle:                  handle,
		HWModuleBaseAddress:     hwModuleBase,
		ClientModuleBaseAddress: clientModuleBase,
	}, nil
}

func (g *GameData) Draw(hdc uintptr) {
	entities := g.Entities
	viewMatrix := g.LocalPlayerInfo.ViewMatrix

	for i, e := range entities {
		curr := e.AnimFrame
		old := g.animFrames[i]
		oldOld := g.oldAnimFrames[i]
		screenW, screenH := g.windowDetails.Size()
		if curr == old && curr == oldOld {
			entities[i].IsStale = true
		} else {
			entities[i].IsStale = false
		}

		g.oldAnimFrames[i] = old
		g.animFrames[i] = curr

		if e.IsAlive == 0 || entities[i].IsStale {
			continue
		}

		drawColor := draw.Green

		team, ok := modelToTeams[e.Model]
		if !ok {
			drawColor = draw.Blue
		} else if team != uint(g.LocalPlayerInfo.Team) {
			drawColor = draw.Red
		}

		feet := overlay.Vector3{X: e.Position.X, Y: e.Position.Y, Z: e.Position.Z - 40}
		head := overlay.Vector3{X: feet.X, Y: feet.Y, Z: feet.Z + 60}

		feetSX, feetSY, feetOnScreen := overlay.WorldToScreen(viewMatrix, feet, screenW, screenH)
		headSX, headSY, headOnScreen := overlay.WorldToScreen(viewMatrix, head, screenW, screenH)

		if !feetOnScreen || !headOnScreen {
			continue
		}

		h := feetSY - headSY
		w := h / 3
		draw.Box3D(hdc, win32.Rect{
			Left:   int32(feetSX - w/2),
			Top:    int32(headSY),
			Right:  int32(feetSX + w/2),
			Bottom: int32(feetSY),
		}, drawColor, 8)

		draw.TextOut(hdc, e.Name, int32(headSX-20), int32(headSY-20), drawColor)
	}
}

func Run(config *config.Config) {
	gameData, err := Init()

	if err != nil {
		log.Fatal("zenith: could not initialize engine: %w", err)
	}
	defer windows.CloseHandle(gameData.Handle)
	fmt.Println("zenith: engine initialized")
	go func() {
		for {
			if input.IsKeyJustPressed(config.ESP.Hotkey) {
				config.ESP.Toggle()
			}
			if config.ESP.IsToggledOn() {
				localPlayerInfo, err := entity.GetLocalPlayerInfo(gameData.Handle, gameData.ClientModuleBaseAddress, gameData.HWModuleBaseAddress)
				if err != nil {
					log.Fatal(err)
				}
				gameData.LocalPlayerInfo = localPlayerInfo
				entities, err := entity.GetEntities(gameData.Handle, gameData.HWModuleBaseAddress)
				if err != nil {
					log.Fatal(err)
				}
				gameData.Entities = entities
				time.Sleep(16 * time.Millisecond)
			}
		}
	}()
	w, err := window.New(windowName)
	if err != nil {
		log.Fatalf("zenith: could not open overlay window %w")
	}
	gameData.windowDetails = w
	fmt.Println(w.Width(), w.Height())
	w.Run(gameData.Draw)
}
