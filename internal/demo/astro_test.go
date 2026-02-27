package demo

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	t.Parallel()
	if got := FormatDuration(3661 * time.Second); got != "01h 01m 01s" {
		t.Fatalf("unexpected format: %q", got)
	}
	if got := FormatDuration(-1 * time.Second); got != "00h 00m 00s" {
		t.Fatalf("unexpected negative format: %q", got)
	}
}

func TestSubsolarLongitude(t *testing.T) {
	t.Parallel()
	noon := time.Date(2026, 2, 27, 12, 0, 0, 0, time.UTC)
	midnight := time.Date(2026, 2, 27, 0, 0, 0, 0, time.UTC)

	if got := SubsolarLongitude(noon); got != 0 {
		t.Fatalf("noon expected 0 got %v", got)
	}
	if got := SubsolarLongitude(midnight); got != 180 {
		t.Fatalf("midnight expected 180 got %v", got)
	}
}

func TestDaylightAtEquinox(t *testing.T) {
	t.Parallel()
	decl := 0.0
	subLon := 0.0

	if !IsDaylight(0, 0, decl, subLon) {
		t.Fatal("subsolar point must be daylight")
	}
	if IsDaylight(0, 181, decl, subLon) {
		t.Fatal("anti-solar point must be night")
	}
}

func TestProjectVisibility(t *testing.T) {
	t.Parallel()
	if _, _, ok := Project(0, 0, 0, 64, 26); !ok {
		t.Fatal("front point should be visible")
	}
	if _, _, ok := Project(0, 180, 0, 64, 26); ok {
		t.Fatal("back point should not be visible")
	}
}
