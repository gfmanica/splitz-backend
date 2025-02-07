package bill

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

func (s *Store) GetBills() ([]types.Bill, error) {
	rows, err := s.db.Query("SELECT * FROM bill")

	if err != nil {
		return nil, err
	}

	bills := make([]types.Bill, 0)

	for rows.Next() {
		bill, err := scanRowIntoBill(rows)

		if err != nil {
			return nil, err
		}

		bills = append(bills, *bill)
	}

	return bills, nil
}

func (s *Store) GetBillById(id int) (*types.Bill, error) {
	rows, err := s.db.Query("SELECT * FROM bill WHERE id_bill = $1", id)

	if err != nil {
		return nil, err
	}

	bill := &types.Bill{}

	for rows.Next() {
		bill, err = scanRowIntoBill(rows)

		if err != nil {
			return nil, err
		}
	}

	if bill.IdBill == 0 {
		return nil, nil
	}

	return bill, nil
}

func scanRowIntoBill(rows *sql.Rows) (*types.Bill, error) {
	u := &types.Bill{}

	err := rows.Scan(
		&u.IdBill,
		&u.DsBill,
		&u.VlBill,
		&u.QtPerson,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}
