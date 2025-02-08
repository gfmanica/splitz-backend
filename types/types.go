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

type CreateRidePayload struct {
	DsRide         string        `json:"dsRide" validate:"required"`
	VlRide         float64       `json:"vlRide" validate:"required"`
	DtInit         time.Time     `json:"dtInit" validate:"required"`
	DtFinish       time.Time     `json:"dtFinish" validate:"required"`
	FgCountWeekend bool          `json:"fgCountWeekend"`
	Payments       []RidePayment `json:"payments" validate:"required"`
}

type UpdateRidePayload struct {
	DsRide           string            `json:"dsRide" validate:"required"`
	VlRide           float64           `json:"vlRide" validate:"required"`
	QtPerson         int               `json:"qtPerson" validate:"required"`
	DtInit           time.Time         `json:"dtInit" validate:"required"`
	DtFinish         time.Time         `json:"dtFinish" validate:"required"`
	FgCountWeekend   bool              `json:"fgCountWeekend"`
	GroupedPresences []GroupedPresence `json:"groupedPresences"`
	Payments         []RidePayment     `json:"payments"`
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

type RideStore interface {
	GetRides() ([]Ride, error)
	GetRideById(id int) (*Ride, error)
	CreateRide(r Ride) error
	// UpdateRide(r Ride) error
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type Bill struct {
	IdBill   int           `json:"idBill"`
	DsBill   string        `json:"dsBill"`
	VlBill   float64       `json:"vlBill"`
	QtPerson float64       `json:"qtPerson"`
	Payments []BillPayment `json:"payments"`
}

type BillPayment struct {
	IdBillPayment   int     `json:"idBillPayment"`
	IdBill          int     `json:"idBill"`
	VlPayment       float64 `json:"vlPayment"`
	FgPayed         bool    `json:"fgPayed"`
	FgCustomPayment bool    `json:"fgCustomPayment"`
	DsPerson        string  `json:"dsPerson"`
}

type Ride struct {
	IdRide           int               `json:"idRide"`
	DsRide           string            `json:"dsRide"`
	VlRide           float64           `json:"vlRide"`
	DtInit           time.Time         `json:"dtInit"`
	DtFinish         time.Time         `json:"dtFinish"`
	FgCountWeekend   bool              `json:"fgCountWeekend"`
	GroupedPresences []GroupedPresence `json:"groupedPresences"`
	Payments         []RidePayment     `json:"payments"`
}

type RidePayment struct {
	IdRidePayment int     `json:"idRidePayment"`
	VlPayment     float64 `json:"vlPayment"`
	FgPayed       bool    `json:"fgPayed"`
	DsPerson      string  `json:"dsPerson"`
}

type Presence struct {
	IdPresence    int       `json:"idPresence"`
	IdRidePayment int       `json:"idRidePayment"`
	QtPresence    int       `json:"qtPresence"`
	DtRide        time.Time `json:"dtRide"`
}

type GroupedPresence struct {
	DtRide    time.Time  `json:"dtRide"`
	Presences []Presence `json:"presences"`
}
