package actions

import (
	"fmt"
	"os"
)

// S3ProgressPrinter returns a callback that renders a progress bar to stderr.
// label describes the operation, e.g. "Uploading" or "Downloading".
func S3ProgressPrinter(label string) func(done, total int) {
	const width = 30
	last := -1
	finished := false
	return func(done, total int) {
		if total <= 0 {
			return
		}
		pct := float64(done) / float64(total)
		filled := int(pct * width)
		if filled != last {
			last = filled
			bar := make([]rune, width)
			for i := range bar {
				if i < filled {
					bar[i] = '='
				} else {
					bar[i] = '.'
				}
			}
			fmt.Fprintf(os.Stderr, "\r%s: [%s] %5.1f%%", label, string(bar), pct*100)
		}
		if done >= total && !finished {
			finished = true
			fmt.Fprintln(os.Stderr)
		}
	}
}
