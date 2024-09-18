package main

import (
	"gochat/chat"
	"gochat/database"
	"gochat/keydb"
	"gochat/models"
	"log"

	"github.com/gliderlabs/ssh"
)

const (
	APP_NAME = "Gochat"
	KeyCtrlC = 3
	KeyCtrlD = 4
)

func main() {
	keydb.InitClient()
	database.Init()
	models.Sessions = make(map[ssh.Session]*models.Room)

	ssh.Handle(func(s ssh.Session) {
		log.Println("New handle for " + s.RemoteAddr().String())
		chat.Render(s)
	})
	log.Println("starting " + APP_NAME + " server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
