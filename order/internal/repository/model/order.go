package model

import (
	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentMethodUnknown       PaymentMethod = "PAYMENT_METHOD_UNKNOWN"
	PaymentMethodCard          PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSBP           PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCreditCard    PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)

type Status string

const (
	StatusPendingPayment Status = "PENDING_PAYMENT"
	StatusPaid           Status = "PAID"
	StatusCanceled       Status = "CANCELED"
)

type Order struct {
	OrderID       uuid.UUID
	UserID        uuid.UUID
	Parts         []uuid.UUID
	TotalPrice    float64
	TransactionID uuid.UUID
	PaymentMethod PaymentMethod
	Status        Status
}
