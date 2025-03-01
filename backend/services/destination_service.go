package services

import (
	"math/rand"
	"time"

	"github.com/shubhsherl/globetrotter/backend/db"
	"github.com/shubhsherl/globetrotter/backend/models"
)

// DestinationService handles destination-related operations
type DestinationService struct {
	db *db.Database
}

// NewDestinationService creates a new destination service
func NewDestinationService(database *db.Database) *DestinationService {
	return &DestinationService{
		db: database,
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetRandomDestination returns a random destination with multiple choice options
func (s *DestinationService) GetRandomDestination() (models.Destination, error) {
	// Implementation depends on your database structure
	// This is a placeholder implementation
	destinations, err := s.db.GetAllDestinations()
	if err != nil {
		return models.Destination{}, err
	}
	
	if len(destinations) == 0 {
		return models.Destination{}, nil
	}
	
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(destinations))
	return destinations[randomIndex], nil
}
