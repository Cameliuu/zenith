package entity

import (
	"bytes"
	"fmt"
	"strings"
	"unsafe"

	"github.com/Cameliuu/zenith/engine/entitiy/utils"
	"github.com/Cameliuu/zenith/engine/offsets"
	"github.com/Cameliuu/zenith/engine/overlay"
	"github.com/Cameliuu/zenith/memory"
	"golang.org/x/sys/windows"
)

/*
================================================================= Entities ================================================================
*/
type PlayerInfo struct {
	Number        int32
	Cmd           string
	Name          string
	Model         string
	IsAlive       int32
	AnimFrame     float32
	Smth          float32
	Position      utils.Vector3
	PrevAnimFrame float32
	IsStale       bool
}

type PlayerInfoRaw struct {
	Number    int32
	Cmd       [256]byte
	Name      [44]byte
	Model     [76]byte
	isAlive   int32
	AnimFrame float32
	Smth      float32
	Position  utils.Vector3
	Pad       [188]byte
}

/*
==================================================================== Local Player =======================================================
*/
func decodeASCII(b []byte) string {
	if i := bytes.IndexByte(b, 0); i != -1 {
		b = b[:i]
	}

	var out strings.Builder
	out.Grow(len(b))

	for _, c := range b {
		if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' {
			out.WriteByte(c)
		}
	}

	return out.String()
}

func ToPlayerInfo(raw PlayerInfoRaw) PlayerInfo {
	return PlayerInfo{
		Number:    raw.Number,
		Cmd:       decodeASCII(raw.Cmd[:]),
		Name:      decodeASCII(raw.Name[:]),
		Model:     decodeASCII(raw.Model[:]),
		IsAlive:   raw.isAlive,
		AnimFrame: raw.AnimFrame,
		Smth:      raw.Smth,
		Position:  raw.Position,
	}
}

type LocalPlayerInfo struct {
	ViewMatrix overlay.ViewMatrix
	Team       int32
	ViewAngles overlay.ViewAngles
	Origin     utils.Vector3
}

func GetLocalPlayerInfo(handle windows.Handle, moduleBaseAddress uintptr, HWModuleBaseAddress uintptr) (LocalPlayerInfo, error) {

	team, err := memory.Read[int32](handle, moduleBaseAddress+offsets.LocalPlayerTeam)

	if err != nil {
		return LocalPlayerInfo{}, fmt.Errorf("zenith: could not read local player :%w", err)
	}

	viewMatrix, err := memory.Read[overlay.ViewMatrix](handle, HWModuleBaseAddress+offsets.ViewMatrix)

	if err != nil {
		return LocalPlayerInfo{}, fmt.Errorf("zenith: could not read local player :%w", err)
	}

	viewAngles, err := memory.Read[overlay.ViewAngles](handle, HWModuleBaseAddress+offsets.Pitch)
	if err != nil {
		return LocalPlayerInfo{}, fmt.Errorf("zenith: could not read local player :%w", err)
	}

	origin, err := memory.Read[utils.Vector3](handle, HWModuleBaseAddress+offsets.LocalPlayerOriginOnline)
	if err != nil {
		return LocalPlayerInfo{}, fmt.Errorf("zenith: could not read local player :%w", err)
	}

	return LocalPlayerInfo{
		Team:       team,
		ViewMatrix: viewMatrix,
		ViewAngles: viewAngles,
		Origin:     origin,
	}, nil

}

func GetEntities(handle windows.Handle, moduleBaseAddress uintptr, animFrames, oldAnimFrames *[64]float32) ([]PlayerInfo, error) {
	buf := make([]byte, 37900)
	entitySize := unsafe.Sizeof(PlayerInfoRaw{})
	count := len(buf) / int(entitySize)

	err := memory.ReadMemoryInBulk(handle, moduleBaseAddress+offsets.EntityList, buf)
	if err != nil {
		return nil, err
	}
	var entities []PlayerInfo
	for i := 0; i < count; i++ {
		offset := uintptr(i) * entitySize
		raw := *(*PlayerInfoRaw)(unsafe.Pointer(&buf[offset]))
		player := ToPlayerInfo(raw)
		if player.Name == "" {
			continue
		}
		entities = append(entities, player)
	}
	for i := range entities {
		curr := entities[i].AnimFrame
		old := animFrames[i]
		oldOld := oldAnimFrames[i]

		if curr == old && curr == oldOld && curr != 0 {
			entities[i].IsStale = true
		}

		oldAnimFrames[i] = old
		animFrames[i] = curr
	}

	return entities, nil
}
