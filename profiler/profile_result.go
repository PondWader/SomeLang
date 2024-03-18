package profiler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ProfileResult struct {
	Name        string
	Duration    time.Duration
	SubPrograms []*ProfileResult
}

// Uses insertion sort to sort the durations of the subprograms from slowest to fastest.
func (pr *ProfileResult) SortSubPrograms() []*ProfileResult {
	subPrograms := make([]*ProfileResult, len(pr.SubPrograms))
	copy(subPrograms, pr.SubPrograms)

	for i := 1; i < len(subPrograms); i++ {
		value := subPrograms[i]
		index := i
		for index > 0 && value.Duration > subPrograms[index-1].Duration {
			subPrograms[index] = subPrograms[index-1]
			index -= 1
		}
		subPrograms[index] = value
	}

	return subPrograms
}

func (pr *ProfileResult) ToSortedStringFormat(indentAmount int) string {
	indent := ""
	for i := 0; i < indentAmount; i++ {
		indent += "\t"
	}
	output := indent + "- " + pr.Name + " " + pr.Duration.String()
	if len(pr.SubPrograms) > 0 {
		output += ":\n"
		subPrograms := pr.SortSubPrograms()
		for _, subProgram := range subPrograms {
			output += subProgram.ToSortedStringFormat(indentAmount + 1)
		}
	} else {
		output += "\n"
	}

	return output
}

func (pr *ProfileResult) ToCsv() string {
	csv, _ := pr.generateCsvOutput(1, 0)
	return csv
}

func (pr *ProfileResult) generateCsvOutput(index int, parentIndex int) (string, int) {
	csv := fmt.Sprint(parentIndex) + "," + pr.Name + "," + fmt.Sprint(pr.Duration.Nanoseconds())
	currentIndex := index
	for _, subProgram := range pr.SubPrograms {
		var output string
		output, currentIndex = subProgram.generateCsvOutput(currentIndex+1, index)
		csv += "\n" + output
	}
	return csv, currentIndex
}

func ParseCsv(csv string) (*ProfileResult, error) {
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	results := make([]*ProfileResult, len(lines))

	for i, line := range lines {
		values := strings.Split(line, ",")

		if len(values) != 3 {
			return nil, errors.New("profiler CSV data is invalid (too many values on line)")
		}
		parentIndex, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, errors.New("profiler CSV data is invalid (failed to read parent index integer)")
		}
		durationNanoseconds, err := strconv.Atoi(values[2])
		if err != nil {
			return nil, errors.New("profiler CSV data is invalid (failed to read duration integer)")
		}

		result := &ProfileResult{
			Name:        values[1],
			Duration:    time.Nanosecond * time.Duration(durationNanoseconds),
			SubPrograms: make([]*ProfileResult, 0),
		}
		results[i] = result

		if parentIndex != 0 {
			if len(results) <= parentIndex {
				return nil, errors.New("profiler CSV data is invalid (parent index out of bounds)")
			}
			parent := results[parentIndex-1]
			if parent == nil {
				return nil, errors.New("profiler CSV data is invalid (could not find parent with index)")
			}
			parent.SubPrograms = append(parent.SubPrograms, result)
		}
	}

	return results[0], nil
}
