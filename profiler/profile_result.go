package profiler

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type ProfileResult struct {
	Name        string
	Duration    time.Duration
	SubPrograms []*ProfileResult
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

func ParseCsv(csv string) *ProfileResult {
	lines := strings.Split(strings.TrimSpace(csv), "\n")
	results := make([]*ProfileResult, len(lines))

	for i, line := range lines {
		values := strings.Split(line, ",")

		if len(values) != 3 {
			log.Fatal("Profiler CSV data is invalid (too many values on line).")
		}
		parentIndex, err := strconv.Atoi(values[0])
		if err != nil {
			log.Fatal("Profiler CSV data is invalid (failed to read parent index integer).")
		}
		durationNanoseconds, err := strconv.Atoi(values[2])
		if err != nil {
			log.Fatal("Profiler CSV data is invalid (failed to read duration integer).")
		}

		result := &ProfileResult{
			Name:        values[1],
			Duration:    time.Nanosecond * time.Duration(durationNanoseconds),
			SubPrograms: make([]*ProfileResult, 0),
		}
		results[i] = result

		if parentIndex != 0 {
			if len(results) <= parentIndex {
				log.Fatal("Profiler CSV data is invalid (parent index out of bounds).")
			}
			parent := results[parentIndex-1]
			if parent == nil {
				log.Fatal("Profiler CSV data is invalid (could not find parent with index).")
			}
			parent.SubPrograms = append(parent.SubPrograms, result)
		}
	}

	return results[0]
}
