package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"rsoi3/src/payment-service/internal/db"
	"rsoi3/src/payment-service/internal/models"
)

type Repository interface {
	GetUserRentalsRepo(all bool) ([]*models.Payment, error)
}

type PaymentRepository struct {
	db *sql.DB
}

const selectPaymentByUID = `SELECT payment_uid, status, price FROM payment where payment_uid = $1;`

func (r *PaymentRepository) GetPaymentByUidRepo(paymentUID string) (*models.Payment, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	var payment models.Payment

	row := r.db.QueryRow(selectPaymentByUID, paymentUID)
	err := row.Scan(&payment.PaymentUID, &payment.Status, &payment.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &payment, err
		}
	}

	return &payment, nil
}

const createPayment = `INSERT INTO payment (payment_uid, status, price) VALUES ($1, $2, $3) RETURNING payment_uid;`

func (r *PaymentRepository) CreatePaymentRepo(payment *models.Payment) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Query(createPayment, payment.PaymentUID, payment.Status, payment.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	return nil
}

const cancelPaymentReserve = `UPDATE payment SET status='CANCELED' WHERE payment_uid=$1;`

func (r *PaymentRepository) CancelPaymentRepo(paymentUid string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Exec(cancelPaymentReserve, paymentUid)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %w", err)
	}

	return nil
}
