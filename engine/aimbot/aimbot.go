package aimbot

import (
	"fmt"
	"math"

	"github.com/Cameliuu/veil/window"
	"github.com/Cameliuu/zenith/engine/config"
	entity "github.com/Cameliuu/zenith/engine/entitiy"
	"github.com/Cameliuu/zenith/engine/entitiy/utils"
	"github.com/Cameliuu/zenith/engine/offsets"
	"github.com/Cameliuu/zenith/engine/overlay"
	"github.com/Cameliuu/zenith/input"
	"github.com/Cameliuu/zenith/memory"
	"golang.org/x/sys/windows"
)

type AimbotInput struct {
	Entities            []entity.PlayerInfo
	LocalPlayerInfo     entity.LocalPlayerInfo
	WindowDetails       *window.Window
	ViewAngles          overlay.ViewAngles
	Handle              windows.Handle
	HWModuleBaseAddress uintptr
	HasTarget           *bool
}

func findClosestPlayer(entities []entity.PlayerInfo, w, h int, localPlayer *entity.LocalPlayerInfo, radius float32) *entity.PlayerInfo {
	bestDist := float32(math.MaxFloat32)
	var closestEnemy *entity.PlayerInfo
	for _, e := range entities {
		if utils.IsEntityInTheSameTeam(e.Model, localPlayer.Team) {
			continue
		}
		centerW, centerH := w/2, h/2
		enemyW, enemyH, onScreen := overlay.WorldToScreen(localPlayer.ViewMatrix, e.Position, w, h)
		if !onScreen || e.IsAlive == 0 || e.IsStale {
			continue
		}

		dist := utils.DistanceBetweenTwoPoints(utils.Vector2{
			X: float32(centerW),
			Y: float32(centerH)},
			utils.Vector2{X: enemyW,
				Y: enemyH})
		if dist >= radius*radius {
			continue
		}
		if dist < bestDist {
			bestDist = dist
			closestEnemy = &e
		}

	}

	return closestEnemy
}

func NormalizeAngle(angle float32) float32 {
	for angle > 180 {
		angle -= 360
	}
	for angle < -180 {
		angle += 360
	}
	return angle
}

func Smooth(current, target, factor float32) float32 {
	delta := NormalizeAngle(target - current)
	return current + delta/factor
}
func CalcAngle(origin, target utils.Vector3) overlay.ViewAngles {

	head := utils.Vector3{
		X: target.X,
		Y: target.Y,
		Z: target.Z + 20,
	}
	dx := head.X - origin.X
	dy := head.Y - origin.Y
	dz := head.Z - origin.Z

	hyp := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	pitch := float32(math.Atan2(float64(-dz), float64(hyp)) * (180.0 / math.Pi))
	yaw := float32(math.Atan2(float64(dy), float64(dx)) * (180.0 / math.Pi))

	return overlay.ViewAngles{Pitch: pitch, Yaw: yaw}
}

func Run(aimbotInput AimbotInput, cfg *config.Aimbot) {
	*aimbotInput.HasTarget = false
	w, h := aimbotInput.WindowDetails.Size()
	memory.Write(aimbotInput.Handle, aimbotInput.HWModuleBaseAddress+0x122E324, utils.Vector3{X: 0, Y: 0, Z: 0})
	closestPlayer := findClosestPlayer(aimbotInput.Entities, w, h, &aimbotInput.LocalPlayerInfo, cfg.Radius)
	if closestPlayer == nil {
		return
	}
	*aimbotInput.HasTarget = true
	angles := CalcAngle(aimbotInput.LocalPlayerInfo.Origin, closestPlayer.Position)
	smoothPitch := Smooth(aimbotInput.LocalPlayerInfo.ViewAngles.Pitch, angles.Pitch, cfg.Smoothness)
	smoothYaw := Smooth(aimbotInput.LocalPlayerInfo.ViewAngles.Yaw, angles.Yaw, cfg.Smoothness)
	if !input.IsKeyPressed(input.VK_LBUTTON) {
		return
	}
	err := memory.Write(aimbotInput.Handle, aimbotInput.HWModuleBaseAddress+offsets.Pitch, smoothPitch)
	if err != nil {
		fmt.Println(err)
	}

	memory.Write(aimbotInput.Handle, aimbotInput.HWModuleBaseAddress+offsets.Yaw, smoothYaw)

	//fmt.Printf("Aiming at %s %f %f", closestPlayer.Name, smoothPitch, smoothYaw)
	if err != nil {
		fmt.Println(err)
	}
}
