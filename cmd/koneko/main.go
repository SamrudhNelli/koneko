package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"koneko/internal/game"
	"koneko/internal/input/hypr"
	"koneko/internal/render/wayland"
	"koneko/internal/render/windows"
	winInput "koneko/internal/input/windows"
	"github.com/lxn/win"
)

func main() {
	if runtime.GOOS == "windows" {
		runOnWindows()
	} else if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		runOnHyprland()
	} else {
		log.Fatal("Sorry, your operating system is not supported yet!")
	}
}

func runOnHyprland() {
	waylandWindow := wayland.SetupWayland()
	defer waylandWindow.Close()

	var lastX int = 100
	var lastY int = 100
	var lastSpriteX int = 64
	var lastSpriteY int = 0
	ticker := time.NewTicker(time.Second / 6)
	defer ticker.Stop()

	for range ticker.C {
		cursorX, cursorY, err := hypr.GetCursorPos()
		if err != nil {
			log.Println(err)
		}
		spriteX, spriteY, currX, currY := game.GetSpriteCoord(lastX, lastY, cursorX, cursorY)

		if spriteX != lastSpriteX || spriteY != lastSpriteY {
			wayland.UpdateInWayland(waylandWindow, min(currX, currX-16), min(currY, currY-40), spriteX/32, spriteY/32)
			lastSpriteX = spriteX
			lastSpriteY = spriteY
		}

		if lastY != currY || lastX != currX {
			wayland.UpdateInWayland(waylandWindow, min(currX, currX-16), min(currY, currY-40), spriteX/32, spriteY/32)
			lastY = currY
			lastX = currX
		}
	}
}

func runOnWindows() {
	state := windows.SetupWindows()

	go func() {
		var lastX int = 100
		var lastY int = 100
		var lastSpriteX int = 64
		var lastSpriteY int = 0
		ticker := time.NewTicker(time.Second / 6)
		defer ticker.Stop()

		for range ticker.C {
			cursorX, cursorY, err := winInput.GetCursorPos()
			if err != nil {
				log.Fatal(err)
			}

			spriteX, spriteY, currX, currY := game.GetSpriteCoord(lastX, lastY, cursorX, cursorY)

			if spriteX != lastSpriteX || spriteY != lastSpriteY {
				state.UpdatePosition(min(currX, currX-16), min(currY, currY-10), spriteX/32, spriteY/32)
				lastSpriteX = spriteX
				lastSpriteY = spriteY
			}

			if lastY != currY || lastX != currX {
				state.UpdatePosition(min(currX, currX-16), min(currY, currY-10), spriteX/32, spriteY/32)
				lastY = currY
				lastX = currX
			}
		}
	}()

	var msg win.MSG
	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}
