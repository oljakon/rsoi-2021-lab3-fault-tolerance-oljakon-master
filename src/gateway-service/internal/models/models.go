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

type CarsLimited struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"pageSize"`
	TotalElements int    `json:"totalElements"`
	Items         *[]Car `json:"items"`
}

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

type Payment struct {
	ID         int    `json:"id"`
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}

type PaymentField struct {
	PaymentUID string `json:"paymentUid"`
	Status     string `json:"status"`
	Price      int    `json:"price"`
}

type CarField struct {
	CarUID             string `json:"carUid"`
	Brand              string `json:"brand"`
	Model              string `json:"model"`
	RegistrationNumber string `json:"registrationNumber"`
}

type RentalInfo struct {
	RentalUID string        `json:"rentalUid"`
	CarUID    string        `json:"carUid"`
	DateFrom  string        `json:"dateFrom"`
	DateTo    string        `json:"dateTo"`
	Status    string        `json:"status"`
	Car       *CarField     `json:"car"`
	Payment   *PaymentField `json:"payment"`
}

type RentalInfoEmptyPayment struct {
	RentalUID string    `json:"rentalUid"`
	CarUID    string    `json:"carUid"`
	DateFrom  string    `json:"dateFrom"`
	DateTo    string    `json:"dateTo"`
	Status    string    `json:"status"`
	Car       *CarField `json:"car"`
	Payment   string    `json:"payment"`
}

type CreateRentalResponse struct {
	RentalUID string        `json:"rentalUid"`
	CarUID    string        `json:"carUid"`
	DateFrom  string        `json:"dateFrom"`
	DateTo    string        `json:"dateTo"`
	Status    string        `json:"status"`
	Payment   *PaymentField `json:"payment"`
}

type RentCarRequest struct {
	CarUID   string `json:"carUid"`
	DateFrom string `json:"dateFrom"`
	DateTo   string `json:"dateTo"`
}

type CancelRentalResponse struct {
	CarUID     string `json:"carUid"`
	PaymentUID string `json:"paymentUid"`
}

func (ri *RentalInfo) Validate() error {
	return validation.ValidateStruct(
		ri,
		validation.Field(&ri.RentalUID, validation.Required, is.UUID),
		validation.Field(&ri.CarUID, validation.Required, is.UUID),
		validation.Field(&ri.DateFrom, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&ri.DateTo, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&ri.Status, validation.Required, validation.In("IN_PROGRESS", "FINISHED", "CANCELED")),
		validation.Field(&ri.Car, validation.Required),
		validation.Field(&ri.Payment, validation.Required),
	)
}

func TestRentalInfo(t *testing.T) *RentalInfo {
	t.Helper()

	rentalUid := uuid.New().String()
	carUid := uuid.New().String()

	return &RentalInfo{
		RentalUID: rentalUid,
		CarUID:    carUid,
		DateFrom:  "2021-10-08",
		DateTo:    "2021-10-11",
		Status:    "IN_PROGRESS",
		Car: &CarField{
			CarUID:             "test",
			Brand:              "test",
			Model:              "test",
			RegistrationNumber: "test",
		},
		Payment: &PaymentField{
			PaymentUID: "test",
			Status:     "PAID",
			Price:      1000,
		},
	}
}
