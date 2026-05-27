package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
	entity "github.com/Cameliuu/zenith/engine/entitiy"
	"github.com/Cameliuu/zenith/engine/overlay"
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
}

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
	drawColor := draw.Green
	for _, e := range entities {
		if e.IsAlive == 0 {
			continue
		}
		feet := overlay.Vector3{X: e.Position.X, Y: e.Position.Y, Z: e.Position.Z - 35}
		head := overlay.Vector3{X: feet.X, Y: feet.Y, Z: feet.Z + 60}

		feetSX, feetSY, feetOnScreen := overlay.WorldToScreen(viewMatrix, feet, overlay.ScreenW, overlay.ScreenH)
		_, headSY, headOnScreen := overlay.WorldToScreen(viewMatrix, head, overlay.ScreenW, overlay.ScreenH)

		if !feetOnScreen || !headOnScreen {
			continue
		}
		h := feetSY - headSY
		w := h / 2.5

		if modelToTeams[e.Model] != uint(g.LocalPlayerInfo.Team) {
			drawColor = draw.Red
		}
		draw.Box(hdc, win32.Rect{
			Left:   int32(feetSX - w/2),
			Top:    int32(headSY),
			Right:  int32(feetSX + w/2),
			Bottom: int32(feetSY),
		}, drawColor)
	}
}

func Run() {
	gameData, err := Init()

	if err != nil {
		log.Fatal("zenith: could not initialize engine: %w", err)
	}
	defer windows.CloseHandle(gameData.Handle)
	fmt.Println("zenith: engine initialized")
	go func() {
		for {
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
			time.Sleep(1 * time.Millisecond)
		}
	}()
	window.Run("Counter-Strike", gameData.Draw)
}
