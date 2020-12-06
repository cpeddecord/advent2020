package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

func parseText(s string) map[string]string {
	list := strings.Split(s, "\n")
	sort.Strings(list)

	ret := map[string]string{}
	for _, v := range list {
		if v == "" {
			continue
		}
		ret[v] = v
	}

	return ret
}

// TODO: iterate row/column with length-based exponent
func tixToCoords(s string) []int {
	row := 128
	col := 8
	for i, r := range s {
		v := string(r)
		if i <= 6 {
			rowNum := 128 / math.Pow(2, float64(i+1))
			if v == "F" {
				row -= int(rowNum)
			}
			continue
		}

		colNum := 8 / math.Pow(2, float64(i-6))
		if v == "L" {
			col -= int(colNum)
		}
	}

	row--
	col--

	return []int{row, col}
}

func coordsToInt(n []int) int {
	return (n[0] * 8) + n[1]
}

func tixToInt(s string) int {
	return coordsToInt(tixToCoords(s))
}

func findHighestID(m map[string]string) (string, int) {
	s := ""
	n := 0
	for _, v := range m {
		tixID := tixToInt(v)
		if tixID > n {
			n = tixID
			s = v
		}
	}

	return s, n
}

func findMissingSeat(m map[string]string) int {
	seatIDs := []int{}
	for _, v := range m {
		seatIDs = append(seatIDs, tixToInt(v))
	}

	sort.Ints(seatIDs)

	missingSeat := 0
	for i, v := range seatIDs {
		if i+1 >= len(seatIDs) {
			break
		}

		if v+1 != seatIDs[i+1] {
			missingSeat = v + 1
			break
		}
	}

	return missingSeat
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	list := parseText(string(buf))
	_, num := findHighestID(list)
	seat := findMissingSeat(list)

	fmt.Println(num)
	fmt.Println(seat)
}
