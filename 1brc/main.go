package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Measurement struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int64
}

func (m *Measurement) Average() float64 {
	return m.Sum / float64(m.Count)
}

func main() {
	start := time.Now()

	measurements, err := os.Open("data/measurements.txt")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println(
				"Please run `./setup.sh <number>` to generate necessary files, " +
					"then run program again")
			return
		}
		panic(err)
	}
	defer measurements.Close()

	const separator = ";"
	data := make(map[string]Measurement)
	scanner := bufio.NewScanner(measurements)
	for scanner.Scan() {
		line := scanner.Text()
		separatorIndex := strings.Index(line, separator)
		location := line[:separatorIndex]
		rawTemperature := line[separatorIndex+1:]
		temperature, _ := strconv.ParseFloat(rawTemperature, 64)

		measurement, ok := data[location]
		if ok {
			measurement.Min = min(temperature, measurement.Min)
			measurement.Max = max(temperature, measurement.Max)
			measurement.Sum += temperature
			measurement.Count++
		} else {
			measurement = Measurement{
				Min:   temperature,
				Max:   temperature,
				Sum:   temperature,
				Count: 1,
			}
		}

		data[location] = measurement
	}

	locations := make([]string, 0, len(data))

	for name := range data {
		locations = append(locations, name)
	}

	sort.Strings(locations)

	fmt.Print("{ ")
	for _, location := range locations {
		measurement := data[location]
		fmt.Printf(
			"%s=%.1f/%.1f/%.1f, ",
			location,
			measurement.Min,
			measurement.Average(),
			measurement.Max,
		)
	}
	fmt.Print("}\n")

	fmt.Printf("Elapsed: %s\n", time.Since(start).String())
}
