package database

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;size:255"`
	Password  string
	Role      string
	Banned    bool
	Connected bool
	Session   ssh.Session    `gorm:"-"`
	Terminal  *term.Terminal `gorm:"-"`
	RoomID    *uint          `gorm:"index"`
	Room      *Room          `gorm:"foreignKey:RoomID;"`
	Color     int            `gorm:"default:32"`
}
