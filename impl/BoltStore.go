package impl

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	"github.com/masatana/go-textdistance"
	"github.com/tiltfactor/simish/domain"
	// Only need SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// SQLStore ..
type SQLStore struct {
	db *gorm.DB
}

// NewSQLStore ...
func NewSQLStore(path string) (*SQLStore, error) {
	db, err := gorm.Open("mysql", path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := db.DB().Ping(); err != nil {
		log.Println(err)
		return nil, err
	}
	if !db.HasTable(new(domain.InputOutput)) {
		db.CreateTable(new(domain.Match))
		db.CreateTable(new(domain.InputOutput))
	}

	return &SQLStore{
		db: db,
	}, nil
}

// SaveMatch ...
func (s SQLStore) SaveMatch(match domain.Match) {
	s.db.Save(&match)
}

// SaveInputOutput ...
func (s SQLStore) SaveInputOutput(io domain.InputOutput) {
	indb := domain.InputOutput{}
	s.db.Model(new(domain.InputOutput)).
		Where("pc_input = ? AND gm_input = ?", io.Input, io.Output).
		First(&indb)

	// If the input pair has not been saved.
	if indb.Input == "" {
		s.db.Save(&io)
	}
}

// Store ..
func (s SQLStore) Store(io domain.InputOutput) error {
	s.db.Save(&io)
	return nil
}

// Response ..
func (s SQLStore) Response(in, room string) (domain.InputOutput, domain.Match, float64, error) {
	pairs := []domain.InputOutput{}
	s.db.Model(new(domain.InputOutput)).Where("room_id = ?", room).Find(&pairs)
	response := domain.InputOutput{}
	var maxScore float64
	match := domain.Match{}
	for _, pair := range pairs {
		indb := pair.Input
		score := textdistance.JaroWinklerDistance(in, indb)

		// dm := domain.Match{}
		// s.db.Model(new(domain.Match)).
		// 	Where("uid = ?", domain.Hash(in, indb)).
		// 	First(&dm)
		//
		// // We need to convert them to floats so they don't get truncated
		// votes := float64(dm.UpVotes) / float64((dm.UpVotes + dm.DownVotes))
		//
		// if (dm.UpVotes + dm.DownVotes) > 0 {
		// 	score *= votes
		// }
		if score > maxScore {
			maxScore = score
			response = pair
			// match = dm
		}
	}
	return response, match, maxScore, nil
}

func (s SQLStore) containsMatch(pair *domain.Match) bool {
	indb := domain.Match{}
	s.db.Model(new(domain.Match)).
		Where("uid = ?", domain.Hash(pair.Input, pair.Match)).
		First(&indb)
	return indb.UID != ""
}

// Upvote ..
func (s SQLStore) Upvote(in, match, room string) error {
	pair := domain.Match{}

	s.db.Model(new(domain.Match)).
		Where("uid = ?", domain.Hash(in, match)).
		Find(&pair)

	fmt.Println(pair.UID)
	if pair.UID != "" {
		pair.UpVotes++
		s.db.Save(&pair)
	} else {
		pair.Room = room
		pair.Input = in
		pair.Match = match
		pair.UID = domain.Hash(in, match)
		pair.UpVotes = 1
		s.db.Save(&pair)
	}
	return nil
}

// Downvote ...
func (s SQLStore) Downvote(in, match, room string) error {
	pair := domain.Match{}

	s.db.Model(new(domain.Match)).
		Where("uid = ?", domain.Hash(in, match)).
		Find(&pair)

	fmt.Println(pair.UID)
	if pair.UID != "" {
		pair.DownVotes++
		s.db.Save(&pair)
	} else {
		pair.Room = room
		pair.Input = in
		pair.Match = match
		pair.UID = domain.Hash(in, match)
		pair.DownVotes = 1
		s.db.Save(&pair)
	}
	return nil
}
