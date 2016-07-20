package domain

import textdistance "github.com/masatana/go-textdistance"

// InputOutput used to map the database structure to the input output pair used by the
// program.
type InputOutput struct {
	int    `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Input  string `json:"pc_input" gorm:"column:pc_input"`
	Output string `json:"gm_response" gorm:"column:gm_response"`
	RoomID int64  `json:"room_id"`
}

// InputOutputStore is the interface that needs to be fulfilled by other store implementations.
type InputOutputStore interface {
	SaveInputOutput(InputOutput) error
	Response(string, int64) (InputOutput, float64)
}

// TableName setting this function satisfies the gorm interface and changes the table
// name.
func (i InputOutput) TableName() string {
	return "input_response_pairs"
}

// NewInputOutput returns a new input output object.
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
