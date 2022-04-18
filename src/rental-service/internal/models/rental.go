package models

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type Rental struct {
	ID         int    `json:"id"`
	RentalUID  string `json:"rentalUid"`
	Username   string `json:"username"`
	PaymentUID string `json:"paymentUid"`
	CarUID     string `json:"car_uid"`
	DateFrom   string `json:"dateFrom"`
	DateTo     string `json:"dateTo"`
	Status     string `json:"status"`
}

type CancelRentalResponse struct {
	CarUID     string `json:"carUid"`
	PaymentUID string `json:"paymentUid"`
}

func (r *Rental) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.RentalUID, validation.Required, is.UUID),
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.PaymentUID, validation.Required, is.UUID),
		validation.Field(&r.CarUID, validation.Required, is.UUID),
		validation.Field(&r.DateFrom, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&r.DateTo, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&r.Status, validation.Required, validation.In("IN_PROGRESS", "FINISHED", "CANCELED")),
	)
}

func TestRental(t *testing.T) *Rental {
	t.Helper()

	rentalUid := uuid.New().String()
	paymentUid := uuid.New().String()
	carUid := uuid.New().String()

	return &Rental{
		RentalUID:  rentalUid,
		Username:   "test",
		PaymentUID: paymentUid,
		CarUID:     carUid,
		DateFrom:   "2021-10-08",
		DateTo:     "2021-10-11",
		Status:     "IN_PROGRESS",
	}
}
