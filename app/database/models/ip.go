package database

type Ip struct {
	ID uint   `gorm:"primaryKey"`
	Ip string `gorm:"uniqueIndex;size:255"`
}
