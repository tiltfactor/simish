package domain

import (
	"math"
)

// voteAsymptoteSlope is the slope of the total-vote function. It determines how
// quickly our confidence increases as we get more votes for a pair.
const voteAsymptoteSlope float64 = 0.2

// initialVoteScore is the score given to a pair with no votes. As the pair is
// up-voted or down-voted, it increases or decreases from that point.
const initialVoteScore float64 = 0.7

// InputOutput used to map the database structure to the input output pair used by the
// program.
type InputOutput struct {
	int        `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	AiCol      int64  `gorm:"column:ai_col"`
	ResultType int64  `gorm:"column:result_type"`
	Input      string `json:"pc_input" gorm:"column:pc_input"`
	Output     string `json:"gm_response" gorm:"column:gm_response"`
	RoomID     int64  `json:"room_id"`
	Upvotes	   float64  `gorm:"column:pc_upvotes_weighted"`
	Downvotes  float64  `gorm:"column:pc_downvotes_weighted"`
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
	userTokens := prepareInput(input)

	for _, pair := range pairs {
		dbTokens := prepareInput(pair.Input)
		rawScore := getRawScore(userTokens, dbTokens)
		voteScore := getVoteScore(pair.Upvotes, pair.Downvotes)
		score := rawScore * voteScore
		if score > maxScore {
			maxScore = score
			response = pair
		}
	}
	return response, maxScore
}

func getVoteScore(upvotes float64, downvotes float64) float64 {
	totalVotes := upvotes + downvotes
	upvoteRatio := upvotes / totalVotes
	if math.IsNaN(upvoteRatio) {
		upvoteRatio = 0
	}

	// totalVoteAsymptote is a value from 0 to 1. When there are no votes, totalAsymptote = 0.
	// As the number of votes increases, totalAsymptote increases towards 1
	totalVoteAsymptote := 1 - 1 / (voteAsymptoteSlope * totalVotes + 1)

	// voteScore is a value from 0 to 1. When there are no votes, voteScore = 0.5.
	// As we get more votes, the voteScore becomes more extreme, increasing towards
	// 0 if there are more down-votes, or 1 if there are more up-votes
	voteScore := initialVoteScore + totalVoteAsymptote * (upvoteRatio - initialVoteScore);
	return voteScore
}
