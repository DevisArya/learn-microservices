package entity

type TransactionDetail struct {
	Id            uint   `gorm:"size:35;primary_key"`
	TransactionId string `gorm:"not null"`
	ScheduleId    uint   `gorm:"not null"`
	Name          string `gorm:"size:255;not null"`
	Price         uint32 `gorm:"not null"`

	// Transaction Transaction `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
