package main

import (
	"bytes"
	_ "embed"
	"image"
	"image/png"
	_ "image/png"
	"koneko/source/hypr"
	"log"
	"os"

	"github.com/diamondburned/gotk4-layer-shell/pkg/gtk4layershell"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

//go:embed assets/koneko.png
var imgBytes []byte

func imageToTexture(img image.Image) *gdk.Texture {
	var buff bytes.Buffer
	err := png.Encode(&buff, img) 
	if err != nil {
		log.Fatal(err)
	}
	bytes := glib.NewBytes(buff.Bytes())
	texture, err := gdk.NewTextureFromBytes(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return texture
}

func cropImage(img image.Image, x, y, w, h int) image.Image {
	sub, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})
	if ok {
		return sub.SubImage(image.Rect(x, y, x + w, y + h))
	}
	return img
}

func transparentWindow(window *gtk.Window) {
	css := "window { background-color: rgba(0,0,0,0); }"
	provider := gtk.NewCSSProvider()
	provider.LoadFromData(css)
	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func main() {

	app := gtk.NewApplication("com.koneko", 0)
	app.ConnectActivate(func() {

		window := gtk.NewWindow()
		app.AddWindow(window)
		gtk4layershell.InitForWindow(window)
		gtk4layershell.SetLayer(window, gtk4layershell.LayerShellLayerOverlay)
		gtk4layershell.SetKeyboardMode(window, gtk4layershell.LayerShellKeyboardModeNone)
		gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeTop, true)
		gtk4layershell.SetAnchor(window, gtk4layershell.LayerShellEdgeLeft, true)
		transparentWindow(window)

		fullImg, _, err := image.Decode(bytes.NewReader(imgBytes))
		if err == nil {
			spriteImg := cropImage(fullImg, 0, 0, 32, 32)
			texture := imageToTexture(spriteImg)
			if texture != nil {
				pic := gtk.NewPictureForPaintable(texture)
				pic.SetCanShrink(false)
				pic.SetContentFit(gtk.ContentFitContain)

				window.SetChild(pic)
			}
		} else {
			log.Fatal(err)
		}

		glib.TimeoutAdd(100, func() bool {
			x, y, err := hypr.GetCursorPos()
			if err != nil {
				log.Println(err)
				return true
			}
			gtk4layershell.SetMargin(window, gtk4layershell.LayerShellEdgeTop, min(y, y - 16))
			gtk4layershell.SetMargin(window, gtk4layershell.LayerShellEdgeLeft, min(x, x - 16))
			return true
		})

		window.Present()
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
