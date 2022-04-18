package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRentalInfo_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *RentalInfo
		isValid bool
	}{
		{
			name: "valid",
			u: func() *RentalInfo {
				return TestRentalInfo(t)
			},
			isValid: true,
		},
		{
			name: "invalid: uid not uuid",
			u: func() *RentalInfo {
				u := TestRentalInfo(t)
				u.RentalUID = "12345"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: empty car field",
			u: func() *RentalInfo {
				u := TestRentalInfo(t)
				u.Car = nil
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: wrong date format",
			u: func() *RentalInfo {
				u := TestRentalInfo(t)
				u.DateFrom = "01.01.2021"
				u.DateFrom = "02.01.2021"
				return u
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}

func TestRentalInfo_PriceCounting(t *testing.T) {
	timeFormat := "2006-01-02"

	rentalInfo := &RentalInfo{
		RentalUID: "test",
		CarUID:    "test",
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

	t.Run("count price", func(t *testing.T) {
		t1, _ := time.Parse(timeFormat, rentalInfo.DateFrom)
		t2, _ := time.Parse(timeFormat, rentalInfo.DateTo)
		diff := t2.Sub(t1)
		days := int(diff.Hours() / 24)
		fullPrice := rentalInfo.Payment.Price * days

		assert.Equal(t, 3000, fullPrice)
	})
}
