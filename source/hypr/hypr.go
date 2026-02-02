package hypr

import (
	"errors"
	"fmt"
	"strconv"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func getSocketPath() (string, error) {
	xdgRuntineDir, flag := os.LookupEnv("XDG_RUNTIME_DIR")
	if !flag {
		return "", errors.New("XDG_RUNTIME_DIR cannot be determined!")
	}
	hyprInstanceSign, flag := os.LookupEnv("HYPRLAND_INSTANCE_SIGNATURE")
	if !flag {
		return "", errors.New("HYPRLAND_INSTANCE_SIGNATURE cannot be determined!")
	}
	return filepath.Join(xdgRuntineDir, "hypr", hyprInstanceSign, ".socket.sock"), nil
}

func GetCursorPos() (int, int, error) {
	socketPath, err := getSocketPath()
	if err != nil {
		fmt.Print(fmt.Errorf("could not locate hypr socket: %s", err))
	}

	connection, err := net.Dial("unix", socketPath)
	if err != nil {
		return 0, 0, fmt.Errorf("connection to hypr socket failed! : %s", err)
	}
	defer connection.Close()

	_, err = connection.Write([]byte("cursorpos"))
	if err != nil {
		return 0, 0, fmt.Errorf("writing to socket failed! : %s", err)
	}
	buff := make([]byte, 1024)
	n, err := connection.Read(buff)
	if err != nil {
		return 0, 0, fmt.Errorf("could not read from socket! %s", err)
	}
	response := string(buff[:n])
	coordinates := strings.Split(response, ",")
	if len(coordinates) != 2 {
		return 0, 0, fmt.Errorf("unexpected value received from the socket! expected size : 2, received : %d", len(coordinates))
	}
	var x, y int
	x, err = strconv.Atoi(strings.TrimSpace(coordinates[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse the coordinate : %s, expected int : %s", coordinates[0], err)
	}
	y, err = strconv.Atoi(strings.TrimSpace(coordinates[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("could not parse the coordinate : %s, expected int : %s", coordinates[1], err)
	}
	return x, y, nil
}