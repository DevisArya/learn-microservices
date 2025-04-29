package entity

type Field struct {
	Id          uint   `gorm:"primary_key"`
	Name        string `gorm:"size:100;not null"`
	Type        string `gorm:"size:50;not null"`
	Description string `gorm:"type:text"`
	Price       uint32 `gorm:"not null"`

	Schedule []Schedule `gorm:"foreignKey:FieldId"`
}
