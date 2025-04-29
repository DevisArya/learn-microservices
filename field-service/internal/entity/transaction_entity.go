package entity

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentRejected      PaymentStatus = "rejected"
	PaymentStatusSuccess PaymentStatus = "success"
)

type Transaction struct {
	// Id string
	TransactionId     string              `gorm:"size:35;primary_key"`
	UserId            uint                `gorm:"not null"`
	OrderId           string              `gorm:"size:100"`
	PaymentType       string              `gorm:"size:100"`
	PaymentUrl        string              `gorm:"size:255"`
	PaymentStatus     string              `gorm:"type:enum('unpaid', 'rejected', 'success')"`
	TotalPrice        uint32              `gorm:"not null"`
	TransactionTime   time.Time           `gorm:"type:datetime;null"`
	SettlementTime    time.Time           `gorm:"type:datetime;null"`
	FraudStatus       string              `gorm:"size:100"`
	TransactionDetail []TransactionDetail `gorm:"foreignKey:TransactionId"`

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
