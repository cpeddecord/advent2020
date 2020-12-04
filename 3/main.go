package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// TODO: separate required/optional fields for easier iteration
type Passport map[string]string

// TODO: be better at regexp, otherwise you get this gross-ass code
func parseText(t string) []Passport {
	ret := []Passport{}

	kv := regexp.MustCompile(`(\w{3}):([\w\d#]*)`)
	boundaries := regexp.MustCompile(`\n|\s`)

	// passports separated by double new line
	splits := strings.Split(t, "\n\n")
	for _, recordStr := range splits {
		passport := Passport{}
		record := boundaries.Split(recordStr, -1)
		for _, v := range record {
			if v == "" {
				continue
			}
			field := kv.FindStringSubmatch(v)
			passport[field[1]] = field[2]
		}

		ret = append(ret, passport)
	}

	return ret
}

func ruleOne(p []Passport) []Passport {
	optionalFields := []string{"pid", "cid"}
	requiredFields := []string{
		"byr",
		"iyr",
		"eyr",
		"hgt",
		"hcl",
		"ecl",
	}

	validPassports := []Passport{}
	// TODO: this is bad and you should feel bad
	for _, v := range p {
		validFieldsCount := 0
		for _, f := range requiredFields {
			if v[f] == "" {
				break
			}

			validFieldsCount++
		}

		if validFieldsCount != len(requiredFields) {
			continue
		}

		isValid := false
		for _, f := range optionalFields {
			if v[f] == "" {
				break
			}

			isValid = true
		}

		if isValid == true {
			validPassports = append(validPassports, v)
		}

	}

	return validPassports
}

func ruleTwo(p []Passport) []Passport {
	hgtRe := regexp.MustCompile(`(\d*)(cm|in)`)
	hclRe := regexp.MustCompile(`#([a-f\d]{6})$`)
	pidRe := regexp.MustCompile(`^\d{9}$`)

	validatorsReq := map[string]func(string) bool{
		// byr (Birth Year) - four digits; at least 1920 and at most 2002.
		"byr": func(s string) bool {
			year, _ := strconv.Atoi(s)
			return year >= 1920 && year <= 2002
		},
		// iyr (Issue Year) - four digits; at least 2010 and at most 2020.
		"iyr": func(s string) bool {
			year, _ := strconv.Atoi(s)
			return year >= 2010 && year <= 2020
		},
		// eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
		"eyr": func(s string) bool {
			year, _ := strconv.Atoi(s)
			return year >= 2020 && year <= 2030
		},
		/*
			hgt (Height) - a number followed by either cm or in:
			If cm, the number must be at least 150 and at most 193.
			If in, the number must be at least 59 and at most 76.
		*/
		"hgt": func(s string) bool {
			if !hgtRe.MatchString(s) {
				return false
			}

			res := hgtRe.FindStringSubmatch(s)
			hgt, _ := strconv.Atoi(res[1])
			unit := res[2]

			if unit == "in" {
				return hgt >= 59 && hgt <= 76
			}

			if unit == "cm" {
				return hgt >= 150 && hgt <= 193
			}

			return false
		},
		// hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
		"hcl": func(s string) bool {
			isValid := hclRe.MatchString(s)
			return isValid
		},
		// ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth.
		"ecl": func(s string) bool {
			for _, v := range []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"} {
				if s == v {
					return true
				}
			}

			return false
		},
	}

	ret := []Passport{}
	for _, v := range p {
		validReqFields := 0
		validPid := pidRe.MatchString(v["pid"])
		if !validPid {
			continue
		}

		for k, f := range v {
			if validatorsReq[k] != nil && validatorsReq[k](f) {
				validReqFields++
			}
		}

		if validReqFields == 6 {
			ret = append(ret, v)
		}

	}

	return ret
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	passports := parseText(string(buf))
	validOnes := ruleOne(passports)
	validTwos := ruleTwo(validOnes)

	fmt.Println(len(validOnes))
	fmt.Println(len(validTwos))
}
