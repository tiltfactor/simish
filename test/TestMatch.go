package test

import (
	"fmt"
	"os"
	"sort"

	"github.com/tiltfactor/simish/domain"
)

func TestMatch(input string, allPairs []domain.InputOutput) {
	if input != "" {
		bestMatch, score := domain.SoftMatch(os.Args[2], allPairs)
		fmt.Printf("Input:\t\t%v\n", os.Args[2])
		fmt.Printf("Matched:\t%v\n", bestMatch.Input)
		fmt.Printf("Response:\t%v\n", bestMatch.Output)
		fmt.Printf("Score:\t\t%v\n\n", score)

	} else {
		matches := []Match{}
		for _, pair := range allPairs {
			pairs := []domain.InputOutput{}
			for _, filterPair := range allPairs {
				if pair.AiCol != filterPair.AiCol {
					pairs = append(pairs, filterPair)
				}
			}
			bestMatch, score := domain.SoftMatch(pair.Input, pairs)
			matches = append(matches,
				Match{pair.Input, bestMatch.Input, bestMatch.Output, score},
			)
		}

		sort.Sort(ByScore(matches))

		for _, match := range matches {
			fmt.Printf("Input:\t\t%v\n", match.Input)
			fmt.Printf("Matched:\t%v\n", match.InputMatch)
			fmt.Printf("Response:\t%v\n", match.Response)
			fmt.Printf("Score:\t\t%v\n\n", match.Score)
		}
	}
}
