package main

import (
	"flag"
	"fmt"
	"os"

	"hello-go/internal/demo"
)

func main() {
	cfg := demo.Config{}
	var compact bool
	var cinematic bool
	flag.BoolVar(&cfg.Live, "live", false, "run the live Hello World from Space demo")
	flag.StringVar(&cfg.Lang, "lang", "auto", "language mode: auto or en")
	flag.StringVar(&cfg.View, "view", "auto", "view mode: auto, compact, cinematic")
	flag.BoolVar(&compact, "compact", false, "shortcut for --view compact")
	flag.BoolVar(&cinematic, "cinematic", false, "shortcut for --view cinematic")
	flag.BoolVar(&cfg.MaxASCII, "max-ascii", false, "maximize ASCII detail (auto-enabled in iTerm)")
	flag.IntVar(&cfg.Frames, "frames", 0, "number of frames to render (0 = infinite)")
	flag.DurationVar(&cfg.Interval, "interval", demo.DefaultInterval, "frame interval (example: 700ms, 1s)")
	flag.Parse()

	if compact && cinematic {
		fmt.Fprintln(os.Stderr, "error: use either --compact or --cinematic, not both")
		os.Exit(2)
	}
	if compact {
		cfg.View = "compact"
	}
	if cinematic {
		cfg.View = "cinematic"
	}

	if !cfg.Live {
		fmt.Println("Hello, world!")
		fmt.Println("Try: go run . --live")
		return
	}

	if err := demo.Run(os.Stdout, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
