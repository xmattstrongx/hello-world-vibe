package demo

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
)

func defaultControls() ControlState {
	return ControlState{
		Meteors:   true,
		Aurora:    true,
		Trail:     true,
		Pulses:    true,
		Scanlines: true,
	}
}

func startInputListener(in io.Reader) <-chan string {
	ch, _, _ := startInputListenerWithCleanup(in)
	return ch
}

func startInputListenerWithCleanup(in io.Reader) (<-chan string, func(), bool) {
	ch := make(chan string, 16)
	if in == nil {
		close(ch)
		return ch, func() {}, false
	}

	if f, ok := in.(*os.File); ok {
		if restore, ok := enableRawMode(f); ok {
			go readBytes(ch, f)
			return ch, restore, true
		}
	}

	go func() {
		defer close(ch)
		sc := bufio.NewScanner(in)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			select {
			case ch <- strings.ToLower(line):
			default:
			}
		}
	}()
	return ch, func() {}, false
}

func applyCommand(ctrl *ControlState, cmd string) {
	switch cmd {
	case " ", "space", "p", "pause":
		ctrl.Paused = !ctrl.Paused
	case "m", "meteor", "meteors":
		ctrl.Meteors = !ctrl.Meteors
	case "a", "aurora":
		ctrl.Aurora = !ctrl.Aurora
	case "t", "trail":
		ctrl.Trail = !ctrl.Trail
	case "c", "city", "pulse", "pulses":
		ctrl.Pulses = !ctrl.Pulses
	case "s", "scan", "scanline", "scanlines":
		ctrl.Scanlines = !ctrl.Scanlines
	case "z", "zoom", "solar":
		ctrl.SolarSystem = !ctrl.SolarSystem
	case "q", "quit", "exit":
		ctrl.Quit = true
	}
}

func readBytes(ch chan<- string, in io.Reader) {
	defer close(ch)
	r := bufio.NewReader(in)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return
		}
		cmd := byteToCommand(b)
		if cmd == "" {
			continue
		}
		select {
		case ch <- cmd:
		default:
		}
	}
}

func byteToCommand(b byte) string {
	switch b {
	case ' ':
		return "space"
	case '\n', '\r', '\t':
		return ""
	default:
		if b >= 'A' && b <= 'Z' {
			return strings.ToLower(string([]byte{b}))
		}
		if b >= 'a' && b <= 'z' {
			return string([]byte{b})
		}
	}
	return ""
}

func enableRawMode(f *os.File) (func(), bool) {
	get := exec.Command("stty", "-g")
	get.Stdin = f
	saved, err := get.Output()
	if err != nil {
		return func() {}, false
	}
	mode := strings.TrimSpace(string(saved))
	if mode == "" {
		return func() {}, false
	}

	set := exec.Command("stty", "raw", "-echo")
	set.Stdin = f
	if err := set.Run(); err != nil {
		return func() {}, false
	}

	restore := func() {
		cmd := exec.Command("stty", mode)
		cmd.Stdin = f
		_ = cmd.Run()
	}
	return restore, true
}
