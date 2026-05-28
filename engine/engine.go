package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/Cameliuu/veil/draw"
	"github.com/Cameliuu/veil/win32"
	"github.com/Cameliuu/veil/window"
	"github.com/Cameliuu/zenith/engine/aimbot"
	"github.com/Cameliuu/zenith/engine/config"
	entity "github.com/Cameliuu/zenith/engine/entitiy"
	"github.com/Cameliuu/zenith/engine/entitiy/utils"
	"github.com/Cameliuu/zenith/engine/overlay"
	"github.com/Cameliuu/zenith/input"
	"github.com/Cameliuu/zenith/memory"
	"golang.org/x/sys/windows"
)

type GameData struct {
	Handle                  windows.Handle
	HWModuleBaseAddress     uintptr
	ClientModuleBaseAddress uintptr
	Entities                []entity.PlayerInfo
	LocalPlayerInfo         entity.LocalPlayerInfo
	animFrames              [64]float32
	oldAnimFrames           [64]float32
	windowDetails           *window.Window
	Config                  *config.Config
	AimbotHasTarget         bool
}

const windowName = "Counter-Strike"

func Init(config *config.Config) (*GameData, error) {
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
		Config:                  config,
	}, nil
}

func (g *GameData) Draw(hdc uintptr) {
	screenW, screenH := g.windowDetails.Size()
	if g.Config.Aimbot.IsToggledOn() {
		radiusColor := draw.White
		if g.AimbotHasTarget {
			radiusColor = draw.Green
		}
		draw.Circle(hdc, screenW/2, screenH/2, int(g.Config.Aimbot.Radius), radiusColor)
	}
	if !g.Config.ESP.IsToggledOn() {
		return
	}
	entities := g.Entities
	viewMatrix := g.LocalPlayerInfo.ViewMatrix

	for i, e := range entities {
		if e.IsAlive == 0 || entities[i].IsStale {
			continue
		}

		drawColor := g.Config.ESP.EnemyColor
		isAllied := utils.IsEntityInTheSameTeam(e.Model, g.LocalPlayerInfo.Team)

		if isAllied {
			drawColor = g.Config.ESP.TeamColor
		}
		feet := utils.Vector3{X: e.Position.X, Y: e.Position.Y, Z: e.Position.Z - 40}
		head := utils.Vector3{X: feet.X, Y: feet.Y, Z: feet.Z + 60}

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
	gameData, err := Init(config)

	if err != nil {
		log.Fatal("zenith: could not initialize engine: %w", err)
	}
	defer windows.CloseHandle(gameData.Handle)

	fmt.Println("zenith: engine initialized")
	w, err := window.New(windowName)
	go func() {
		for {
			if input.IsKeyJustPressed(config.ESP.Hotkey) {
				config.ESP.Toggle()
			}
			if input.IsKeyJustPressed(config.Aimbot.Hotkey) {
				config.Aimbot.Toggle()
			}
			localPlayerInfo, err := entity.GetLocalPlayerInfo(gameData.Handle, gameData.ClientModuleBaseAddress, gameData.HWModuleBaseAddress)
			if err != nil {
				log.Fatal(err)
			}
			gameData.LocalPlayerInfo = localPlayerInfo
			entities, err := entity.GetEntities(gameData.Handle, gameData.HWModuleBaseAddress, &gameData.animFrames, &gameData.oldAnimFrames)
			if err != nil {
				log.Fatal(err)
			}
			gameData.Entities = entities
			if gameData.Config.Aimbot.IsToggledOn() {

				aimbot.Run(aimbot.AimbotInput{
					Entities:            gameData.Entities,
					LocalPlayerInfo:     gameData.LocalPlayerInfo,
					WindowDetails:       gameData.windowDetails,
					ViewAngles:          gameData.LocalPlayerInfo.ViewAngles,
					HWModuleBaseAddress: gameData.HWModuleBaseAddress,
					Handle:              gameData.Handle,
					HasTarget:           &gameData.AimbotHasTarget,
				}, config.Aimbot)
			}
			time.Sleep(16 * time.Millisecond)
		}
	}()
	if err != nil {
		log.Fatalf("zenith: could not open overlay window %w", err)
	}
	gameData.windowDetails = w
	w.Run(gameData.Draw)
}
