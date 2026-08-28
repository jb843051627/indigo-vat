package codec

import (
	"encoding/csv"
	"github.com/jb843051627/indigo-vat/internal/model"
	"io"
	"sort"
	"strconv"
)

func WriteSamples(w io.Writer, samples []model.Sample, locationName string) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{"sample_id", "taken_at", "hue", "ph", "temperature", "status", "observer"}); err != nil {
		return err
	}
	location := Location(locationName)
	ordered := make([]model.Sample, len(samples))
	copy(ordered, samples)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].TakenAt.Before(ordered[j].TakenAt) })
	for _, sample := range ordered {
		if err := writer.Write([]string{sample.ID, Format(sample.TakenAt, location), strconv.FormatFloat(sample.Hue, 'f', 2, 64), strconv.FormatFloat(sample.PH, 'f', 2, 64), strconv.FormatFloat(sample.Temperature, 'f', 2, 64), sample.Status, sample.Observer}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
