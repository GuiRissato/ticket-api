package domain

type EventRepository interface {
	ListEvents() ([]Event, error)
	FindEventByID(eventID string) (*Event, error)
	FindSpotsByEventID(eventID string) ([]*Spot, error)
	FindSpotByName(eventID, spotName string) (*Spot, error)
	CreateEvent(event *Event) error    //pegar no projeto
	CreateSpot(spot *Spot) error       //pegar no projeto
	CreateTicket(ticket *Ticket) error //pegar no projeto
	ReserveSpot(spotID, ticketID string) error
}
