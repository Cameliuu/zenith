package entity

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/Cameliuu/zenith/engine/offsets"
	"github.com/Cameliuu/zenith/memory"
	"golang.org/x/sys/windows"
)

/*
================================================================ Utilities ================================================================
*/
type ViewMatrix struct {
	Matrix [16]float32
}
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

/*
================================================================= Entities ================================================================
*/
type PlayerInfo struct {
	Number    int32
	Cmd       string
	Name      string
	Model     string
	AnimFrame float32
	Smth      float32
	Position  Vector3
}

type PlayerInfoRaw struct {
	Number    int32
	Cmd       [256]byte
	Name      [44]byte
	Model     [80]byte
	AnimFrame float32
	Smth      float32
	Position  Vector3
	Pad       [188]byte
}

/*
==================================================================== Local Player =======================================================
*/

func ToPlayerInfo(raw PlayerInfoRaw) PlayerInfo {
	return PlayerInfo{
		Number:    raw.Number,
		Cmd:       string(bytes.Trim(raw.Cmd[:], "\x00")),
		Name:      string(bytes.Trim(raw.Name[:], "\x00")),
		Model:     string(bytes.Trim(raw.Model[:], "\x00")),
		AnimFrame: raw.AnimFrame,
		Smth:      raw.Smth,
		Position:  raw.Position,
	}
}

type LocalPlayerInfo struct {
	ViewMatrix ViewMatrix
	Team       int32
}

func GetLocalPlayerInfo(handle windows.Handle, moduleBaseAddress uintptr) (LocalPlayerInfo, error) {

	team, err := memory.Read[int32](handle, moduleBaseAddress+offsets.LocalPlayerTeam)

	if err != nil {
		return LocalPlayerInfo{}, fmt.Errorf("zenith: could not read local player :%w", err)
	}

	viewMatrix, err := memory.Read[ViewMatrix](handle, moduleBaseAddress+offsets.ViewMatrix)

	if err != nil {
		return LocalPlayerInfo{}, fmt.Errorf("zenith: could not read local player :%w", err)
	}

	return LocalPlayerInfo{
		Team:       team,
		ViewMatrix: viewMatrix,
	}, nil

}
func GetEntities(handle windows.Handle, moduleBaseAddress uintptr) ([]PlayerInfo, error) {
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

	return entities, nil
}
