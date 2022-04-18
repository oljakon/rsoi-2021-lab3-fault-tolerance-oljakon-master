package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"rsoi3/src/rental-service/internal/db"
	"rsoi3/src/rental-service/internal/models"
)

type Repository interface {
	GetUserRentalsRepo(all bool) ([]*models.Rental, error)
}

type RentalRepository struct {
	db *sql.DB
}

const selectUserRentals = `SELECT rental_uid, payment_uid, car_uid, date_from, date_to, status FROM rental WHERE username = $1;`

func (r *RentalRepository) GetUserRentalsRepo(username string) ([]*models.Rental, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	var rentals []*models.Rental

	rows, err := r.db.Query(selectUserRentals, username)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	for rows.Next() {
		rental := new(models.Rental)
		rental.Username = username
		if err := rows.Scan(&rental.RentalUID, &rental.PaymentUID, &rental.CarUID, &rental.DateFrom, &rental.DateTo, &rental.Status); err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}
		rentals = append(rentals, rental)
	}

	defer rows.Close()

	return rentals, nil
}

const selectRentalByUid = `SELECT payment_uid, car_uid, date_from, date_to, status FROM rental where rental_uid = $1`

func (r *RentalRepository) GetRentalRepo(rentalUid string) (*models.Rental, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	var rental models.Rental

	row := r.db.QueryRow(selectRentalByUid, rentalUid)
	rental.RentalUID = rentalUid
	err := row.Scan(&rental.PaymentUID, &rental.CarUID, &rental.DateFrom, &rental.DateTo, &rental.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &rental, err
		}
	}

	return &rental, nil
}

const createRental = `INSERT INTO rental (rental_uid, username, payment_uid, car_uid, date_from, date_to, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING rental_uid;`

func (r *RentalRepository) CreateRentalRepo(rental *models.Rental) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Query(createRental, rental.RentalUID, rental.Username, rental.PaymentUID, rental.CarUID, rental.DateFrom, rental.DateTo, rental.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}

	return nil
}

const (
	endRental         = `UPDATE rental SET status='FINISHED' WHERE rental_uid=$1;`
	selectRentalByUID = `SELECT car_uid FROM rental where rental_uid = $1;`
)

func (r *RentalRepository) EndRentalRepo(rentalUid string) (string, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Query(endRental, rentalUid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
	}

	var carUid string
	err = r.db.QueryRow(selectRentalByUID, rentalUid).Scan(&carUid)
	if err != nil {
		return "", fmt.Errorf("failed to execute the query: %w", err)
	}

	return carUid, nil
}

const (
	cancelRental           = `UPDATE rental SET status='CANCELED' WHERE rental_uid=$1;`
	selectRentalCredsByUID = `SELECT car_uid, payment_uid FROM rental where rental_uid = $1;`
)

func (r *RentalRepository) CancelRentalRepo(rentalUid string) (*models.CancelRentalResponse, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Query(cancelRental, rentalUid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	canceled := new(models.CancelRentalResponse)
	err = r.db.QueryRow(selectRentalCredsByUID, rentalUid).Scan(&canceled.CarUID, &canceled.PaymentUID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %w", err)
	}

	return canceled, nil
}
