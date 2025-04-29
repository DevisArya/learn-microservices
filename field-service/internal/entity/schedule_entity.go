package entity

import (
	"time"
)

type ScheduleStatus string

const (
	ScheduleStatusAvailable ScheduleStatus = "available"
	ScheduleStatusReserved  ScheduleStatus = "reserved"
	ScheduleStatusSold      ScheduleStatus = "sold"
)

type Schedule struct {
	Id                uint              `gorm:"primaryKey"`
	UserId            uint              `gorm:"not null"`
	FieldId           uint              `gorm:"not null"`
	Date              time.Time         `gorm:"not null"`
	Status            ScheduleStatus    `gorm:"type:enum('available', 'reserved', 'sold')"`
	TransactionDetail TransactionDetail `gorm:"foreignKey:ScheduleId"`

	User  User  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Field Field `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
