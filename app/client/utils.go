package client

import (
	dbmodels "gochat/database/models"
	"gochat/database/services"
	"gochat/models"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/pterm/pterm"
	xterm "golang.org/x/term"
)

func GracefulExit(s ssh.Session, code int, needLog bool, currentUser *dbmodels.User) {
	if currentUser != nil {
		if currentUser.Room != nil {
			err := services.LeaveRoom(currentUser)
			if err != nil {
				log.Println(err.Error())
			}
		}
		for i, user := range models.ConnectedUsers {
			if s == user.Session {
				models.ConnectedUsers = append(models.ConnectedUsers[:i], models.ConnectedUsers[i+1:]...)
				break
			}
		}
		services.DisconnectUser(*currentUser)
	}

	if needLog {
		log.Println(s.User() + " | " + s.RemoteAddr().String() + " disconnected ")
	}
	s.Exit(code)
}

func SendMessage(message string, style *pterm.Style, term *xterm.Terminal) {
	formattedMessage := style.Sprint(message)
	term.Write([]byte(formattedMessage))
}
