package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// .....#............#....#####.##
func parseText(t string) [][]int {
	splits := strings.Split(t, "\n")

	ret := [][]int{}
	for _, l := range splits {
		if len(l) == 0 {
			continue
		}

		row := []int{}
		for _, c := range l {
			val := 0
			if string(c) == "#" {
				val = 1
			}

			row = append(row, val)
		}
		ret = append(ret, row)
	}

	return ret
}

func getValueFromGrid(g [][]int, coords map[string]int) int {
	row := g[coords["x"]]
	multiplier := coords["y"] / len(row)
	relativeY := coords["y"] - (multiplier * len(row))

	return g[coords["x"]][relativeY]
}

func getTobogganTreeStrikes(g [][]int, x int, y int) int {
	coords := map[string]int{"x": 0, "y": 0}

	treeStrikes := 0
	for len(g) > coords["x"] {
		treeStrikes += getValueFromGrid(g, coords)

		coords["x"] += x
		coords["y"] += y
	}

	return treeStrikes
}

func main() {
	buf, err := ioutil.ReadFile("./data2.txt")
	if err != nil {
		panic(err)
	}

	lines := parseText(string(buf))

	firstRun := getTobogganTreeStrikes(lines, 1, 3)
	fmt.Println(firstRun)

	ruleSet := [][]int{
		[]int{1, 1},
		[]int{1, 3},
		[]int{1, 5},
		[]int{1, 7},
		[]int{2, 1},
	}

	trees := 1
	for _, v := range ruleSet {
		trees = trees * getTobogganTreeStrikes(lines, v[0], v[1])
	}

	fmt.Println(trees)

}
