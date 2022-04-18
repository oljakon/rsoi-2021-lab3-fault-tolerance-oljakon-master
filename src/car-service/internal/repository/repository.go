package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"rsoi3/src/car-service/internal/db"
	"rsoi3/src/car-service/internal/models"
)

type Repository interface {
	GetAvailableCarsRepo(all bool) ([]*models.Car, error)
	GetCarByUidRepo(carUID string) (*models.Car, error)
	ReserveCarRepo(carUID string) (*models.Car, error)
}

type CarRepository struct {
	db *sql.DB
}

const (
	selectAllCars       = `SELECT id, car_uid, brand, model, registration_number, power, price, type, availability FROM cars;`
	selectAvailableCars = `SELECT id, car_uid, brand, model, registration_number, power, price, type, availability FROM cars WHERE availability IS TRUE;`
)

func (r *CarRepository) GetAvailableCarsRepo(all bool) ([]*models.Car, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	var cars []*models.Car

	if all {
		rows, err := r.db.Query(selectAllCars)
		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}

		for rows.Next() {
			car := new(models.Car)
			if err := rows.Scan(&car.ID, &car.CarUID, &car.Brand, &car.Model, &car.RegistrationNumber, &car.Power, &car.Price, &car.Type, &car.Availability); err != nil {
				return nil, fmt.Errorf("failed to execute the query: %w", err)
			}
			cars = append(cars, car)
		}

		defer rows.Close()
	} else {
		rows, err := r.db.Query(selectAvailableCars)
		if err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}

		for rows.Next() {
			car := new(models.Car)
			if err := rows.Scan(&car.ID, &car.CarUID, &car.Brand, &car.Model, &car.RegistrationNumber, &car.Power, &car.Price, &car.Type, &car.Availability); err != nil {
				return nil, fmt.Errorf("failed to execute the query: %w", err)
			}
			cars = append(cars, car)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("failed to execute the query: %w", err)
		}

		defer rows.Close()
	}

	return cars, nil
}

const selectCarByUID = `SELECT car_uid, brand, model, registration_number FROM cars WHERE car_uid = $1;`

func (r *CarRepository) GetCarByUidRepo(carUID string) (*models.Car, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	var car models.Car

	row := r.db.QueryRow(selectCarByUID, carUID)
	err := row.Scan(&car.CarUID, &car.Brand, &car.Model, &car.RegistrationNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &car, err
		}
	}

	return &car, nil
}

const (
	reserveCarByUID     = `UPDATE cars SET availability=false WHERE car_uid=$1;`
	selectCarPriceByUID = `SELECT price FROM cars where car_uid = $1;`
)

func (r *CarRepository) ReserveCarRepo(carUID string) (int, error) {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Exec(reserveCarByUID, carUID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute the query: %w", err)
	}

	var price int
	err = r.db.QueryRow(selectCarPriceByUID, carUID).Scan(&price)
	if err != nil {
		return 0, fmt.Errorf("failed to execute the query: %w", err)
	}

	return price, nil
}

const endCarReserve = `UPDATE cars SET availability=true WHERE car_uid=$1;`

func (r *CarRepository) EndCarReserveRepo(carUID string) error {
	r.db = db.CreateConnection()
	defer r.db.Close()

	_, err := r.db.Exec(endCarReserve, carUID)
	if err != nil {
		return fmt.Errorf("failed to execute the query: %w", err)
	}

	return nil
}
