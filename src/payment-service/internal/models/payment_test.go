package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRentalInfo_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *Payment
		isValid bool
	}{
		{
			name: "valid",
			u: func() *Payment {
				return TestPayment(t)
			},
			isValid: true,
		},
		{
			name: "invalid: uid not uuid",
			u: func() *Payment {
				u := TestPayment(t)
				u.PaymentUID = "12345"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: unsupported status",
			u: func() *Payment {
				u := TestPayment(t)
				u.Status = "REVERSED"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: empty status",
			u: func() *Payment {
				u := TestPayment(t)
				u.Status = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: zero price",
			u: func() *Payment {
				u := TestPayment(t)
				u.Price = 0
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
