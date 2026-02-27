package demo

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	total := int(d.Seconds())
	if total < 0 {
		total = 0
	}
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	return fmt.Sprintf("%02dh %02dm %02ds", h, m, s)
}
