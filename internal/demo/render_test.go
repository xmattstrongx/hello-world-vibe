package demo

import (
	"strings"
	"testing"
	"time"
)

func TestRenderFrameAvoidsSingleLineCrop(t *testing.T) {
	t.Parallel()

	out := RenderFrame(FrameData{
		Now:      time.Date(2026, 3, 10, 14, 53, 2, 0, time.UTC),
		Lang:     "en",
		ViewMode: "cinematic",
		Controls: defaultControls(),
	}, 120, 1)

	if lines := strings.Count(out, "\n") + 1; lines < 10 {
		t.Fatalf("expected multi-line frame, got %d lines", lines)
	}
	if !strings.Contains(out, "HELLO, WORLD FROM SPACE") {
		t.Fatal("expected frame header in output")
	}
}
