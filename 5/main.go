package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Group struct {
	raw         string
	peopleCount int
	questions   map[string]int
}

func parseText(s string) []Group {
	splits := strings.Split(s, "\n\n")

	ret := []Group{}
	for _, g := range splits {
		peepCount := 1
		questions := map[string]int{}
		for _, v := range strings.TrimSpace(g) {
			char := string(v)

			if char == "\n" {
				peepCount++
				continue
			}

			questions[char]++
		}

		ret = append(ret, Group{
			g,
			peepCount,
			questions,
		})
	}

	return ret
}

func addGroupQuestions(g []Group) int {
	ret := 0
	for _, v := range g {
		ret += len(v.questions)
	}

	return ret
}

func getGroupCoherencyScore(g Group) int {
	ret := 0

	for _, v := range g.questions {
		if v == g.peopleCount {
			ret++
		}
	}

	return ret
}

func addGlobalCoherencyScores(g []Group) int {
	scores := 0

	for _, v := range g {
		scores += getGroupCoherencyScore(v)
	}

	return scores
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	groups := parseText(string(buf))

	num := addGroupQuestions(groups)
	scores := addGlobalCoherencyScores(groups)

	fmt.Println(num)
	fmt.Println(scores)
}
