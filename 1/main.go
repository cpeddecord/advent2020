package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type line struct {
	min      int
	max      int
	letter   string
	password string
}

// 17-19 p: pwpzpfbrcpppjppbmppp
func parseText(t string) []line {
	re := regexp.MustCompile(`(\d*)-(\d*) (\w): (\w*)`)
	splits := strings.Split(t, "\n")
	ret := []line{}

	for _, v := range splits {
		t := re.FindStringSubmatch(v)

		if len(t) == 5 {
			min, _ := strconv.Atoi(string(t[1]))
			max, _ := strconv.Atoi(string(t[2]))
			ret = append(ret, line{
				min,
				max,
				string(t[3]),
				string(t[4]),
			})
		}
	}

	return ret
}

// TODO: bail out if max exceeded
func validateRuleOne(l line) int {
	count := 0
	for _, v := range l.password {
		if string(v) == l.letter {
			count++
		}
	}

	if count <= l.max && count >= l.min {
		return 1
	}

	return 0
}

func validateRuleTwo(l line) int {
	count := 0
	if string(l.password[l.min-1]) == l.letter {
		count++
	}
	if string(l.password[l.max-1]) == l.letter {
		count++
	}

	if count == 2 || count == 0 {
		return 0
	}

	return 1
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	lines := parseText(string(buf))

	validRuleOnes := 0
	validRuleTwos := 0
	for _, l := range lines {
		validRuleOnes += validateRuleOne(l)
		validRuleTwos += validateRuleTwo(l)
	}

	fmt.Println(validRuleOnes)
	fmt.Println(validRuleTwos)
}
