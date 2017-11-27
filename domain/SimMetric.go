package domain

import (
	"regexp"
	"strings"

	textdistance "github.com/masatana/go-textdistance"
)

const wordWeight float64 = 0.7
const actionWeight float64 = 1
const distWeight float64 = 0.2
const upWeight float64 = 0
const downWeight float64 = 0

var reg = regexp.MustCompile("([.,?!*]|[[:blank:]])+")

func prepareInput(input string) []string {
	tokens := strings.Split(strings.ToLower(input), " ")

	// https://github.com/golang/go/wiki/SliceTricks#filtering-without-allocating
	filteredTokens := tokens[:0]
	for _, token := range tokens {
		// Filter out whitespace and some punctuation
		token = reg.ReplaceAllString(token, "")

		// If the token is empty, ignore it
		if len(token) == 0 {
			continue
		}

		// Check the token against our list of stop words
		_, ok := stopWords[token]
		if !ok {
			filteredTokens = append(filteredTokens, token)
		}
	}

	return filteredTokens
}

func rawScore(userTokens, dbTokens []string) float64 {
	matchedWords, matchedActions, totalActions := 0, 0, 0

	// The array counter and arrays will be used to track matches
	// to evaluate how similar the Token order is.
	userArray := make([]rune, len(userTokens))
	dbArray := make([]rune, len(dbTokens))
	var arrayCounter rune = 0x20

	for index := range userTokens {
		userArray[index] = arrayCounter
		arrayCounter++
	}
	for index := range dbTokens {
		dbArray[index] = arrayCounter
		arrayCounter++
	}

	totalTokens := len(dbTokens) + len(userTokens)

	// Count actions from the database input
	for _, dbToken := range dbTokens {
		if dbToken[0] == '#' {
			totalActions++
		}
	}

	for userIndex, userToken := range userTokens {

		// Count actions from the user input
		if userToken[0] == '#' {
			totalActions++
		}

		for dbIndex, dbToken := range dbTokens {
			if userToken == dbToken {
				userArray[userIndex] = arrayCounter
				dbArray[dbIndex] = arrayCounter
				arrayCounter++

				// Increment word/action matches
				if userToken[0] == '#' {
					matchedActions++
				} else {
					matchedWords++
				}

				// Remove the matched token from the inner array to prevent additional matches
				dbTokens = append(dbTokens[:dbIndex], dbTokens[dbIndex+1:]...)
				break
			}
		}
	}

	totalWords := totalTokens - totalActions

	// The matched ratio must be multiplied by two because each matched word
	// appears twice (once in each input).
	wordMatch := 2 * float64(matchedWords) / float64(totalWords);
	var actionMatch float64
	if totalActions > 0 {
		actionMatch = 2 * float64(matchedActions) / float64(totalActions)
	}

	// Find the difference in word order between the two strings
	userString := string(userArray)
	dbString := string(dbArray)
	dist := textdistance.JaroWinklerDistance(userString, dbString)

	score := (wordMatch * wordWeight + actionMatch * actionWeight + dist * distWeight) /
		(wordWeight + actionWeight + distWeight)

	return score
}
