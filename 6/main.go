package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

var table = make(map[string]node)

type node struct {
	id       string
	children map[string]child
	parents  map[string]string
}

type child struct {
	quantity int
	nodeID   string
}

var childRe = regexp.MustCompile(`((\d{1,2}) (.+?) (?:bag,|bags,|bag.|bags.|$))`)

func createNode(s string) {
	keyAndChildren := strings.Split(s, " bags contain ")
	childrenStrs := childRe.FindAllStringSubmatch(keyAndChildren[1], -1)

	nodeKey := strings.TrimSpace(keyAndChildren[0])
	nodeChildren := make(map[string]child)

	for _, c := range childrenStrs {
		key := ""
		value := 0
		for i := len(c); i > len(c)-2; i-- {
			if i == len(c) {
				key = strings.TrimSpace(c[i-1])
				continue
			}
			value, _ = strconv.Atoi(c[i-1])
		}

		nodeChildren[key] = child{value, key}

		if _, ok := table[key]; ok {
			table[key].parents[nodeKey] = nodeKey
		} else {
			table[key] = node{
				id:       key,
				parents:  map[string]string{nodeKey: nodeKey},
				children: make(map[string]child),
			}
		}
	}

	// well...
	if _, ok := table[nodeKey]; !ok {
		table[nodeKey] = node{
			id:       nodeKey,
			children: nodeChildren,
			parents:  make(map[string]string),
		}
	} else {
		table[nodeKey] = node{
			id:       nodeKey,
			children: nodeChildren,
			parents:  table[nodeKey].parents,
		}
	}
}

// TODO: optimize with regex
// TODO: be better at coding so you don't look like a caveman
func createNodes(s string) {
	splits := strings.Split(s, "\n")
	for _, n := range splits {
		if n == "" {
			continue
		}

		createNode(n)
	}
}

func findAllParents(bagColor string) map[string]string {
	n := table[bagColor]
	colors := map[string]string{}

	queue := []string{}
	for k := range n.parents {
		queue = append(queue, k)
	}

	// TODO: optimize this shit
	for len(queue) > 0 {
		current := queue[0]

		if _, ok := colors[current]; !ok {
			colors[current] = current
		}

		for k := range table[current].parents {
			if _, ok := colors[k]; !ok {

				queue = append(queue, k)
			}
		}

		queue = queue[1:]
	}

	return colors
}

type layer struct {
	mult int
	node node
}

func getChildrenNumber(l layer) (int, []layer) {
	mult := l.mult

	childQuant := 0
	toQueue := []layer{}
	for _, v := range l.node.children {
		childQuant += v.quantity
		toQueue = append(toQueue, layer{
			mult * v.quantity,
			table[v.nodeID],
		})
	}

	num := mult * childQuant
	return num, toQueue
}

func findChildrenNumber(bagName string) int {
	queue := []layer{layer{
		1,
		table[bagName],
	}}

	num := 0
	for len(queue) > 0 {
		toAdd, toQueue := getChildrenNumber(queue[0])

		num += toAdd
		queue = append(queue, toQueue...)
		queue = queue[1:]
	}

	return num
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	createNodes(string(buf))
	foundParents := findAllParents("shiny gold")
	fmt.Println(len(foundParents))

	childrenQuant := findChildrenNumber("shiny gold")
	fmt.Println(childrenQuant)
}
