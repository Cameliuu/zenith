package memory

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func FindProcessIdByName(targetProc string) (uint32, error) {
	procHandle, procErr := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)

	if procErr != nil {
		return 0, fmt.Errorf("zenith: could not create process snapshot: %w", procErr)
	}
	defer windows.CloseHandle(procHandle)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	procErr = windows.Process32First(procHandle, &procEntry)

	for procErr == nil {
		if targetProc == windows.UTF16ToString(procEntry.ExeFile[:]) {
			return procEntry.ProcessID, nil
		}
		procErr = windows.Process32Next(procHandle, &procEntry)
	}
	return 0, fmt.Errorf("zenith: Could not find process: %w", procErr)
}

// TO-DO: Add handle hijacking
func OpenHandleToProcessByPID(PID uint32) (windows.Handle, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_VM_WRITE|
		windows.PROCESS_VM_READ,
		false,
		PID)
	if err != nil {
		return 0, fmt.Errorf("zenith: could not open handle to process: %w", err)
	}

	defer windows.CloseHandle(handle)
	return handle, nil
}

// TO-DO: Open just one snapshot and search for more modules
func FindModuleBaseAddressByPid(PID uint32, moduleName string) (uintptr, error) {
	modHandle, modErr := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPMODULE32|windows.TH32CS_SNAPMODULE, PID)

	if modErr != nil {
		return 0, fmt.Errorf("zenith: could not module snapshot: %w", modErr)
	}

	defer windows.CloseHandle(modHandle)

	var modEntry windows.ModuleEntry32
	modEntry.Size = uint32(unsafe.Sizeof(modEntry))
	modErr = windows.Module32First(modHandle, &modEntry)

	for modErr == nil {
		if windows.UTF16ToString(modEntry.Module[:]) == moduleName {
			return modEntry.ModBaseAddr, nil
		}
		modErr = windows.Module32Next(modHandle, &modEntry)
	}
	return 0, fmt.Errorf("zenith: Could not find process: %w", modErr)
}
