/*
This is meant to be the most naive and straight forward implemenation that passes the tests
*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

type FixedPoint int

// / Convert out fixed 2 decimal type to a f64 with 1 decimal place thats properly rounded
func roundToFloat(num FixedPoint) float64 {
	hund := num % 10
	num = num / 10
	if hund >= 5 {
		num += 1
	}
	return float64(num) / 10.0
}

type Station struct {
	sum, min, max FixedPoint
	count         uint
}

func stationToString(name string, station Station) string {
	nMin := float64(station.min) / 100.0
	nMax := float64(station.max) / 100.0
	nAvg := roundToFloat(station.sum / FixedPoint(station.count))
	return fmt.Sprintf("%s=%.1f/%.1f/%.1f", string(name), nMin, nAvg, nMax)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func min(a, b FixedPoint) FixedPoint {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b FixedPoint) FixedPoint {
	if a > b {
		return a
	} else {
		return b
	}
}

func main() {
	parse(os.Stdout)
}

func parse(output io.Writer) {
	file, err := os.Open("./measurements.txt")
	check(err)
	defer file.Close()
	reader := bufio.NewScanner(file)
	stations := make(map[string]*Station)
	for reader.Scan() {
		name, temp := parseLine(reader.Bytes())
		s := stations[name]
		if s != nil {
			s.min = min(s.min, temp)
			s.max = max(s.max, temp)
			s.sum += temp
			s.count += 1
		} else {
			stations[name] = &Station{
				min:   temp,
				max:   temp,
				sum:   temp,
				count: 1,
			}
		}
	}
	names := make([]string, 0, len(stations))
	for name := range stations {
		names = append(names, name)
	}

	sort.Strings(names)
	fmt.Fprint(output, "{")
	first := true
	for _, name := range names {
		if !first {
			fmt.Fprint(output, ", ")
		} else {
			first = false
		}
		fmt.Fprint(output, stationToString(name, *stations[name]))
	}
	fmt.Fprintln(output, "}")
}

func parseLine(line []byte) (string, FixedPoint) {
	semiPos := bytes.IndexByte(line, ';')
	name := string(line[:semiPos])
	temp := parseBytesToFixedPoint(line[semiPos+1:])
	return name, temp
}

func parseBytesToFixedPoint(in []byte) FixedPoint {
	result := 0
	neg := 1
	for i := 0; i < len(in); i++ {
		if in[i] == '-' {
			neg = -1
		} else if in[i] >= '0' && in[i] <= '9' {
			digit := int(in[i] - '0')
			result = result*10 + (digit)
		}
	}

	return FixedPoint(result * neg * 10)
}
