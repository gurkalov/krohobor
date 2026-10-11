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

const logTimeLayout = "2006-01-02 15:04:05"

// LogOperationStart logs the start of an operation (e.g. "Restore", "Dump create").
func LogOperationStart(operation string) {
	fmt.Fprintf(os.Stderr, "%s  %s started\n", time.Now().Format(logTimeLayout), operation)
}

// LogOperationFinished logs the finish of an operation: the current time and the
// total elapsed time of the whole operation.
func LogOperationFinished(operation string, duration time.Duration) {
	fmt.Fprintf(os.Stderr, "%s  %s finished in %s\n", time.Now().Format(logTimeLayout), operation, duration.Round(time.Millisecond))
}

// LogDumpComplete logs the dump file written to disk: the current time and
// the size of the dump.
func LogDumpComplete(size int) {
	fmt.Fprintf(os.Stderr, "%s  Dumped %s\n", time.Now().Format(logTimeLayout), tsv.HumanSize(size))
}

// LogArchiveComplete logs the finish of archiving a dump before upload:
// the current time, the archived size and the time it took.
func LogArchiveComplete(size int, duration time.Duration) {
	fmt.Fprintf(os.Stderr, "%s  Archived %s in %s\n", time.Now().Format(logTimeLayout), tsv.HumanSize(size), duration.Round(time.Millisecond))
}

// LogDownloadComplete logs the finish of the S3 download between the download
// and restore steps: the current time and the total bytes downloaded.
func LogDownloadComplete(size int) {
	fmt.Fprintf(os.Stderr, "%s  Downloaded %s\n", time.Now().Format(logTimeLayout), tsv.HumanSize(size))
}

// LogUnarchiveComplete logs the finish of extracting the downloaded archive:
// the current time, the unarchived size and the time it took.
func LogUnarchiveComplete(size int, duration time.Duration) {
	fmt.Fprintf(os.Stderr, "%s  Unzipped %s in %s\n", time.Now().Format(logTimeLayout), tsv.HumanSize(size), duration.Round(time.Millisecond))
}
