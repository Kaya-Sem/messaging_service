package main

import (
	"fmt"
	"os"
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

func getWidth() uint {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return uint(ws.Col)
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

// readInput reads a keypress in raw mode
func readInput() ([]byte, int) {
	buf := make([]byte, 3)
	oldState, err := enableRawMode()
	if err != nil {
		fmt.Printf("Failed to enable raw mode: %v\n", err)
		return buf, -1
	}
	defer disableRawMode(oldState)

	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		fmt.Printf("Error reading key: %v\n", err)
		return buf, -1
	}
	return buf, n
}

func main() {
	fmt.Println("Press any key (Ctrl+C to exit):")

	fmt.Printf("%s", RESET)

	for {
		buf, n := readInput()

		fmt.Printf("%s", RESET)
		if n > 0 {
			width := getWidth()
			fmt.Printf("%s", "┏")
			fmt.Printf("%s", strings.Repeat("━", int(width-2)))
			fmt.Printf("%s", "┓")

			fmt.Printf("Captured: %v (bytes: %v)", string(buf[:n]), buf[:n])

			// Example: detect Escape key (27)
			if buf[0] == 27 {
				if n == 1 {
					fmt.Println("ESC key pressed")
				} else {
					// Special keys like arrows produce escape sequences
					fmt.Println("Special key sequence")
				}
			}

			// Example: detect Ctrl+D (4)
			if buf[0] == 4 {
				fmt.Println("Ctrl+D pressed, exiting...")
				return
			}
		}
	}
}
