package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRentalInfo_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		u       func() *Car
		isValid bool
	}{
		{
			name: "valid",
			u: func() *Car {
				return TestCar(t)
			},
			isValid: true,
		},
		{
			name: "invalid: uid not uuid",
			u: func() *Car {
				u := TestCar(t)
				u.CarUID = "12345"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: unsupported type",
			u: func() *Car {
				u := TestCar(t)
				u.Type = "CABRIOLET"
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: empty registration number",
			u: func() *Car {
				u := TestCar(t)
				u.RegistrationNumber = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid: zero price",
			u: func() *Car {
				u := TestCar(t)
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
