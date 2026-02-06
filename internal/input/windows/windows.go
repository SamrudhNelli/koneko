package windows

import(
	"github.com/lxn/win"
	"errors"
)

func GetCursorPos() (int, int, error) {
	var pt win.POINT
	flag := win.GetCursorPos(&pt)
	if !flag {
		return 0, 0, errors.New("unable to get cursor position")
	}
	return int(pt.X), int(pt.Y), nil
}
