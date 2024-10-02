package chat

import (
	"fmt"
	"gochat/client"
	database "gochat/database/models"
	"gochat/database/services"
	"gochat/keydb"
	"gochat/models"
	"gochat/moderation"
	"gochat/terminals"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/pterm/pterm"
	xterm "golang.org/x/term"
)

const (
	MSG_ALERT = `#################################################################
#                   _    _           _   _                      #
#                  / \  | | ___ _ __| |_| |                     #
#                 / _ \ | |/ _ \ '__| __| |                     #
#                / ___ \| |  __/ |  | |_|_|                     #
#               /_/   \_\_|\___|_|   \__(_)                     #
#                                                               #
#  You are entering into a secured area! Your IP, Login Time,   #
#   Username has been noted and has been sent to the server     #
#                       administrator                           #
#   This service is restricted to authorized users only. All    #
#            activities on this system are logged.              #
#################################################################
`
	MSG_WELCOME = "Welcome to gochat !"
)

func Render(s ssh.Session) {
	term := xterm.NewTerminal(s, models.SuccessStyle.Sprint(fmt.Sprintf("%s> ", s.User())))

	success, currentUser := Login(s, term)
	if !success {
		return
	}
	for {
		line, err := term.ReadLine()

		if err != nil && err == io.EOF {
			log.Println("Session interrupt requested by " + currentUser.Username)
			client.GracefulExit(s, 1, true, currentUser)
			break
		}

		if err != nil {
			term.Write([]byte("Invalid entry, please refer to usage with /help\n"))
		}

		if len(line) > 0 {
			if strings.HasPrefix(line, "/admin") {
				handleAdminCommands(line, currentUser)
			} else if strings.HasPrefix(line, "/user") {
				handleUserCommands(line, currentUser)
			} else if strings.HasPrefix(line, "/room") {
				handleRoomCommands(line, currentUser)
			} else if string(line[0]) == "/" {
				handleMainCommands(line, currentUser)
			} else {
				if currentUser.Room != nil {
					services.SendRoomMessage(currentUser.Room, currentUser, "> ", line)
				} else {
					term.Write([]byte("Invalid entry, please refer to usage with /help\n"))
				}
			}
		}
	}
}

func Login(s ssh.Session, term *xterm.Terminal) (bool, *database.User) {
	moderation.CheckBan(s, term)
	clearScreen(term)

	client.SendMessage(MSG_ALERT, models.AlertStyle, term)
	serverPass, err := term.ReadPassword("Password: ")
	if err != nil {
		client.SendMessage(err.Error(), models.ErrorStyle, term)
		log.Println(err.Error() + " " + s.User() + " | " + s.RemoteAddr().String())
		client.GracefulExit(s, 1, true, nil)
		return false, nil
	}

	access, currentUser := services.CheckUser(s.User(), serverPass)
	if !access {
		client.SendMessage("Permission denied\n", models.ErrorStyle, term)
		log.Println("Permission denied | User: " + s.User() + " | " + s.RemoteAddr().String())
		client.GracefulExit(s, 1, true, nil)
		return false, nil
	}
	terminals.TermIndex++
	services.ConnectUser(s, term, currentUser)
	setErr := keydb.SetKey(currentUser.Username, strconv.Itoa(terminals.TermIndex))
	if setErr != nil {
		client.SendMessage(setErr.Error()+"\n", models.ErrorStyle, term)
		log.Println("Error setting keydb terminal: " + setErr.Error() + " | " + s.RemoteAddr().String())
		client.GracefulExit(s, 1, true, nil)
		terminals.TermIndex--
		return false, nil
	}
	terminals.Terminals = append(terminals.Terminals, term)

	success, keyErr := keydb.AddEncryptionKey(*currentUser)
	if !success {
		client.SendMessage(keyErr.Error()+"\n", models.ErrorStyle, term)
		log.Println("Error setting keydb user key: " + keyErr.Error() + " | " + s.RemoteAddr().String())
		client.GracefulExit(s, 1, true, nil)
		terminals.TermIndex--
		return false, nil
	}
	log.Println("Authentication successfull for " + currentUser.Username + " | " + s.RemoteAddr().String())
	// test, err := client.EncryptData(*currentUser, "Test")
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println("Test encryption for user " + currentUser.Username + ": " + test)
	// testDecrypt, err := client.DecryptData(*currentUser, test)
	// if err != nil {
	// 	log.Println(testDecrypt)
	// }
	// log.Println("Test decryption for user " + currentUser.Username + ": " + testDecrypt)

	clearScreen(term)

	client.SendMessage(MSG_WELCOME+"\n", pterm.NewStyle(pterm.FgCyan), term)
	return true, currentUser
}
