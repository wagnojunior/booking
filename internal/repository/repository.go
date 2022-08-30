package repository

import "github.com/wagnojunior/booking/internal/models"

type DatabaseRepo interface {
	AllUders() bool

	InsertReservation(res models.Reservation) error
}
