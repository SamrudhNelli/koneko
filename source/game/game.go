package game

var spriteSets = map[string][] []int{
	"idle":         {{96, 96}},
	"alert":        {{224, 96}},
	"scratchSelf":  {{160, 0}, {192, 0}, {224, 0}},
	"scratchWallN": {{0, 0}, {0, 32}},
	"scratchWallS": {{224, 32}, {192, 64}},
	"scratchWallE": {{64, 64}, {64, 96}},
	"scratchWallW": {{128, 0}, {128, 32}},
	"tired":        {{96, 64}},
	"sleeping":     {{64, 0}, {64, 32}},
	"N":            {{32, 64}, {32, 96}},
	"NE":           {{0, 64}, {0, 96}},
	"E":            {{96, 0}, {96, 32}},
	"SE":           {{160, 32}, {160, 64}},
	"S":            {{192, 96}, {224, 64}},
	"SW":           {{160, 96}, {192, 32}},
	"W":            {{128, 64}, {128, 96}},
	"NW":           {{32, 0}, {32, 32}},
}

var konekoSpeed int = 16
var frameCnt int = 0

func GetSpriteCoord(lastX, lastY, cursorX, cursorY int) (spriteX, spriteY, currX, currY int) {
	diffX := max(lastX - cursorX, cursorX - lastX)
	diffY := max(lastY - cursorY, cursorY - lastY)
	if diffY <= 16 || diffX <= 16 {
		spriteX = spriteSets["idle"][frameCnt % len(spriteSets["idle"])][0]
		spriteY = spriteSets["idle"][frameCnt % len(spriteSets["idle"])][1]
		currX = lastX
		currY = lastY
		return
	} else if  diffX > konekoSpeed && diffY > konekoSpeed {
		// NE, SE, SW, NW
		if lastX > cursorX && lastY > cursorY {

		} else if lastX < cursorX && lastY < cursorY {

		} else if lastX > cursorX && lastY < cursorY {

		} else {

		}
	} else if diffX > konekoSpeed {
		// E, W
		if lastX > cursorX {

		} else {

		}
	} else if diffY > konekoSpeed {
		// N, S
		if lastY > cursorY {

		}
	}
	return
}