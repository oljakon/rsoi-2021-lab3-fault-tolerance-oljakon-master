package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"rsoi3/src/gateway-service/internal/models"
)

func GetUserRentalsRequest(rentalServiceAddress, username string) (*[]models.Rental, error) {
	requestURL := fmt.Sprintf(rentalServiceAddress + "/api/v1/rental")

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("failed to create an http request")
		return nil, err
	}

	req.Header.Set("X-User-Name", username)

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	rentals := &[]models.Rental{}
	err = json.NewDecoder(res.Body).Decode(rentals)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return rentals, nil
}

func GetUserRentalRequest(rentalServiceAddress, username, rentalUid string) (*models.Rental, error) {
	requestURL := fmt.Sprintf(rentalServiceAddress+"/api/v1/rental/%s", rentalUid)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("failed to create an http request")
		return nil, err
	}

	req.Header.Set("X-User-Name", username)

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	log.Println(res.Body)
	rental := &models.Rental{}
	err = json.NewDecoder(res.Body).Decode(rental)
	log.Println(rental)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return rental, nil
}

func CreateRental(rentalServiceAddress, carUid, dateFrom, dateTo, username, paymentUid string) (string, error) {
	requestURL := fmt.Sprintf(rentalServiceAddress + "/api/v1/rental")

	uid := uuid.New().String()

	rental := &models.Rental{
		RentalUID:  uid,
		Username:   username,
		PaymentUID: paymentUid,
		CarUID:     carUid,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Status:     "IN_PROGRESS",
	}

	data, err := json.Marshal(rental)
	if err != nil {
		return "", fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	return uid, nil
}

func EndRental(rentalServiceAddress, rentalUid string) (string, error) {
	requestURL := fmt.Sprintf(rentalServiceAddress+"/api/v1/rental/%s/finish", rentalUid)

	req, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	var carUid string
	err = json.NewDecoder(res.Body).Decode(&carUid)
	if err != nil {
		return "", fmt.Errorf("failed to encode response: %w", err)
	}

	return carUid, nil
}

func CancelRental(rentalServiceAddress, rentalUid string) (*models.CancelRentalResponse, error) {
	requestURL := fmt.Sprintf(rentalServiceAddress+"/api/v1/rental/%s", rentalUid)

	req, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	canceledRental := &models.CancelRentalResponse{}
	err = json.NewDecoder(res.Body).Decode(canceledRental)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return canceledRental, nil
}
