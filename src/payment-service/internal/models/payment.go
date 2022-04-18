package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"testing"
)

type Payment struct {
	ID         int    `json:"id"`
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}

func (p *Payment) Validate() error {
	return validation.ValidateStruct(
		p,
		validation.Field(&p.PaymentUID, validation.Required, is.UUID),
		validation.Field(&p.Status, validation.Required, validation.In("PAID", "CANCELED")),
		validation.Field(&p.Price, validation.Required, validation.Min(1)),
	)
}

func TestPayment(t *testing.T) *Payment {
	t.Helper()

	uid := uuid.New().String()

	return &Payment{
		PaymentUID: uid,
		Status:     "PAID",
		Price:      1000,
	}
}
