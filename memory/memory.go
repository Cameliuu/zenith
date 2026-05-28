package memory

import (
	"bytes"
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
func Read[T any](handle windows.Handle, address uintptr) (T, error) {
	var result T
	err := windows.ReadProcessMemory(
		handle,
		address,
		(*byte)(unsafe.Pointer(&result)),
		unsafe.Sizeof(result),
		nil,
	)
	return result, err
}

func ReadString(handle windows.Handle, address uintptr, length uint) (string, error) {
	buf := make([]byte, length)

	err := windows.ReadProcessMemory(handle,
		address,
		&buf[0],
		uintptr(length),
		nil)

	if err != nil {
		return "", nil
	}

	return string(bytes.TrimRight(buf, "\x00")), nil
}
func Write[T any](handle windows.Handle, address uintptr, value T) error {
	var written uintptr
	return windows.WriteProcessMemory(
		handle,
		address,
		(*byte)(unsafe.Pointer(&value)),
		unsafe.Sizeof(value),
		&written,
	)
}

// TO-DO: Add handle hijacking
func OpenHandleToProcessByPID(PID uint32) (windows.Handle, error) {
	handle, err := windows.OpenProcess(windows.PROCESS_VM_WRITE|
		windows.PROCESS_VM_READ|
		windows.PROCESS_VM_OPERATION,
		false,
		PID)
	if err != nil {
		return 0, fmt.Errorf("zenith: could not open handle to process: %w", err)
	}

	return handle, nil
}

func ReadMemoryInBulk(handle windows.Handle, moduleBaseAddress uintptr, buf []byte) error {
	err := windows.ReadProcessMemory(handle, moduleBaseAddress, &buf[0], uintptr(len(buf)), nil)
	if err != nil {
		return fmt.Errorf("zenith: error reading memory in bulk: %w", err)
	}
	return nil
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
