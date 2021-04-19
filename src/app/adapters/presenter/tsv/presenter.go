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
	} else if v, ok := val.([]domain.Table); ok {
		for _, record := range v {
			records = append(records, []string{record.Name, strconv.Itoa(record.Count), strconv.Itoa(record.Size)})
		}
	}

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
