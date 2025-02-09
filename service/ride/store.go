package ride

import (
	"database/sql"
	"time"

	"github.com/gfmanica/splitz-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetRides() ([]types.Ride, error) {
	rows, err := s.db.Query("SELECT id_ride, ds_ride, vl_ride,  dt_init, dt_finish, fg_count_weekend, qt_ride FROM ride")

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
	rows, err := s.db.Query("SELECT id_ride, ds_ride, vl_ride,  dt_init, dt_finish, fg_count_weekend, qt_ride  FROM ride WHERE id_ride = $1", id)
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

func (s *Store) UpdateRide(ridePayload types.Ride) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Atualiza os dados da raiz do ride
	_, err = tx.Exec(`
		UPDATE ride SET ds_ride = $1, vl_ride = $2, dt_init = $3, dt_finish = $4, fg_count_weekend = $5
		WHERE id_ride = $6
	`, ridePayload.DsRide, ridePayload.VlRide, ridePayload.DtInit, ridePayload.DtFinish, ridePayload.FgCountWeekend, ridePayload.IdRide)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Lida com os pagamentos
	// Recupera os pagamentos existentes
	rows, err := tx.Query(`SELECT id_ride_payment FROM ride_payment WHERE id_ride = $1`, ridePayload.IdRide)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingPayments := make(map[int]bool)
	for rows.Next() {
		var id int
		rows.Scan(&id)
		existingPayments[id] = true
	}
	rows.Close()

	payloadPaymentsIDs := make(map[int]bool)
	// Para cada pagamento enviado, atualiza ou insere
	for _, p := range ridePayload.Payments {
		if p.IdRidePayment != 0 {
			// Atualiza pagamento existente
			payloadPaymentsIDs[p.IdRidePayment] = true
			_, err = tx.Exec(`
				UPDATE ride_payment SET ds_person = $1, fg_payed = $2
				WHERE id_ride_payment = $3
			`, p.DsPerson, p.FgPayed, p.IdRidePayment)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Insere pagamento novo
			var newID int
			err = tx.QueryRow(`
				INSERT INTO ride_payment (vl_payment, ds_person, fg_payed, id_ride)
				VALUES ($1, $2, $3, $4) RETURNING id_ride_payment
			`, 0, p.DsPerson, false, ridePayload.IdRide).Scan(&newID)
			if err != nil {
				tx.Rollback()
				return err
			}
			payloadPaymentsIDs[newID] = true
		}
	}
	// Deleta pagamentos que foram removidos
	for id := range existingPayments {
		if !payloadPaymentsIDs[id] {
			_, err = tx.Exec(`DELETE FROM ride_payment WHERE id_ride_payment = $1`, id)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Lida com as presenças
	// Se houver alteração nas datas ou na flag de fim de semana, delete todas as presenças
	_, err = tx.Exec(`
		DELETE FROM presence 
		WHERE id_ride_payment IN (SELECT id_ride_payment FROM ride_payment WHERE id_ride = $1)
	`, ridePayload.IdRide)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insere as presenças conforme a nova grade, considerando possíveis atualizações enviadas no groupedPresences
	currentDate := ridePayload.DtInit
	for !currentDate.After(ridePayload.DtFinish) {
		// Aplica a regra de fim de semana
		if ridePayload.FgCountWeekend || (currentDate.Weekday() != time.Saturday && currentDate.Weekday() != time.Sunday) {
			for id := range payloadPaymentsIDs {
				// Procura se há qtPresence para este dia e payment enviado do front (se não, usa 0)
				qtPresence := 0
				found := false
				for _, gp := range ridePayload.GroupedPresences {
					if gp.DtRide.Equal(currentDate) {
						for _, ps := range gp.Presences {
							if ps.IdRidePayment == id {
								qtPresence = ps.QtPresence
								found = true
								break
							}
						}
					}
					if found {
						break
					}
				}
				_, err = tx.Exec(`
					INSERT INTO presence (id_ride_payment, dt_ride, qt_presence)
					VALUES ($1, $2, $3)
				`, id, currentDate, qtPresence)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Recalcula e atualiza os valores de cada pagamento (vl_payment)
	// A lógica agrupa por dia e, em cada dia, distribui (vlRide * qtRide) proporcionalmente à qtPresence
	paymentTotals := make(map[int]float64)
	date := ridePayload.DtInit
	for !date.After(ridePayload.DtFinish) {
		if ridePayload.FgCountWeekend || (date.Weekday() != time.Saturday && date.Weekday() != time.Sunday) {
			rows, err := tx.Query(`
				SELECT id_ride_payment, qt_presence 
				FROM presence 
				WHERE dt_ride = $1
			`, date)
			if err != nil {
				tx.Rollback()
				return err
			}
			dailyTotal := 0
			dailyPresences := make(map[int]int)
			for rows.Next() {
				var pid, qt int
				rows.Scan(&pid, &qt)
				dailyPresences[pid] = qt
				dailyTotal += qt
			}
			rows.Close()
			if dailyTotal > 0 {
				dailyCost := ridePayload.VlRide * float64(ridePayload.QtRide)
				for pid, qt := range dailyPresences {
					share := (float64(qt) / float64(dailyTotal)) * dailyCost
					paymentTotals[pid] += share
				}
			}
		}
		date = date.AddDate(0, 0, 1)
	}

	// Atualiza os vl_payment dos pagamentos
	for pid, total := range paymentTotals {
		_, err = tx.Exec(`
			UPDATE ride_payment SET vl_payment = $1 WHERE id_ride_payment = $2
		`, total, pid)
		if err != nil {
			tx.Rollback()
			return err
		}
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
		&ride.QtRide,
	)

	if err != nil {
		return nil, err
	}

	return ride, nil
}
