package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRentalInfo_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *Rental
		isValid bool
	}{
		{
			name: "valid",
			u: func() *Rental {
				return TestRental(t)
			},
			isValid: true,
		},
		{
			name: "invalid: uid not uuid",
			u: func() *Rental {
				u := TestRental(t)
				u.RentalUID = "12345"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: unsupported status",
			u: func() *Rental {
				u := TestRental(t)
				u.Status = "STATED"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: invalid date format",
			u: func() *Rental {
				u := TestRental(t)
				u.DateFrom = "01.01.2021"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: empty username",
			u: func() *Rental {
				u := TestRental(t)
				u.Username = ""
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
