package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// 48914340
func find(d []int, x int) []int {
	indexes := []int{}
	for i := 0; i < x; i++ {
		indexes = append(indexes, i)
	}

	ret := 0

	nums := []int{}

	for ret == 0 {
		sum := 0
		for _, v := range indexes {
			sum += d[v]
		}

		if sum == 2020 {
			fmt.Println(sum)
		}

	}

	for i, v := range d {
		for ii := i + 1; ii < len(d); ii++ {
			for iii := ii + 1; iii < len(d); iii++ {
				if v+d[ii]+d[iii] == 2020 {
					return []int{v, d[ii], d[iii]}
				}
			}
		}
	}

	return nums
}

func main() {
	buf, err := ioutil.ReadFile("./data.txt")
	if err != nil {
		panic(err)
	}

	txt := string(buf)
	splits := strings.Split(txt, "\n")

	var nums []int
	for _, v := range splits {
		num, _ := strconv.Atoi(v)
		nums = append(nums, num)
	}

	res := find(nums, 3)

	ret := 1
	for _, v := range res {
		ret = ret * v
	}

	fmt.Println(ret)
}
