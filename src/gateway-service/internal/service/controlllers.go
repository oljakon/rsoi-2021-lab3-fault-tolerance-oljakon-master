package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"rsoi3/src/gateway-service/internal/models"
)

const timeFormat = "2006-01-02"

var NeedsRetry = 0

func UsersRentalWithPaymentController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username string) ([]*models.RentalInfo, error) {
	rentals, err := GetUserRentalsRequest(rentalServiceAddress, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get users rentals: %w\n", err)
	}

	rentalsInfo := make([]*models.RentalInfo, 0)

	for _, rental := range *rentals {
		paymentUID := rental.PaymentUID
		payment, err := GetPayment(paymentServiceAddress, paymentUID)
		if err != nil {
			fmt.Printf("failed to get payment: %v", err)
			return nil, err
		}

		carUID := rental.CarUID
		car, err := GetCar(carServiceAddress, carUID)
		if err != nil {
			fmt.Printf("failed to get car: %v", err)
			return nil, err
		}

		parsedDateFrom := rental.DateFrom[:10]
		parsedDateTo := rental.DateTo[:10]

		rentalInfo := &models.RentalInfo{
			RentalUID: rental.RentalUID,
			CarUID:    rental.CarUID,
			DateFrom:  parsedDateFrom,
			DateTo:    parsedDateTo,
			Status:    rental.Status,
			Car: &models.CarField{
				CarUID:             car.CarUID,
				Brand:              car.Brand,
				Model:              car.Model,
				RegistrationNumber: car.RegistrationNumber,
			},
			Payment: &models.PaymentField{
				PaymentUID: payment.PaymentUID,
				Status:     payment.Status,
				Price:      payment.Price,
			},
		}

		rentalsInfo = append(rentalsInfo, rentalInfo)
	}

	return rentalsInfo, nil
}

func UsersRentalFullInfoController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username, rentalUID string) (*models.RentalInfo, error) {
	rental, err := GetUserRentalRequest(rentalServiceAddress, username, rentalUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rental: %w\n", err)
	}

	paymentUID := rental.PaymentUID
	payment, err := GetPayment(paymentServiceAddress, paymentUID)
	if err != nil {
		log.Println("payment service unavailable: ", err)
	}

	carUID := rental.CarUID
	car, err := GetCar(carServiceAddress, carUID)
	if err != nil {
		fmt.Printf("failed to get car: %v", err)
		return nil, err
	}

	parsedDateFrom := rental.DateFrom[:10]
	parsedDateTo := rental.DateTo[:10]

	rentalInfo := &models.RentalInfo{
		RentalUID: rental.RentalUID,
		CarUID:    rental.CarUID,
		DateFrom:  parsedDateFrom,
		DateTo:    parsedDateTo,
		Status:    rental.Status,
		Car: &models.CarField{
			CarUID:             car.CarUID,
			Brand:              car.Brand,
			Model:              car.Model,
			RegistrationNumber: car.RegistrationNumber,
		},
		Payment: &models.PaymentField{
			PaymentUID: payment.PaymentUID,
			Status:     payment.Status,
			Price:      payment.Price,
		},
	}

	return rentalInfo, nil
}

func RentCarController(rentalServiceAddress, paymentServiceAddress, carServiceAddress, username string, rentInfo *models.RentCarRequest) (*models.CreateRentalResponse, error) {
	price, err := ReserveCar(carServiceAddress, rentInfo.CarUID)
	if err != nil {
		return nil, fmt.Errorf("failed to rent car by uid: %w\n", err)
	}

	t1, _ := time.Parse(timeFormat, rentInfo.DateFrom)
	t2, _ := time.Parse(timeFormat, rentInfo.DateTo)
	diff := t2.Sub(t1)
	days := int(diff.Hours() / 24)
	fullPrice := price * days

	paymentUid, err := CreatePayment(paymentServiceAddress, fullPrice)
	if err != nil {
		if errors.Is(err, ErrPayment503) {
			return nil, ErrPayment503
		}
		err = EndCarReserve(carServiceAddress, rentInfo.CarUID)
		if err != nil {
			return nil, fmt.Errorf("failed to end rental: %w\n", err)
		}
		return nil, fmt.Errorf("failed to rent car by uid: failed to create payment %w\n", err)
	}

	rentalUid, err := CreateRental(rentalServiceAddress, rentInfo.CarUID, rentInfo.DateFrom, rentInfo.DateTo, username, paymentUid)
	if err != nil {
		err = EndCarReserve(carServiceAddress, rentInfo.CarUID)
		if err != nil {
			return nil, fmt.Errorf("failed to end rental: %w\n", err)
		}
		return nil, fmt.Errorf("failed to rent car by uid: : failed to create rental %w\n", err)
	}

	createRentalResponse := &models.CreateRentalResponse{
		RentalUID: rentalUid,
		CarUID:    rentInfo.CarUID,
		DateFrom:  rentInfo.DateFrom,
		DateTo:    rentInfo.DateTo,
		Status:    "IN_PROGRESS",
		Payment: &models.PaymentField{
			PaymentUID: paymentUid,
			Status:     "PAID",
			Price:      fullPrice,
		},
	}

	return createRentalResponse, nil
}

func EndRentalController(rentalServiceAddress, carServiceAddress, rentalUid string) error {
	carUid, err := EndRental(rentalServiceAddress, rentalUid)
	if err != nil {
		return fmt.Errorf("failed to end rental: %w\n", err)
	}

	err = EndCarReserve(carServiceAddress, carUid)
	if err != nil {
		return fmt.Errorf("failed to end rental: %w\n", err)
	}

	return nil
}

func CancelRentalController(rentalServiceAddress, carServiceAddress, paymentServiceAddress, rentalUid string) error {
	canceledRental, err := CancelRental(rentalServiceAddress, rentalUid)
	if err != nil {
		return fmt.Errorf("failed to cancel rental: %w\n", err)
	}

	err = EndCarReserve(carServiceAddress, canceledRental.CarUID)
	if err != nil {
		return fmt.Errorf("failed to cancel car reserve: %w\n", err)
	}

	NeedsRetry = 1
	err = CancelPayment(paymentServiceAddress, canceledRental.PaymentUID)
	if err != nil {
		fmt.Println("cancel payment: payment service unavailable")
		go RetryPaymentCancel(paymentServiceAddress, canceledRental.PaymentUID)
		return fmt.Errorf("failed to cancel payment: %w\n", err)
	}

	return nil
}

func RetryPaymentCancel(paymentServiceAddress, PaymentUID string) {
	i := 0
	retries := 100

	for i < retries {
		err := CancelPayment(paymentServiceAddress, PaymentUID)
		if err == nil {
			return
		}
		time.Sleep(10 * time.Second)
		i += 1
	}
}
