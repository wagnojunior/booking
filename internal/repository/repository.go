package repository

import (
	"time"

	"github.com/wagnojunior/booking/internal/models"
)

type DatabaseRepo interface {
	AllUders() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error)
}
