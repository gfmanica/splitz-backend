package bill

import (
	"database/sql"
	"fmt"

	"github.com/gfmanica/splitz-backend/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetBills() ([]types.Bill, error) {
	rows, err := s.db.Query("SELECT * FROM bill ORDER BY id_bill DESC")

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
	rows, err := s.db.Query("SELECT id_bill, ds_bill, vl_bill, qt_person FROM bill WHERE id_bill = $1", id)
	if err != nil {
		return nil, err
	}

	bill := &types.Bill{}
	if rows.Next() {
		err := rows.Scan(&bill.IdBill, &bill.DsBill, &bill.VlBill, &bill.QtPerson)
		if err != nil {
			return nil, err
		}
	}

	paymentRows, err := s.db.Query("SELECT id_bill_payment, vl_payment,  ds_person, fg_payed, fg_custom_payment, id_bill FROM bill_payment WHERE id_bill = $1 ORDER BY id_bill_payment ASC ", id)
	if err != nil {
		return nil, err
	}

	bill.Payments = make([]types.BillPayment, 0)
	for paymentRows.Next() {
		payment := types.BillPayment{}
		err := paymentRows.Scan(&payment.IdBillPayment, &payment.VlPayment, &payment.DsPerson, &payment.FgPayed, &payment.FgCustomPayment, &payment.IdBill)
		if err != nil {
			return nil, err
		}
		bill.Payments = append(bill.Payments, payment)
	}

	if bill.IdBill == 0 {
		return nil, nil
	}

	return bill, nil
}

func (s *Store) CreateBill(billPayload types.Bill) (*types.Bill, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	var id int

	err = tx.QueryRow("INSERT INTO bill (ds_bill, vl_bill, qt_person) VALUES ($1, $2, $3) RETURNING id_bill", billPayload.DsBill, billPayload.VlBill, billPayload.QtPerson).Scan(&id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	fmt.Print(id)

	totalVlPayments := 0.0

	for _, payment := range billPayload.Payments {
		totalVlPayments += payment.VlPayment
	}

	personVlBill := (billPayload.VlBill - totalVlPayments) / float64(billPayload.QtPerson)

	for i := 0; i < int(billPayload.QtPerson); i++ {
		dsPerson := fmt.Sprintf("Pessoa %d", i+1)
		vlPayment := personVlBill
		fgCustomPayment := false

		if i < len(billPayload.Payments) {
			if billPayload.Payments[i].DsPerson != "" {
				dsPerson = billPayload.Payments[i].DsPerson
			}
			if billPayload.Payments[i].VlPayment != 0 {
				vlPayment = billPayload.Payments[i].VlPayment
				fgCustomPayment = true
			}
		}

		_, err := tx.Exec("INSERT INTO bill_payment (vl_payment, ds_person, fg_payed, fg_custom_payment, id_bill) VALUES ($1, $2, $3, $4, $5)",
			vlPayment,
			dsPerson,
			false,
			fgCustomPayment,
			id,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return s.GetBillById(id)
}

func (s *Store) UpdateBill(billPayload types.Bill) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Atualizar a descrição, valor e quantidade de pessoas do bill
	_, err = tx.Exec("UPDATE bill SET ds_bill = $1, vl_bill = $2, qt_person = $3 WHERE id_bill = $4",
		billPayload.DsBill, billPayload.VlBill, billPayload.QtPerson, billPayload.IdBill)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Obter pagamentos existentes do banco de dados
	rows, err := tx.Query("SELECT id_bill_payment FROM bill_payment WHERE id_bill = $1", billPayload.IdBill)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingPayments := make(map[int]bool)
	for rows.Next() {
		var idBillPayment int
		err := rows.Scan(&idBillPayment)
		if err != nil {
			tx.Rollback()
			return err
		}
		existingPayments[idBillPayment] = true
	}

	// Atualizar ou inserir pagamentos
	for _, payment := range billPayload.Payments {
		if payment.IdBillPayment != 0 {
			// Atualizar pagamento existente
			_, err := tx.Exec("UPDATE bill_payment SET vl_payment = $1, ds_person = $2, fg_payed = $3, fg_custom_payment = $4 WHERE id_bill_payment = $5",
				payment.VlPayment, payment.DsPerson, payment.FgPayed, payment.FgCustomPayment, payment.IdBillPayment)
			if err != nil {
				tx.Rollback()
				return err
			}
			delete(existingPayments, payment.IdBillPayment)
		} else {
			// Inserir novo pagamento
			_, err := tx.Exec("INSERT INTO bill_payment (vl_payment, ds_person, fg_payed, fg_custom_payment, id_bill) VALUES ($1, $2, $3, $4, $5)",
				payment.VlPayment, payment.DsPerson, payment.FgPayed, payment.FgCustomPayment, billPayload.IdBill)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Excluir pagamentos que não foram enviados no payload
	for idBillPayment := range existingPayments {
		_, err := tx.Exec("DELETE FROM bill_payment WHERE id_bill_payment = $1", idBillPayment)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Inserir novos pagamentos até a quantidade total de pessoas (QtPerson)
	totalVlPayments := 0.0
	for _, payment := range billPayload.Payments {
		totalVlPayments += payment.VlPayment
	}

	personVlBill := (billPayload.VlBill - totalVlPayments) / float64(billPayload.QtPerson)
	for i := len(billPayload.Payments); i < int(billPayload.QtPerson); i++ {
		dsPerson := fmt.Sprintf("Pessoa %d", i+1)
		vlPayment := personVlBill
		fgCustomPayment := false

		_, err := tx.Exec("INSERT INTO bill_payment (vl_payment, ds_person, fg_payed, fg_custom_payment, id_bill) VALUES ($1, $2, $3, $4, $5)",
			vlPayment, dsPerson, false, fgCustomPayment, billPayload.IdBill)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Recalcular valores dos pagamentos não customizados
	rows, err = tx.Query("SELECT id_bill_payment, vl_payment, ds_person, fg_payed, fg_custom_payment, id_bill FROM bill_payment WHERE id_bill = $1", billPayload.IdBill)
	if err != nil {
		tx.Rollback()
		return err
	}

	payments := make([]types.BillPayment, 0)
	for rows.Next() {
		payment := types.BillPayment{}
		err := rows.Scan(&payment.IdBillPayment, &payment.VlPayment, &payment.DsPerson, &payment.FgPayed, &payment.FgCustomPayment, &payment.IdBill)
		if err != nil {
			tx.Rollback()
			return err
		}
		payments = append(payments, payment)
	}

	totalCustomPayments := 0.0
	nonCustomPaymentsCount := 0
	for _, payment := range payments {
		if payment.FgCustomPayment {
			totalCustomPayments += payment.VlPayment
		} else {
			nonCustomPaymentsCount++
		}
	}

	remainingAmount := billPayload.VlBill - totalCustomPayments
	nonCustomPaymentValue := remainingAmount / float64(nonCustomPaymentsCount)

	for _, payment := range payments {
		if !payment.FgCustomPayment {
			_, err := tx.Exec("UPDATE bill_payment SET vl_payment = $1 WHERE id_bill_payment = $2",
				nonCustomPaymentValue, payment.IdBillPayment)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteBill(id int) error {
	_, err := s.db.Exec("DELETE FROM bill WHERE id_bill = $1", id)
	if err != nil {
		return err
	}

	return nil
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
