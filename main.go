package main

import (
	"flag"
	"fmt"
	"math/rand"
)

var (
	numIter  = flag.Int("iter", 10000, "number of iterations")
	faces    = flag.Int("d", 6, "die faces")
	groupSum = flag.Int("group", 10, "sum to remove from the dice")
	nMin     = flag.Int("nmin", 1, "minimum number of dice")
	nMax     = flag.Int("nmax", 0, "maximum number of dice")
	perm     = flag.Bool("perm", false, "use slower brute force algorithm")
)

func main() {
	flag.Parse()

	for numDice := *nMin; *nMax <= 0 || numDice <= *nMax; numDice++ {
		success := 0
		dice := make([]int, numDice)
		count := make([]int, *faces+1)
		for i := 0; i < *numIter; i++ {
			for c := range count {
				count[c] = 0
			}
			for d := 0; d < len(dice); d++ {
				dice[d] = rand.Intn(*faces) + 1
				count[dice[d]]++
			}

			var s bool
			if *perm {
				s = forEachPerm(len(dice), dice, func(dd []int) bool {
					sum := 0
					for _, d := range dd {
						sum += d
						if sum == *groupSum {
							sum = 0
							continue
						}
						if sum >= *groupSum {
							return false
						}
					}
					return true
				})
			} else {
				s = removeN(count, *groupSum)
			}
			if s {
				success++
			}
		}
		fmt.Printf("%d, %f\n", numDice, float64(success)/float64(*numIter))
	}
}

func removeN(count []int, n int) bool {
	sum := 0
	for i, c := range count {
		sum += i * c
	}
	for sum > n {
		if !tryRemove(count, n) {
			return false
		}
		sum -= n
	}
	return true
}

func tryRemove(count []int, n int) bool {
	if n == 0 {
		return true
	}
	i := len(count) - 1
	if i <= n && count[i] > 0 {
		count[i]--
		if tryRemove(count, n-i) {
			return true
		}
		count[i]++
	}
	if len(count) > 1 {
		return tryRemove(count[:i], n)
	}
	return false
}

func forEachPerm(k int, dice []int, fn func([]int) bool) bool {
	// https://en.wikipedia.org/wiki/Heap%27s_algorithm
	if k == 1 {
		return fn(dice)
	}

	if forEachPerm(k-1, dice, fn) {
		return true
	}

	for i := 0; i < k-1; i++ {
		if k%2 == 0 {
			dice[i], dice[k-1] = dice[k-1], dice[i]
		} else {
			dice[0], dice[k-1] = dice[k-1], dice[0]
		}
		if forEachPerm(k-1, dice, fn) {
			return true
		}
	}
	return false
}
