package input

import "golang.org/x/sys/windows"

var (
	user32DLL            = windows.NewLazyDLL("user32.DLL")
	procGetAsyncKeyState = user32DLL.NewProc("GetAsyncKeyState")
)

const (
	VK_LBUTTON = 0x01 // left mouse
	VK_RBUTTON = 0x02 // right mouse
	VK_MBUTTON = 0x04 // middle mouse
	VK_INSERT  = 0x2D // insert
	VK_DELETE  = 0x2E // delete
	VK_SPACE   = 0x20 // space
	VK_SHIFT   = 0x10 // shift
	VK_CONTROL = 0x11 // ctrl
	VK_ALT     = 0x12 // alt
	VK_F1      = 0x70 // F1
	VK_F2      = 0x71 // F2
)

var KeyBinds = map[string]int{
	"F1":  0x70,
	"F2":  0x71,
	"F3":  0x72,
	"F4":  0x73,
	"F5":  0x74,
	"F6":  0x75,
	"F7":  0x76,
	"F8":  0x77,
	"F9":  0x78,
	"F10": 0x79,
}

var KeyNames = []string{"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10"}

var KeyNameFromCode = map[int]string{
	0x70: "F1",
	0x71: "F2",
	0x72: "F3",
	0x73: "F4",
	0x74: "F5",
	0x75: "F6",
	0x76: "F7",
	0x77: "F8",
	0x78: "F9",
	0x79: "F10",
}

func IsKeyJustPressed(vk int) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return r&0x0001 != 0
}

func IsKeyPressed(vk int) bool {
	r, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return r&0x8000 != 0
}
