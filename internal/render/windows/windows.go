package windows

import (
	"runtime"
	"syscall"
	"unsafe"
	"log"
	"image"
	_ "image/png"
	"image/draw"
	"os"

	"github.com/lxn/win"
)

const (
	AC_SRC_OVER = 0x00
	ULW_ALPHA   = 0x00000002
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procUpdateLayeredWindow = user32.NewProc("UpdateLayeredWindow")
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
	state.loadSprites("assets/koneko.png")

	instance := win.GetModuleHandle(nil)
	className, _ := syscall.UTF16PtrFromString("konekoClass")
	wc := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		LpfnWndProc:   syscall.NewCallback(WndProc),
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

func (s *WindowsState) UpdatePosition(x, y, spriteX, spriteY int) {
	spriteIndex := spriteY*8 + spriteX
	copy(s.PixelData[:], s.Frames[spriteIndex])

	ptDst := win.POINT{X: int32(x), Y: int32(y)}
	ptSrc := win.POINT{X: 0, Y: 0}
	size := win.SIZE{CX: s.Width, CY: s.Height}
	blend := win.BLENDFUNCTION{BlendOp: AC_SRC_OVER, SourceConstantAlpha: 255, AlphaFormat: win.AC_SRC_ALPHA}

	procUpdateLayeredWindow.Call(
		uintptr(s.window), uintptr(s.ScreenDC), uintptr(unsafe.Pointer(&ptDst)), uintptr(unsafe.Pointer(&size)),
		uintptr(s.MemDC), uintptr(unsafe.Pointer(&ptSrc)), 0, uintptr(unsafe.Pointer(&blend)), uintptr(ULW_ALPHA),
	)
}

func (win *WindowsState) loadSprites(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open PNG: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Unable to decode PNG: %v", err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	win.Frames = make([][]byte, 32)
	for i := range win.Frames {
		win.Frames[i] = make([]byte, 32*32*4)
	}

	for spriteY := 0; spriteY < 4; spriteY++ {
		for spriteX := 0; spriteX < 8; spriteX++ {
			for y := 0; y < 32; y++ {
				for x := 0; x < 32; x++ {
					spriteIndex := spriteY*8 + spriteX
					pixOffset := rgba.PixOffset(spriteX*32+x, spriteY*32+y)
					r, g, b, a := int(rgba.Pix[pixOffset]), int(rgba.Pix[pixOffset+1]), int(rgba.Pix[pixOffset+2]), int(rgba.Pix[pixOffset+3])
					r, g, b = (r*a)/255, (g*a)/255, (b*a)/255

					frameOffest := (y*32 + x) * 4
					win.Frames[spriteIndex][frameOffest+0] = byte(b)
					win.Frames[spriteIndex][frameOffest+1] = byte(g)
					win.Frames[spriteIndex][frameOffest+2] = byte(r)
					win.Frames[spriteIndex][frameOffest+3] = byte(a)
				}
			}
		}
	}
}

func WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == win.WM_DESTROY {
		win.PostQuitMessage(0)
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}