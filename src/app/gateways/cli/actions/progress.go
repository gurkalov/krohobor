package actions

import (
	"fmt"
	"os"
)

// S3ProgressPrinter returns a callback that renders an upload progress bar
// to stderr. It is used to track the progress of loading a dump to S3 storage.
func S3ProgressPrinter() func(done, total int) {
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
			fmt.Fprintf(os.Stderr, "\rUploading: [%s] %5.1f%%", string(bar), pct*100)
		}
		if done >= total && !finished {
			finished = true
			fmt.Fprintln(os.Stderr)
		}
	}
}
