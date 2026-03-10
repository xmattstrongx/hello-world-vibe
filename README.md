# hello-go

A polished Go CLI demo that evolves `Hello, world!` into a live terminal experience:

- ASCII globe animation
- iTerm-friendly cinematic ASCII rendering
- Layered space animations (starfield, meteor showers, aurora)
- Daylight-aware city greetings
- ISS position marker (real-time) with fading orbit trail
- City pulse rings and scanline boot sequence
- Local next-sunrise countdown
- Responsive layout that adapts to terminal size

## Quickstart

```bash
go run .
go run . --live
```

Run a fixed number of frames for demos/CI:

```bash
go run . --live --frames 2
```

## CLI Flags

- `--live`: run the animated demo
- `--lang`: `auto` (city-local greetings) or `en`
- `--view`: `auto`, `compact`, or `cinematic`
- `--compact`: shortcut for compact view
- `--cinematic`: shortcut for cinematic view
- `--max-ascii`: force maximum ASCII detail (auto-enabled in iTerm)
- `--frames`: number of frames to render (`0` = infinite)
- `--interval`: frame interval (default `700ms`)

## Interactive Controls (Live Mode)

Single-key controls work instantly in a real terminal (raw mode). If stdin is not interactive, fallback is command + Enter.

- `p` or `space`: pause/resume animation
- `m`: toggle meteors
- `a`: toggle aurora
- `t`: toggle ISS trail
- `c`: toggle city pulses
- `s`: toggle scanlines
- `q`: quit
- `Ctrl+C`: quit cleanly and restore terminal state

## Project Layout

- `cmd/hello-go`: executable entrypoint
- `internal/demo`: core demo logic (rendering, astro math, API clients)
- `main.go`: tiny root entrypoint for `go run .`

## Development

```bash
make fmt
make test
make run-live
make build
```

## Notes

This project intentionally keeps dependencies to the Go standard library.
Network-backed features degrade gracefully if APIs are unavailable.
