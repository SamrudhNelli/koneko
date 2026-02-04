package main

import (
	"log"

	"koneko/protocols/wlr"

	"github.com/neurlang/wayland/wl"
	"github.com/neurlang/wayland/wlclient"
)

type winState struct {
	appID string
	title string
	// pImage *image.RGBA
	width, height int32
	// frame *image.RGBA
	exit bool
	display *wl.Display
	registry *wl.Registry
	shm *wl.Shm
	zwlr *wlr.ZwlrLayerShellV1
	compositor *wl.Compositor
	// wmBase *xdg.WmBase
	seat *wl.Seat
	surface *wl.Surface
	// xdgSurface  *xdg.Surface
	// xdgTopLevel *xdg.Toplevel
	keyboard *wl.Keyboard
	pointer *wl.Pointer
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

}

func run(window *winState) {
	registry, err := window.display.GetRegistry()
	if err != nil {
		log.Fatalf("Unable to connect to Global registery! %v", err)
	}
	window.registry = registry
	registry.AddGlobalHandler(window)
	_ = wlclient.DisplayRoundtrip(window.display)
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