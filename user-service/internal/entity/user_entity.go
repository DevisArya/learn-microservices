package entity

type Role string

const (
	RoleUser      Role = "user"
	RoleOperator  Role = "operator"
	RoleSuperUser Role = "super user"
)

type User struct {
	Id          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Email       string `gorm:"size:255;unique;not null"`
	Password    string `gorm:"size:255;not null"`
	PhoneNumber string `gorm:"size:20"`
	Role        Role   `gorm:"type:enum('user', 'operator', 'super user')"`
}
