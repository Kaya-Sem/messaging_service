package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unsafe"
)

const (
	CLEAR = "\033[2J"
	RESET = "\033c"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getDimensions() (uint, uint) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col), uint(ws.Row)
}

type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

func enableRawMode() (*termios, error) {
	fd := int(syscall.Stdin)
	var oldState termios

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(syscall.TCGETS),
		uintptr(unsafe.Pointer(&oldState)))

	if errno != 0 {
		return nil, errno
	}

	newState := oldState

	// Disable canonical mode (ICANON) and echo (ECHO)
	newState.Lflag &^= syscall.ICANON | syscall.ECHO

	// Set the new terminal attributes
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&newState)))
	if errno != 0 {
		return nil, errno
	}

	return &oldState, nil

}

// disableRawMode restores the terminal to its original state
func disableRawMode(oldState *termios) error {
	fd := int(syscall.Stdin)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(fd), uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(oldState)))
	if errno != 0 {
		return errno
	}
	return nil
}

func readInput() ([]byte, int) {
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		fmt.Printf("Error reading key: %v\n", err)
		return buf, -1
	}
	return buf, n
}

func getFilledChatWindow(_ []string, rows, columns int) [][]string {
	grid := make([][]string, rows)
	for i := range grid {
		grid[i] = make([]string, columns)
	}

	for i := range columns {
		grid[rows-1][i] = "━"
	}

	return grid
}

func gridToString(grid [][]string) string {
	var sb strings.Builder

	rows := len(grid)
	cols := len(grid[0])

	for r := range rows {
		for c := range cols {
			cell := grid[r][c]
			if cell == "" {
				sb.WriteString(" ")
			} else {
				sb.WriteString(cell)
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func main() {
	fmt.Printf("%s", RESET)

	// Enable raw mode once at the start
	oldState, err := enableRawMode()
	if err != nil {
		fmt.Printf("Failed to enable raw mode: %v\n", err)
		return
	}
	defer disableRawMode(oldState) // Ensure terminal is restored on exit

	// Set up signal handler to restore terminal on interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		disableRawMode(oldState)
		fmt.Print(RESET)
		os.Exit(0)
	}()

	for {
		buf, n := readInput()

		if n > 0 {

			fmt.Printf("%s", RESET)
			cols, rows := getDimensions()

			grid := getFilledChatWindow(make([]string, 0), int(rows), int(cols))
			screen := gridToString(grid)
			fmt.Print(screen)

		}

		if buf[0] == 4 {
			fmt.Println("Ctrl+D pressed, exiting...")
			return
		}
	}
}
