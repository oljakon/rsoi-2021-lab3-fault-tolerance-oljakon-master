package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"rsoi3/src/gateway-service/internal/models"
)

var ErrPayment503 = errors.New("payment service unavailable")

func GetPayment(paymentServiceAddress string, paymentUID string) (*models.Payment, error) {
	requestURL := fmt.Sprintf(paymentServiceAddress+"/api/v1/payment/%s", paymentUID)

	var errorPayment = &models.Payment{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Println("failed to create an http request")
		return errorPayment, err
	}

	client := &http.Client{
		Timeout: time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return errorPayment, fmt.Errorf("failed request to payment service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	payment := &models.Payment{}
	err = json.NewDecoder(res.Body).Decode(payment)
	if err != nil {
		return errorPayment, fmt.Errorf("failed to decode response: %w", err)
	}

	return payment, nil
}

func CreatePayment(paymentServiceAddress string, price int) (string, error) {
	requestURL := fmt.Sprintf(paymentServiceAddress + "/api/v1/payment")

	uid := uuid.New().String()

	payment := &models.Payment{
		PaymentUID: uid,
		Status:     "PAID",
		Price:      price,
	}

	data, err := json.Marshal(payment)
	if err != nil {
		return "", fmt.Errorf("encoding error: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: time.Second,
	}
	res, err := client.Do(req)
	if res.StatusCode == http.StatusServiceUnavailable {
		return "", ErrPayment503
	}
	if err != nil {
		return "", fmt.Errorf("failed request to payment service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	return uid, nil
}

func CancelPayment(paymentServiceAddress, paymentUid string) error {
	requestURL := fmt.Sprintf(paymentServiceAddress+"/api/v1/payment/%s", paymentUid)

	req, err := http.NewRequest(http.MethodPatch, requestURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create an http request: %w", err)
	}

	client := &http.Client{
		Timeout: time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed request to rental service: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(res.Body)

	return nil
}
