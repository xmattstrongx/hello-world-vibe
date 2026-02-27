package demo

import (
	"io"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

type ttyWinsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func terminalSize(out io.Writer) (int, int) {
	if f, ok := out.(*os.File); ok {
		ws := &ttyWinsize{}
		_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
		if errno == 0 && ws.Col > 0 && ws.Row > 0 {
			return int(ws.Col), int(ws.Row)
		}
	}
	if c, r, ok := sizeFromEnv(); ok {
		return c, r
	}
	return 100, 36
}

func sizeFromEnv() (int, int, bool) {
	cols, errC := strconv.Atoi(os.Getenv("COLUMNS"))
	rows, errR := strconv.Atoi(os.Getenv("LINES"))
	if errC != nil || errR != nil || cols <= 0 || rows <= 0 {
		return 0, 0, false
	}
	return cols, rows, true
}

func resolveViewMode(termW, termH int, selected string) string {
	if selected == "compact" || selected == "cinematic" {
		return selected
	}
	if termW < 100 || termH < 30 {
		return "compact"
	}
	return "cinematic"
}

func globeViewport(termW, termH int, mode string) (width, height, leftPad int) {
	if termW < 40 {
		termW = 40
	}
	if termH < 16 {
		termH = 16
	}

	maxHeight := termH - 9
	if mode == "compact" {
		maxHeight = termH - 7
	}
	if mode == "cinematic" {
		maxHeight = termH - 11
	}
	if maxHeight < 8 {
		maxHeight = 8
	}
	if maxHeight > 42 {
		maxHeight = 42
	}

	maxWidth := termW - 2
	if maxWidth < 30 {
		maxWidth = 30
	}

	// Keep the globe roughly proportional in monospace terminals.
	width = maxHeight * 3
	if mode == "compact" {
		width = maxHeight * 2
	}
	if mode == "cinematic" {
		width = maxHeight * 4
	}
	if width > maxWidth {
		width = maxWidth
	}
	if width < 30 {
		width = 30
	}

	height = width / 3
	if height > maxHeight {
		height = maxHeight
	}
	if height < 8 {
		height = 8
	}

	leftPad = (termW - width) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	return width, height, leftPad
}
