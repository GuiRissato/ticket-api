package domain

import (
	"errors"
	"fmt"
)

type spotService struct{}

func NewSpotService() *spotService {
	return &spotService{}
}

var (
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
)

func (s *spotService) GenerateSpots(e *Event, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	for i := range quantity {
		spotName := fmt.Sprintf("%c%d", 'A'+i/10, i%10+1)
		spot, err := NewSpot(e, spotName)
		if err != nil {
			return err
		}
		e.Spots = append(e.Spots, *spot)
	}
	return nil
}
