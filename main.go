package main

import (
	"log"
	"time"

	"koneko/source/wayland"
	"koneko/source/game"
	"koneko/source/hypr"
)

func main() {
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
			wayland.UpdateInWayland(waylandWindow, min(currX, currX - 16), min(currY, currY - 40), spriteX / 32, spriteY / 32)
			lastSpriteX = spriteX
			lastSpriteY = spriteY
		}

		if lastY != currY || lastX != currX {
			wayland.UpdateInWayland(waylandWindow, min(currX, currX - 16), min(currY, currY - 40), spriteX / 32, spriteY / 32)
			lastY = currY
			lastX = currX
		}
	}
}
