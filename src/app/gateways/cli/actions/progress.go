package actions

import (
	"fmt"
	"krohobor/app/adapters/presenter/tsv"
	"os"
	"time"
)

// ProgressPrinter returns a callback that renders a progress bar to stderr.
// label describes the operation, e.g. "Uploading" or "Downloading".
func ProgressPrinter(label string) func(done, total int) {
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

// LogDownloadComplete logs the finish of the S3 download between the download
// and restore steps: the current time and the total bytes downloaded.
func LogDownloadComplete(size int) {
	fmt.Fprintf(os.Stderr, "%s  Downloaded %s\n", time.Now().Format("2006-01-02 15:04:05"), tsv.HumanSize(size))
}
