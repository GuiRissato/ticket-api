package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidSpotNumber   = errors.New("invalid spot number")
	ErrSpotNotFound        = errors.New("spot not found")
	ErrSpotAlreadyReserved = errors.New("spot already reserved")
	ErrSpotNameRequired    = errors.New("spot name is required")
	ErrSpotNameLength      = errors.New("spot name must be at least 2 characters long")
	ErrStartWithLetter     = errors.New("spot name must start with a letter")
	ErrEndWithNumber       = errors.New("spot name must end with a number")
)

type SpotStatus string

const (
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusSold      SpotStatus = "sold"
)

type Spot struct {
	ID       string
	EventId  string
	Name     string
	Status   SpotStatus
	TicketID string
}

func NewSpot(e *Event, name string) (*Spot, error) {
	spot := &Spot{
		ID:      uuid.New().String(),
		EventId: e.ID,
		Name:    name,
		Status:  SpotStatusAvailable,
	}

	if err := spot.Validade(); err != nil {
		return nil, err
	}
	// mesma coisa do que o de cima
	// v := spot.Validade()
	// if v != nil {
	// 	return nil, v
	// }
	return spot, nil
}

func (s *Spot) Validade() error {
	if len(s.Name) == 0 {
		return ErrSpotNameRequired
	}
	if len(s.Name) < 2 {
		return ErrSpotNameLength
	}
	if s.Name[0] < 'A' || s.Name[0] > 'Z' {
		return ErrStartWithLetter
	}
	if s.Name[1] < '0' || s.Name[1] > '9' {
		return ErrEndWithNumber
	}
	return nil
}

func (s *Spot) Reserve(ticketID string) error {
	if s.Status == SpotStatusSold {
		return ErrSpotAlreadyReserved
	}
	s.Status = SpotStatusSold
	s.TicketID = ticketID
	return nil
}
