/*
Baseline but read into a large buffer and process instead of line reader
*/
package main

import (
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
	stations := make(map[string]*Station)
	buffer := make([]byte, 4096*1024)
	startReadAt := 0
	for {
		count, err := file.Read(buffer[startReadAt:])
		if err == io.EOF && count == 0 {
			break
		}
		start := 0
		for {
			newLinePos := bytes.IndexByte(buffer[start:count+startReadAt], '\n')
			if newLinePos == -1 {
				break
			}
			name, temp := parseLine(buffer[start : start+newLinePos])
			nameStr := string(name)
			start += newLinePos + 1
			s := stations[nameStr]
			if s != nil {
				s.min = min(s.min, temp)
				s.max = max(s.max, temp)
				s.sum += temp
				s.count += 1
			} else {
				stations[nameStr] = &Station{
					min:   temp,
					max:   temp,
					sum:   temp,
					count: 1,
				}
			}
		}

		//check if we ended before all the bytes read
		if start != count+startReadAt {
			remainder := buffer[start:]
			var i int
			var b byte
			for i, b = range remainder {
				buffer[i] = b
			}
			startReadAt = i + 1
		} else {
			startReadAt = 0
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

func parseLine(line []byte) ([]byte, FixedPoint) {
	semiPos := bytes.IndexByte(line, ';')
	name := line[:semiPos]
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
