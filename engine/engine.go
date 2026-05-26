package engine

import (
	"fmt"

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

func Run() error {
	_, err := Init()

	if err != nil {
		return err
	}
	fmt.Println("zenith: engine initialized")
	return nil
}
