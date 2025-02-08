package types

import (
	"time"
)

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateBillPayload struct {
	DsBill   string        `json:"dsBill" validate:"required"`
	VlBill   float64       `json:"vlBill" validate:"required"`
	QtPerson float64       `json:"qtPerson" validate:"required"`
	Payments []BillPayment `json:"payments,omitempty"`
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(u User) error
}

type BillStore interface {
	GetBills() ([]Bill, error)
	GetBillById(id int) (*Bill, error)
	CreateBill(b Bill) error
	UpdateBill(b Bill) error
}

type Bill struct {
	IdBill   int           `json:"idBill"`
	DsBill   string        `json:"dsBill"`
	VlBill   float64       `json:"vlBill"`
	QtPerson float64       `json:"qtPerson"`
	Payments []BillPayment `json:"payments"`
}

type BillPayment struct {
	IdPayment       int          `json:"idPayment"`
	IdBill          int          `json:"idBill"`
	VlPayment       float64      `json:"vlPayment"`
	FgPayed         bool         `json:"fgPayed"`
	FgCustomPayment bool         `json:"fgCustomPayment,omitempty"`
	DsPerson        string       `json:"dsPerson"`
}
