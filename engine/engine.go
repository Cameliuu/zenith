package engine

import (
	"fmt"
	"log"
	"time"

	entity "github.com/Cameliuu/zenith/engine/entitiy"
	"github.com/Cameliuu/zenith/memory"
	"golang.org/x/sys/windows"
)

type GameData struct {
	handle                  windows.Handle
	hwModuleBaseAddress     uintptr
	clientModuleBaseAddress uintptr
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
		handle:                  handle,
		hwModuleBaseAddress:     hwModuleBase,
		clientModuleBaseAddress: clientModuleBase,
	}, nil
}

func Run() {
	gameData, err := Init()

	if err != nil {
		log.Fatal("zenith: could not initialize engine: %w", err)
	}
	defer windows.CloseHandle(gameData.handle)
	fmt.Println("zenith: engine initialized")

	for {
		localPlayerInfo, err := entity.GetLocalPlayerInfo(gameData.handle, gameData.clientModuleBaseAddress)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(localPlayerInfo)

		_, err = entity.GetEntities(gameData.handle, gameData.hwModuleBaseAddress)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
