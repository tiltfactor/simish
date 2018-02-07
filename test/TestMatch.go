package test

import (
	"fmt"
	"sort"

	"github.com/tiltfactor/simish/domain"
)

// RunSoftMatch tests the softmatch algorithm with the input string against the
// pairs given. If no input is given, each of the pairs will be used as an input.
func RunSoftMatch(args []string, allPairs []domain.InputOutput) {
	if len(args) == 2 {
		input := args[0]
		pairs := []domain.InputOutput {domain.NewInputOutput(args[1], "", 1)}
		bestMatch, score := domain.SoftMatch(args[0], pairs)
		fmt.Printf("Input:\t\t%v\n", input)
		fmt.Printf("Matched:\t%v\n", bestMatch.Input)
		fmt.Printf("Score:\t\t%v\n\n", score)

	} else if len(args) == 1 {
		input := args[0]
		bestMatch, score := domain.SoftMatch(input, allPairs)
		fmt.Printf("Input:\t\t%v\n", input)
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
