package models

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type Car struct {
	ID                 int    `json:"id"`
	CarUID             string `json:"carUid"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	RegistrationNumber string `json:"registrationNumber"`
	Power              int    `json:"power"`
	Price              int    `json:"price"`
	Type               string `json:"type"`
	Availability       bool   `json:"available"`
}

func (c *Car) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.CarUID, validation.Required, is.UUID),
		validation.Field(&c.Brand, validation.Required),
		validation.Field(&c.Model, validation.Required),
		validation.Field(&c.RegistrationNumber, validation.Required),
		validation.Field(&c.Power, validation.Required, validation.Min(1)),
		validation.Field(&c.Price, validation.Required, validation.Min(1)),
		validation.Field(&c.Type, validation.Required, validation.In("SEDAN", "SUV", "MINIVAN", "ROADSTER")),
		validation.Field(&c.Availability, validation.Required, validation.In(true, false)),
	)
}

func TestCar(t *testing.T) *Car {
	t.Helper()

	uid := uuid.New().String()

	return &Car{
		CarUID:             uid,
		Brand:              "test",
		Model:              "test",
		RegistrationNumber: "test",
		Power:              1000,
		Price:              1000,
		Type:               "SEDAN",
		Availability:       true,
	}
}
