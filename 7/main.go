package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type instruction struct {
	name     string
	quantity int
	lineNo   int
}

func parseText(s string) []instruction {
	ret := []instruction{}

	re := regexp.MustCompile(`(\w{3}) ((?:\+|\-)\d*)`)

	strs := re.FindAllStringSubmatch(s, -1)
	for i, line := range strs {
		quantity, _ := strconv.Atoi(line[2])
		ret = append(ret, instruction{
			name:     line[1],
			quantity: quantity,
			lineNo:   i,
		})
	}

	return ret
}

func runInstruction(inst instruction, acc *int, set map[int]int, order *[]int, next *int) error {
	if v, ok := set[inst.lineNo]; ok {
		if v == -1 {
			return errors.New("duplicate instruction found")
		}
	}

	*order = append(*order, inst.lineNo)
	set[inst.lineNo] = -1

	switch op := inst.name; op {
	case "acc":
		*acc += inst.quantity
		*next++
	case "jmp":
		*next += inst.quantity
	case "nop":
		*next++
	}

	return nil
}

func runInstructions(instructionSet []instruction) (int, instruction, error) {
	acc := 0
	currentInstructionLine := 0
	linesExecuted := make(map[int]int)
	executionOrder := []int{}

	for currentInstructionLine <= len(instructionSet) {
		err := runInstruction(
			instructionSet[currentInstructionLine],
			&acc,
			linesExecuted,
			&executionOrder,
			&currentInstructionLine,
		)

		if err != nil {
			lastInstruction := instructionSet[executionOrder[len(executionOrder)-1]]
			return acc, lastInstruction, err
		}
	}

	return acc, instruction{}, nil
}

// TODO: optimize by re-executing from error location
func runInstructionsAndFix(instructionSet []instruction) (int, error) {
	_, lastInstruction, _ := runInstructions(instructionSet)

	newInstructions := make([]instruction, len(instructionSet))
	copy(newInstructions, instructionSet)

	re := strings.NewReplacer("jmp", "nop", "nop", "jmp")
	newInstructions[lastInstruction.lineNo] = instruction{
		lineNo:   lastInstruction.lineNo,
		quantity: lastInstruction.quantity,
		name:     re.Replace(lastInstruction.name),
	}

	finalAcc, _, err := runInstructions(newInstructions)

	if err == nil {
		return finalAcc, nil
	}

	return 0, errors.New("oh noes")
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	instructionSet := parseText(string(buf))

	acc, _, _ := runInstructions(instructionSet)
	fmt.Println(acc)

	finalAcc, _ := runInstructionsAndFix(instructionSet)
	fmt.Println(finalAcc)
}
