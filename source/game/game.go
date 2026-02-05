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
	frameCnt++
	diffX := max(lastX - cursorX, cursorX - lastX)
	diffY := max(lastY - cursorY, cursorY - lastY)
	spriteX = spriteSets["idle"][frameCnt % len(spriteSets["idle"])][0]
	spriteY = spriteSets["idle"][frameCnt % len(spriteSets["idle"])][1]
	currX = lastX
	currY = lastY
	if diffY <= 25 && diffX <= 25 {
		return
	} else if  diffX > konekoSpeed && diffY > konekoSpeed {
		if lastX > cursorX && lastY > cursorY {
			// NW
			spriteX = spriteSets["NW"][frameCnt % len(spriteSets["NW"])][0]
			spriteY = spriteSets["NW"][frameCnt % len(spriteSets["NW"])][1]
			currX = lastX - konekoSpeed
			currY = lastY - konekoSpeed
		} else if lastX < cursorX && lastY < cursorY {
			// SE
			spriteX = spriteSets["SE"][frameCnt % len(spriteSets["SE"])][0]
			spriteY = spriteSets["SE"][frameCnt % len(spriteSets["SE"])][1]
			currX = lastX + konekoSpeed
			currY = lastY + konekoSpeed
		} else if lastX > cursorX && lastY < cursorY {
			// SW
			spriteX = spriteSets["SW"][frameCnt % len(spriteSets["SW"])][0]
			spriteY = spriteSets["SW"][frameCnt % len(spriteSets["SW"])][1]
			currX = lastX - konekoSpeed
			currY = lastY + konekoSpeed
		} else {
			// NE
			spriteX = spriteSets["NE"][frameCnt % len(spriteSets["NE"])][0]
			spriteY = spriteSets["NE"][frameCnt % len(spriteSets["NE"])][1]
			currX = lastX + konekoSpeed
			currY = lastY - konekoSpeed
		}
	} else if diffX > konekoSpeed {
		if lastX > cursorX {
			// W
			spriteX = spriteSets["W"][frameCnt % len(spriteSets["W"])][0]
			spriteY = spriteSets["W"][frameCnt % len(spriteSets["W"])][1]
			currX = lastX - konekoSpeed
			currY = lastY
		} else {
			// E
			spriteX = spriteSets["E"][frameCnt % len(spriteSets["E"])][0]
			spriteY = spriteSets["E"][frameCnt % len(spriteSets["E"])][1]
			currX = lastX + konekoSpeed
			currY = lastY
		}
	} else if diffY > konekoSpeed {
		if lastY > cursorY {
			// N
			spriteX = spriteSets["N"][frameCnt % len(spriteSets["N"])][0]
			spriteY = spriteSets["N"][frameCnt % len(spriteSets["N"])][1]
			currX = lastX
			currY = lastY - konekoSpeed
		} else {
			// S
			spriteX = spriteSets["S"][frameCnt % len(spriteSets["S"])][0]
			spriteY = spriteSets["S"][frameCnt % len(spriteSets["S"])][1]
			currX = lastX
			currY = lastY + konekoSpeed
		}
	}
	return
}