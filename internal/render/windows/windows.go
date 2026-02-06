package windows

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

const (
	AC_SRC_OVER = 0x00
	ULW_ALPHA   = 0x00000002
)

var (
	gdi32                   = syscall.NewLazyDLL("gdi32.dll")
	procCreateDIBSection    = gdi32.NewProc("CreateDIBSection")
)

type WindowsState struct {
	window      win.HWND
	MemDC     win.HDC
	PixelData *[32 * 32 * 4]byte
	Frames    [][]byte
	Width     int32
	Height    int32
	ScreenDC  win.HDC
}

func init() {
	runtime.LockOSThread()
}

func SetupWindows() *WindowsState {
	state := &WindowsState{Width: 32, Height: 32}

	instance := win.GetModuleHandle(nil)
	className, _ := syscall.UTF16PtrFromString("konekoClass")
	wc := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(wndProc),
		HInstance:     instance,
		LpszClassName: className,
		HCursor:       win.LoadCursor(0, (*uint16)(unsafe.Pointer(uintptr(win.IDC_ARROW)))),
	}
	win.RegisterClassEx(&wc)

	state.window = win.CreateWindowEx(
		win.WS_EX_LAYERED|win.WS_EX_TOPMOST|win.WS_EX_TOOLWINDOW|win.WS_EX_TRANSPARENT,
		className, nil,
		win.WS_POPUP|win.WS_VISIBLE,
		0, 0, 32, 32,
		0, 0, instance, nil,
	)

	state.ScreenDC = win.GetDC(0)
	state.MemDC = win.CreateCompatibleDC(state.ScreenDC)

	bi := win.BITMAPINFOHEADER{
		BiSize:        uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
		BiWidth:       state.Width,
		BiHeight:      -state.Height,
		BiPlanes:      1,
		BiBitCount:    32,
		BiCompression: win.BI_RGB,
	}
	var bits unsafe.Pointer
	hBitmap, _, _ := procCreateDIBSection.Call(uintptr(state.MemDC), uintptr(unsafe.Pointer(&win.BITMAPINFO{BmiHeader: bi})), 0, uintptr(unsafe.Pointer(&bits)), 0, 0,)
	win.SelectObject(state.MemDC, win.HGDIOBJ(hBitmap))
	state.PixelData = (*[32 * 32 * 4]byte)(bits)

	return state
}

func wndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == win.WM_DESTROY {
		win.PostQuitMessage(0)
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}