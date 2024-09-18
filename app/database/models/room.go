package database

type Room struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"uniqueIndex;size:255"`
	Type     int
	History  []Message `gorm:"foreignKey:RoomID"`
	Users    []User    `gorm:"many2many:room_users;constraint:OnDelete:CASCADE;"`
	Password string
	Motd     string
}
