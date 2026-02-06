package windows

import (
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	user32GetSystemMetrics = user32.NewProc("GetSystemMetrics")
	user32DefWindowProc = user32.NewProc("DefWindowProcW")
	user32GetMessage = user32.NewProc("GetMessageW")
)

func runOnWindows() {
	instance, err := windows.GetModuleHandle(nil)
	if err != nil {
		log.Fatalf("Unable to acces module handle: %v", err)
	}

	className, _ := windows.UTF16PtrFromString("konekoClass")
	wc := windows.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(windows.WNDCLASSEX{})),
		LpfnWndProc:   windows.NewCallback(DefWindowProc),
		HInstance:     instance,
		LpszClassName: className,
		HbrBackground: windows.HBRUSH(windows.GetStockObject(windows.WHITE_BRUSH)),
	}

	if errCode := windows.RegisterClassEx(&wc); errCode == 0 {
		log.Fatal("Unable to register class!")
	}

	winStyle := windows.WS_EX_LAYERED | windows.WS_EX_TRANSPARENT | windows.WS_EX_TOPMOST | windows.WS_EX_TOOLWINDOW
	style := windows.WS_POPUP

	screenWidth := getSystemMetrics(0)  // SM_CXSCREEN
	screenHeight := getSystemMetrics(1) // SM_CYSCREEN

	window, err := windows.CreateWindowEx(uint32(winStyle), className, nil, uint32(style), 0, 0, int32(screenWidth), int32(screenHeight), 0, 0, instance, nil)

	if err != nil {
		log.Fatalf("Unable to create a window: %v", err)
	}

	windows.ShowWindow(window, windows.SW_SHOW)

	var msg windows.MSG
	for {
		if r, _, _ := user32GetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0); r == 0 {
			break
		}
		windows.TranslateMessage(&msg)
		windows.DispatchMessage(&msg)
	}
}

func DefWindowProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	r, _, _ := user32DefWindowProc.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return r
}

func getSystemMetrics(index int) int {
	ret, _, _ := user32GetSystemMetrics.Call(uintptr(index))
	return int(ret)
}