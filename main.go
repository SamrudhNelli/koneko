package main

import (
	"os"
	"koneko/source/hypr"
	"log"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
)

func makeTransparent(win *gtk.Window) {
	css := "window { background-color: rgba(0,0,0,0); }"
	provider := gtk.NewCSSProvider()
	provider.LoadFromData(css)
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func main() {
	x, y, err := hypr.GetCursorPos()
	if err != nil {
		log.Fatal(err)
	}

	app := gtk.NewApplication("com.koneko", 0)
	app.ConnectActivate(func() {
		win := gtk.NewWindow()
		app.AddWindow(win)
		gtk4layershell.InitForWindow(win)
		gtk4layershell.SetLayer(win, gtk4layershell.LayerShellLayerOverlay)
		gtk4layershell.SetKeyboardMode(win, gtk4layershell.LayerShellKeyboardModeNone)
		gtk4layershell.SetAnchor(win, gtk4layershell.LayerShellEdgeTop, true)
		gtk4layershell.SetAnchor(win, gtk4layershell.LayerShellEdgeLeft, true)
		gtk4layershell.SetMargin(win, gtk4layershell.LayerShellEdgeTop, y)
		gtk4layershell.SetMargin(win, gtk4layershell.LayerShellEdgeLeft, x)
		makeTransparent(win)

		texture, err := gdk.NewTextureFromFilename("assets/koneko.png")
		if err == nil {
			pic := gtk.NewPictureForPaintable(texture)
			pic.SetContentFit(gtk.ContentFitContain)
			pic.SetCanShrink(false)
			win.SetChild(pic)
		} else {
			log.Fatal("error: could not locate the image!")
		}
		win.Present()
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}