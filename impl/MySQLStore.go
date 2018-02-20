package impl

import (
	"github.com/jinzhu/gorm"
	"github.com/tiltfactor/simish/domain"
	// Only need SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// SQLStore ..
type SQLStore struct {
	db *gorm.DB
}

// NewSQLStore returns an initialized sql store. It takes the db path as its only arugment and
// returns an error if it cannot connect to the database.
func NewSQLStore(path string) (*SQLStore, error) {
	db, err := gorm.Open("mysql", path)
	if err != nil {
		return nil, err
	}

	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	if !db.HasTable(new(domain.InputOutput)) {
		db.CreateTable(new(domain.InputOutput))
	}

	return &SQLStore{
		db: db,
	}, nil
}

// SaveInputOutput saves the provided input output pair
func (s SQLStore) SaveInputOutput(io domain.InputOutput) error {
	indb := domain.InputOutput{}
	s.db.Model(new(domain.InputOutput)).
		Where("pc_input = ? AND gm_input = ?", io.Input, io.Output).
		First(&indb)

	// If the input pair has not been saved.
	if indb.Input == "" {
		s.db.Save(&io)
	}
	return nil
}

// Response gets the available pairs from the database and runs the SoftMatch algorithm
// returning the found pair and the score of the pair.
func (s SQLStore) Response(in string, room int64) (domain.InputOutput, float64) {
	pairs := []domain.InputOutput{}
	s.db.Model(&domain.InputOutput{}).Where("room_id = ? AND NOT disabled", room).Find(&pairs)
	return domain.SoftMatch(in, pairs)
}

// GetAllPairs is used for testing to retrieve the pairs for a given room
func (s SQLStore) GetAllPairs(room int64) []domain.InputOutput {
	pairs := []domain.InputOutput{}
	s.db.Model(&domain.InputOutput{}).Where("room_id = ? AND NOT disabled", room).Find(&pairs)
	return pairs
}
