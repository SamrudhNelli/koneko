package wayland

import (
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"

	"koneko/protocols/wlr"

	sys "github.com/neurlang/wayland/os"
	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
)

type winState struct {
	appID         string
	title         string
	pImage        *image.RGBA
	data          []byte
	width, height int32
	frame         *image.RGBA
	exit          bool
	display       *wl.Display
	registry      *wl.Registry
	shm           *wl.Shm
	zwlr          *wlr.ZwlrLayerShellV1
	layerSurface  *wlr.ZwlrLayerSurfaceV1
	compositor    *wl.Compositor
	surface       *wl.Surface
	buff          *wl.Buffer
	frames        [][]byte
}

func loadImage(path string) (*image.RGBA, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)
	return rgba, nil
}

func (window *winState) loadImageIntoFrames() {
	window.frames = make([][]byte, 32)
	for i := range window.frames {
		window.frames[i] = make([]byte, 32*32*4)
	}
	for spriteX := range 8 {
		for spriteY := range 4 {
			for y := range 32 {
				for x := range 32 {
					pixel := (y * 32 * 4) + (x * 4)
					imgCoord := window.pImage.PixOffset(x+spriteX*32, y+spriteY*32)
					window.frames[spriteY*8+spriteX][pixel+0] = window.pImage.Pix[imgCoord+2]
					window.frames[spriteY*8+spriteX][pixel+1] = window.pImage.Pix[imgCoord+1]
					window.frames[spriteY*8+spriteX][pixel+2] = window.pImage.Pix[imgCoord+0]
					window.frames[spriteY*8+spriteX][pixel+3] = window.pImage.Pix[imgCoord+3]
				}
			}
		}
	}
}

func SetupWayland() *winState {
	window := &winState{
		title: "koneko",
		appID: "koneko",
	}

	var err error
	window.pImage, err = loadImage("assets/koneko.png")
	if err != nil {
		log.Fatalf("Unable to load image : %v", err)
	}

	window.loadImageIntoFrames()

	display, err := wl.Connect("")
	if err != nil {
		log.Fatalf("Unable to connect to wayland server! %v", err)
	}
	window.display = display
	display.AddErrorHandler(window)

	registry, err := window.display.GetRegistry()
	if err != nil {
		log.Fatalf("Unable to connect to Global registery! %v", err)
	}
	window.registry = registry
	registry.AddGlobalHandler(window)

	_ = wlclient.DisplayRoundtrip(window.display)

	window.surface, err = window.compositor.CreateSurface()
	if err != nil {
		log.Printf("Unable to create a zwlr surface! %v", err)
	}

	window.layerSurface, err = window.zwlr.GetLayerSurface(window.surface, nil, wlr.ZwlrLayerShellV1LayerOverlay, "koneko")
	if err != nil {
		log.Printf("Unable to create a zwlr layer: %v", err)
	}

	return window
}

func (window *winState) Close() {
	if window.shm != nil {
		window.releaseShm()
	}
	if window.registry != nil {
		window.releaseRegistry()
	}
	window.display.Context().Close()
}

func UpdateInWayland(window *winState, x, y, spriteX, spriteY int) {
	window.height = 32
	window.width = 32
	window.layerSurface.SetSize(32, 32)
	window.layerSurface.SetAnchor(wlr.ZwlrLayerSurfaceV1AnchorTop | wlr.ZwlrLayerSurfaceV1AnchorLeft)
	window.layerSurface.SetKeyboardInteractivity(0)

	window.layerSurface.AddConfigureHandler(window)
	window.surface.Commit()

	_ = wlclient.DisplayRoundtrip(window.display)

	window.InitBuffer()

	copy(window.data, window.frames[spriteY*8+spriteX])

	window.layerSurface.SetMargin(int32(y), 0, 0, int32(x))
	window.surface.Attach(window.buff, 0, 0)
	window.surface.Damage(0, 0, 32, 32)
	window.surface.Commit()

	if err := wlclient.DisplayDispatch(window.display); err != nil {
		return
	}
}

func (win *winState) HandleZwlrLayerSurfaceV1Configure(e wlr.ZwlrLayerSurfaceV1ConfigureEvent) {
	win.layerSurface.AckConfigure(e.Serial)
}

func (win *winState) HandleRegistryGlobal(r wl.RegistryGlobalEvent) {
	switch r.Interface {
	case "wl_shm":
		shm := wl.NewShm(win.display.Context())
		err := win.registry.Bind(r.Name, r.Interface, r.Version, shm)
		if err != nil {
			log.Fatalf("Unable to bind wl_shm interface: %v", err)
		}
		win.shm = shm
	case "zwlr_layer_shell_v1":
		zwlr := wlr.NewZwlrLayerShellV1(win.display.Context())
		err := win.registry.Bind(r.Name, r.Interface, r.Version, zwlr)
		if err != nil {
			log.Fatalf("Unable to bind zwlr_layer_shell_v1: %v", err)
		}
		win.zwlr = zwlr
	case "wl_compositor":
		compositor := wl.NewCompositor(win.display.Context())
		err := win.registry.Bind(r.Name, r.Interface, r.Version, compositor)
		if err != nil {
			log.Fatalf("Unable to bind wl_compositor: %v", err)
		}
		win.compositor = compositor
	}
}

func (win *winState) releaseShm() {
	win.shm.Unregister()
	win.shm = nil
}

func (win *winState) releaseRegistry() {
	win.registry.RemoveGlobalHandler(win)
	win.registry.Unregister()
	win.registry = nil
}

func (win *winState) InitBuffer() {
	stride := win.width * 4
	size := stride * win.height

	file, err := sys.CreateAnonymousFile(int64(size))
	if err != nil {
		log.Fatalf("Unable to create a temporary file!: %v", err)
	}
	defer file.Close()

	data, err := sys.Mmap(int(file.Fd()), 0, int(size), sys.ProtRead|sys.ProtWrite, sys.MapShared)
	if err != nil {
		log.Fatalf("Unable to create a memory mapping: %v", err)
	}
	win.data = data

	pool, err := win.shm.CreatePool(file.Fd(), size)
	if err != nil {
		_ = sys.Munmap(data)
		_ = file.Close()
		log.Fatalf("Unable to create a shared memory pool: %v", err)
	}
	defer pool.Destroy()

	buff, err := pool.CreateBuffer(0, win.width, win.height, stride, wl.ShmFormatArgb8888)
	if err != nil {
		_ = pool.Destroy()
		_ = sys.Munmap(data)
		_ = file.Close()
		log.Fatalf("Unable to create wl.Buffer from shm pool: %v", err)
	}
	win.buff = buff
	win.surface.Attach(buff, 0, 0)
	win.surface.Damage(0, 0, win.width, win.height)
	win.surface.Commit()
}

func (win *winState) HandleDisplayError(e wl.DisplayErrorEvent) {
	log.Fatalf("Display error event: %v", e)
}
