package tsv

import (
	"encoding/csv"
	"fmt"
	"krohobor/app/domain"
	"os"
	"strconv"
)

type Presenter struct {
}

func (p *Presenter) Print(val interface{}) error {
	records := make([][]string, 0)

	if v, ok := val.(string); ok {
		fmt.Println(v)
		return nil
	}

	if v, ok := val.([]string); ok {
		for _, record := range v {
			records = append(records, []string{record})
		}
	} else if v, ok := val.([]domain.Database); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, strconv.Itoa(record.Size)})
		}
	} else if v, ok := val.([]domain.Dump); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, strconv.Itoa(record.Size)})
		}
	} else if v, ok := val.([]domain.Table); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, strconv.Itoa(record.Count), strconv.Itoa(record.Size)})
		}
	}

	return p.write(records)
}

// PrintHuman prints records with human-readable sizes instead of raw bytes.
func (p *Presenter) PrintHuman(val interface{}) error {
	records := make([][]string, 0)

	if v, ok := val.(string); ok {
		fmt.Println(v)
		return nil
	}

	if v, ok := val.([]string); ok {
		for _, record := range v {
			records = append(records, []string{record})
		}
	} else if v, ok := val.([]domain.Database); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, HumanSize(record.Size)})
		}
	} else if v, ok := val.([]domain.Dump); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, HumanSize(record.Size)})
		}
	} else if v, ok := val.([]domain.Table); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, strconv.Itoa(record.Count), HumanSize(record.Size)})
		}
	}

	return p.write(records)
}

func (p *Presenter) write(records [][]string) error {
	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'
	for _, record := range records {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

// HumanSize formats a byte count as a human-readable string (e.g. "1.5 MB").
func HumanSize(size int) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := int64(size) / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
}
