# hello-go Agent Guide

## Project intent
- Keep this as an advanced, terminal-first "Hello, world" experience.
- Prioritize cinematic ASCII output that is readable and stable in iTerm and other modern TTYs.
- Preserve graceful degradation: if network or TTY features are missing, the app should still run.

## Runtime guarantees
- Live mode must support clean termination via `Ctrl+C` and always restore terminal state.
- Raw-input mode must keep single-key controls responsive without breaking quit behavior.
- Rendering should remain ASCII-only unless explicitly requested otherwise.

## Architecture boundaries
- CLI wiring: `main.go` and `cmd/hello-go/main.go`.
- Demo runtime loop, controls, and terminal UI: `internal/demo/demo.go`, `internal/demo/input.go`.
- Rendering and visual effects: `internal/demo/render.go`, `internal/demo/anim.go`, `internal/demo/viewport.go`.
- External data/API access: `internal/demo/api.go`.

## Best practices for future agents
- Prefer small, focused changes; avoid broad rewrites of the render loop.
- Add tests for control/input behavior when changing key handling.
- Validate both compact and cinematic modes when touching viewport/render logic.
- Keep dependency footprint minimal (standard library first).
- Update `README.md` flags and controls when behavior changes.
- Do not remove fallback behavior for non-interactive stdin or small terminals.
