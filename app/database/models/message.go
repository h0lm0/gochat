package database

type Message struct {
	ID      uint `gorm:"primaryKey"`
	FromID  uint
	From    User   `gorm:"foreignKey:FromID;constraint:OnDelete:CASCADE;"` // Relation avec User
	Message []byte //`gorm:"size:256"`
	RoomID  uint
	Room    Room `gorm:"constraint:OnDelete:CASCADE;"`
}
