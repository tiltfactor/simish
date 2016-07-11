package domain

import (
	"bytes"
	"crypto/sha1"
	"fmt"
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
	Input  string `json:"pc_input"`
	Output string `json:"gm_response"`
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

// Hash ...
func Hash(parts ...string) string {
	var buffer bytes.Buffer
	for _, part := range parts {
		buffer.WriteString(part)
	}
	return fmt.Sprintf("%x", sha1.Sum(buffer.Bytes()))
}
