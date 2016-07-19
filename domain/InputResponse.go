package domain

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	textdistance "github.com/masatana/go-textdistance"
)

// InputOutputStore ...
type InputOutputStore interface {
	Store(InputOutput) error
	Response(string, string) (InputOutput, error)
	Upvote(InputOutput) error
	Downvote(InputOutput) error
}

// InputOutput ...
type InputOutput struct {
	int    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Input  string `json:"pc_input" gorm:"column:pc_input"`
	Output string `json:"gm_response" gorm:"column:gm_response"`
	RoomID int64  `json:"room_id"`
}

// TableName setting this function satisfies the gorm interface and changes the table
// name.
func (i InputOutput) TableName() string {
	return "input_response_pairs"
}

// Match ...
type Match struct {
	ID        int    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	UID       string `json:"uid" sql:"unique"`
	Input     string `json:"input"`
	Match     string `json:"match"`
	Room      string `json:"room"`
	UpVotes   int    `json:"upvotes"`
	DownVotes int    `json:"downvotes"`
}

// NewMatch ...
func NewMatch(input, match, room string) Match {
	return Match{
		Input: input,
		Match: match,
		UID:   Hash(input, match),
		Room:  room,
	}
}

// NewInputOutput ...
func NewInputOutput(input, output string, room int64) InputOutput {
	return InputOutput{
		Input:  input,
		Output: output,
		RoomID: room,
	}
}

// SoftMatch is the actual algorithm that is used to match two inputs.
// it takes the user's input string and a slice of the InputOutput pairs
// it should match against.
func SoftMatch(input string, pairs []InputOutput) (InputOutput, float64) {
	response := InputOutput{}

	var maxScore float64
	for _, pair := range pairs {
		indb := pair.Input
		score := textdistance.JaroWinklerDistance(input, indb)
		if score > maxScore {
			maxScore = score
			response = pair
		}
	}
	return response, maxScore
}

// Hash ...
func Hash(parts ...string) string {
	var buffer bytes.Buffer
	for _, part := range parts {
		buffer.WriteString(part)
	}
	return fmt.Sprintf("%x", sha1.Sum(buffer.Bytes()))
}
