package ride

import (
	"database/sql"

	"github.com/gfmanica/splitz-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetRides() ([]types.Ride, error) {
	rows, err := s.db.Query("SELECT id_ride, ds_ride, vl_ride,  dt_init, dt_finish, fg_count_weekend FROM ride")

	if err != nil {
		return nil, err
	}

	rides := make([]types.Ride, 0)

	for rows.Next() {
		ride, err := scanRowIntoRide(rows)

		if err != nil {
			return nil, err
		}

		rides = append(rides, *ride)
	}

	return rides, nil
}

func (s *Store) GetRideById(id int) (*types.Ride, error) {
	rows, err := s.db.Query("SELECT id_ride, ds_ride, vl_ride,  dt_init, dt_finish, fg_count_weekend  FROM ride WHERE id_ride = $1", id)
	if err != nil {
		return nil, err
	}

	ride := &types.Ride{}
	for rows.Next() {
		ride, err = scanRowIntoRide(rows)

		if err != nil {
			return nil, err
		}
	}

	paymentRows, err := s.db.Query("SELECT id_ride_payment, vl_payment,  ds_person, fg_payed FROM ride_payment WHERE id_ride = $1", id)
	if err != nil {
		return nil, err
	}

	ride.Payments = make([]types.RidePayment, 0)
	for paymentRows.Next() {
		payment := types.RidePayment{}
		err := paymentRows.Scan(&payment.IdRidePayment, &payment.VlPayment, &payment.DsPerson, &payment.FgPayed)
		if err != nil {
			return nil, err
		}
		ride.Payments = append(ride.Payments, payment)
	}

	presenceRows, err := s.db.Query(`
		SELECT p.id_presence, p.id_ride_payment, p.qt_presence, p.dt_ride 
		FROM presence p
		INNER JOIN ride_payment bp ON p.id_ride_payment = bp.id_ride_payment
		WHERE bp.id_ride = $1`, id)
	if err != nil {
		return nil, err
	}

	presences := make([]types.Presence, 0)

	for presenceRows.Next() {
		presence := types.Presence{}
		err := presenceRows.Scan(&presence.IdPresence, &presence.IdRidePayment, &presence.QtPresence, &presence.DtRide)
		if err != nil {
			return nil, err
		}
		presences = append(presences, presence)
	}

	// Agrupar presenças
	groupMap := make(map[string][]types.Presence)
	for _, p := range presences {
		key := p.DtRide.Format("2006-01-02")
		groupMap[key] = append(groupMap[key], p)
	}

	for _, ps := range groupMap {
		ride.GroupedPresences = append(ride.GroupedPresences, types.GroupedPresence{
			DtRide:    ps[0].DtRide,
			Presences: ps,
		})
	}

	if ride.IdRide == 0 {
		return nil, nil
	}

	return ride, nil
}

func (s *Store) CreateRide(ridePayload types.Ride) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var id int
	err = tx.QueryRow(`
		INSERT INTO ride (ds_ride, vl_ride,  dt_init, dt_finish, fg_count_weekend)
		VALUES ($1, $2, $3, $4, $5) RETURNING id_ride
	`,
		ridePayload.DsRide,
		ridePayload.VlRide,
		ridePayload.DtInit,
		ridePayload.DtFinish,
		ridePayload.FgCountWeekend,
	).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}

	var paymentIDs []int
	for i := 0; i < len(ridePayload.Payments); i++ {
		var pid int
		err := tx.QueryRow(`
			INSERT INTO ride_payment (vl_payment, ds_person, fg_payed, id_ride)
			VALUES ($1, $2, $3, $4)
			RETURNING id_ride_payment
		`,
			0,
			ridePayload.Payments[i].DsPerson,
			false,
			id,
		).Scan(&pid)
		if err != nil {
			tx.Rollback()
			return err
		}
		paymentIDs = append(paymentIDs, pid)
	}

	currentDate := ridePayload.DtInit
	for !currentDate.After(ridePayload.DtFinish) {
		if ridePayload.FgCountWeekend ||
			(currentDate.Weekday() != 6 && currentDate.Weekday() != 0) {
			for _, id := range paymentIDs {
				_, err := tx.Exec(`
					INSERT INTO presence (id_ride_payment, dt_ride, qt_presence)
					VALUES ($1, $2, 0)
				`, id, currentDate)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func scanRowIntoRide(rows *sql.Rows) (*types.Ride, error) {
	ride := &types.Ride{}

	err := rows.Scan(
		&ride.IdRide,
		&ride.DsRide,
		&ride.VlRide,
		&ride.DtInit,
		&ride.DtFinish,
		&ride.FgCountWeekend,
	)

	if err != nil {
		return nil, err
	}

	return ride, nil
}
