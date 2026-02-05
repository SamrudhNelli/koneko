package wayland

import (
	"log"
	"image"

	"koneko/protocols/wlr"

	sys "github.com/neurlang/wayland/os"
	"github.com/neurlang/wayland/wl"
	// "github.com/neurlang/wayland/external/swizzle"
	"github.com/neurlang/wayland/wlclient"
)

type winState struct {
	appID string
	title string
	pImage *image.RGBA
	data []byte
	width, height int32
	frame *image.RGBA
	exit bool
	display *wl.Display
	registry *wl.Registry
	shm *wl.Shm
	zwlr *wlr.ZwlrLayerShellV1
	layerSurface *wlr.ZwlrLayerSurfaceV1
	compositor *wl.Compositor
	surface *wl.Surface
	buff *wl.Buffer
}

func main() {
	window := &winState{
		title: "koneko",
		appID: "koneko",
	}

	display, err := wl.Connect("")
	if err != nil {
		log.Fatalf("Unable to connect to wayland server! %v", err)
	}
	window.display = display
	display.AddErrorHandler(window)
	
	run(window)
	
	if window.shm != nil {
		window.releaseShm()
	}
	if window.registry != nil {
		window.releaseRegistry()
	}

	window.display.Context().Close()
}

func run(window *winState) {
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

	window.height = 32
	window.width = 32
    window.layerSurface.SetSize(32, 32)
    window.layerSurface.SetAnchor(wlr.ZwlrLayerSurfaceV1AnchorTop|wlr.ZwlrLayerSurfaceV1AnchorLeft)
    window.layerSurface.SetKeyboardInteractivity(0)
    window.surface.Commit()

    window.InitBuffer()

    for {
        for i := 0; i < len(window.data); i += 4 {
             window.data[i] = 0x00
             window.data[i+1] = 0x00
             window.data[i+2] = 0xFF
             window.data[i+3] = 0xFF
        }

        window.surface.Damage(0, 0, 32, 32)
        window.surface.Commit()

        if err := wlclient.DisplayDispatch(window.display); err != nil {
			break
		}
    }
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

	buff, err := pool.CreateBuffer(0, win.width, win.height, stride, wl.ShmFormatAbgr8888)
	if err != nil {
		_ = pool.Destroy()
		_ = sys.Munmap(data)
		_ = file.Close()
		log.Fatalf("Unable to create wl.Buffer from shm pool: %v", err)
	}
	win.buff = buff

	// copy(data, win.frame.Pix)
	// if err := swizzle.BGRA(data); err != nil {
	// 	log.Printf("unable to convert RGBA to BGRA: %v", err)
	// }

	win.surface.Attach(buff, 0, 0)
	win.surface.Damage(0, 0, win.width, win.height)
    win.surface.Commit()
}