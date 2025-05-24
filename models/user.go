package models

type User struct {
	ID       string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`

	Incomes  []Income  `gorm:"foreignKey:UserID"`
	Expenses []Expense `gorm:"foreignKey:UserID"`
}
