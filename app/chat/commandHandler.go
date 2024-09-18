package chat

import (
	"gochat/client"
	dbmodels "gochat/database/models"
	"gochat/database/services"
	"gochat/models"
	"gochat/moderation"
	"log"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
	xterm "golang.org/x/term"
)

func handleAdminCommands(line string, currentUser *dbmodels.User) {
	if !moderation.CheckAdmin(currentUser) {
		client.SendMessage("Unauthorized access\n", models.ErrorStyle, currentUser.Terminal)
		return
	}
	parts := strings.Fields(line)
	if len(parts) < 2 {
		client.SendMessage("\n Administration usage:\n", models.CommandStyle, currentUser.Terminal)
		currentUser.Terminal.Write([]byte(models.AdminHelpMsg()))
		return
	}

	command := parts[1]
	switch command {
	case models.KickCmd:
		if len(parts) < 3 {
			client.SendMessage("Invalid user\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		user := parts[2]
		reason := ""
		if len(parts) > 3 {
			reason = parts[3]
		}
		currentUser.Terminal.Write([]byte(moderation.KickUser(currentUser.Session, user, reason)))
	case models.BanCmd:
		if len(parts) < 3 {
			client.SendMessage("Invalid user\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		user := parts[2]
		reason := ""
		if len(parts) > 3 {
			reason = parts[3]
		}
		currentUser.Terminal.Write([]byte(services.BanUser(currentUser.Session, user, reason)))
	case models.BanIPCmd:
		if len(parts) < 3 {
			client.SendMessage("Invalid user\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		user := parts[2]
		reason := ""
		if len(parts) > 3 {
			reason = parts[3]
		}
		currentUser.Terminal.Write([]byte(moderation.BanIP(currentUser.Session, user, reason)))
	case models.UnbanCmd:
		if len(parts) < 3 {
			client.SendMessage("Invalid username\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		username := parts[2]
		if services.UnbanUser(username) {
			client.SendMessage("User "+username+" unbanned\n", models.SuccessStyle, currentUser.Terminal)
			log.Println("User " + username + " unbanned by " + currentUser.Username + " | " + currentUser.Session.RemoteAddr().String())
		} else {
			client.SendMessage("User "+username+" not found", models.ErrorStyle, currentUser.Terminal)
		}
	case models.UnbanIPCmd:
		if len(parts) < 3 {
			client.SendMessage("Invalid IP\n", pterm.NewStyle(pterm.Bold, pterm.FgRed), currentUser.Terminal)
			return
		}
		ip := parts[2]
		currentUser.Terminal.Write([]byte(moderation.UnbanIP(currentUser.Session, ip)))
	case models.ListBannedCmd:
		client.SendMessage("\n Banned IPs list:\n", models.CommandStyle, currentUser.Terminal)
		currentUser.Terminal.Write([]byte(moderation.ListBanned() + "\n"))
	default:
		client.SendMessage("Unknown admin command\n", models.ErrorStyle, currentUser.Terminal)
	}
}

func handleUserCommands(line string, currentUser *dbmodels.User) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		client.SendMessage("\n User usage:\n", models.CommandStyle, currentUser.Terminal)
		currentUser.Terminal.Write([]byte(models.UserHelpMsg()))
		return
	}

	command := parts[1]
	switch command {
	case models.ListCmd:
		err, connectedUsers, disconnectedUsers := services.ListUsers()
		if err != "" {
			client.SendMessage(err, models.ErrorStyle, currentUser.Terminal)
			return
		}
		client.SendMessage("\n Connected users:\n\n", models.CommandStyle, currentUser.Terminal)
		for _, user := range connectedUsers {
			client.SendMessage(" "+user.Username+"\n", models.SuccessStyle, currentUser.Terminal)
		}
		client.SendMessage("\n Disconnected users:\n\n", models.CommandStyle, currentUser.Terminal)
		for _, user := range disconnectedUsers {
			client.SendMessage(" "+user.Username+"\n", models.ErrorStyle, currentUser.Terminal)
		}
		currentUser.Terminal.Write([]byte("\n"))
	case models.CreateCmd:
		if !moderation.CheckAdmin(currentUser) {
			client.SendMessage("Unauthorized access\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		if len(parts) < 3 {
			client.SendMessage("Invalid user\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		user := parts[2]
		password, pErr := currentUser.Terminal.ReadPassword("User password: ")
		if pErr != nil {
			client.SendMessage(pErr.Error(), models.ErrorStyle, currentUser.Terminal)
			log.Println(pErr.Error() + " " + currentUser.Username + " | " + currentUser.Session.RemoteAddr().String())
			client.GracefulExit(currentUser.Session, 1, true, currentUser)
			return
		}

		client.SendMessage("Does the user need admin privilleges ? (y/N)\n", pterm.NewStyle(pterm.FgWhite), currentUser.Terminal)
		role, rErr := currentUser.Terminal.ReadLine()
		if rErr != nil {
			client.SendMessage(rErr.Error(), models.ErrorStyle, currentUser.Terminal)
			log.Println(rErr.Error() + " " + currentUser.Username + " | " + currentUser.Session.RemoteAddr().String())
			client.GracefulExit(currentUser.Session, 1, true, currentUser)
			return
		}
		userRole := ""
		switch role {
		case "y", "Y":
			userRole = models.RoleAdmin
		case "n", "N", "":
			userRole = models.RoleGuest
		}

		success, err := services.CreateUser(user, password, userRole)
		if success {
			client.SendMessage("User "+user+" created\n", models.SuccessStyle, currentUser.Terminal)
			log.Println("User " + user + " (" + userRole + ") created by " + currentUser.Username + " | " + currentUser.Session.RemoteAddr().String())
		} else {
			client.SendMessage("User "+user+" not created: "+err+"\n", models.ErrorStyle, currentUser.Terminal)
			log.Println("User " + user + " (" + userRole + ") failed to create by " + currentUser.Username + ":" + err + " | " + currentUser.Session.RemoteAddr().String())
		}
	case models.DeleteCmd:
		if !moderation.CheckAdmin(currentUser) {
			client.SendMessage("Unauthorized access\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		if len(parts) < 3 {
			client.SendMessage("Invalid user\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		user := parts[2]
		success, err := services.DeleteUser(user)
		if !success {
			client.SendMessage(err+"\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		client.SendMessage("User "+user+" deleted\n", models.SuccessStyle, currentUser.Terminal)
		log.Println("User " + user + " deleted by " + currentUser.Username + " | " + currentUser.Session.RemoteAddr().String())
	default:
		client.SendMessage("Unknown user command\n", models.ErrorStyle, currentUser.Terminal)
	}
}

func handleRoomCommands(line string, currentUser *dbmodels.User) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		client.SendMessage("\n Room usage:\n", models.CommandStyle, currentUser.Terminal)
		currentUser.Terminal.Write([]byte(models.RoomHelpMsg()))
		return
	}

	command := parts[1]
	switch command {
	case models.ListCmd:
		publicRooms, privateRooms, err := services.ListRooms()
		if err != nil {
			client.SendMessage("Error retrieving rooms: \n"+err.Error(), models.ErrorStyle, currentUser.Terminal)
		} else {
			client.SendMessage("\n Public rooms:\n", models.CommandStyle, currentUser.Terminal)
			for _, room := range publicRooms {
				client.SendMessage(" ■ "+room.Name+"\n", models.EntryStyle, currentUser.Terminal)
				for _, user := range room.Users {
					if user == room.Users[len(room.Users)-1] {
						client.SendMessage(" │\n", models.RoomStyle, currentUser.Terminal)
						client.SendMessage(" └─ ", models.RoomStyle, currentUser.Terminal)
						client.SendMessage(user.Username+"\n", models.BasicStyle, currentUser.Terminal)
					} else {
						client.SendMessage(" │\n", models.RoomStyle, currentUser.Terminal)
						client.SendMessage(" ├─ ", models.RoomStyle, currentUser.Terminal)
						client.SendMessage(user.Username+"\n", models.BasicStyle, currentUser.Terminal)
					}
				}
				client.SendMessage("\n", models.UserStyle, currentUser.Terminal)
			}
			client.SendMessage("\n Private rooms:\n", models.CommandStyle, currentUser.Terminal)
			for _, room := range privateRooms {
				client.SendMessage(" ■ "+room.Name+"\n", models.EntryStyle, currentUser.Terminal)
			}
			client.SendMessage("\n", models.UserStyle, currentUser.Terminal)
		}
		// currentUser.Terminal.Write([]byte(models.ListRooms() + "\n"))
	case models.JoinCmd:
		parts := strings.Fields(line)
		if len(parts) < 3 {
			client.SendMessage("Invalid room\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		roomName := parts[2]
		err := services.JoinRoom(roomName, currentUser)
		if err != "" {
			client.SendMessage(err+"\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
	case models.CreateCmd:
		if !moderation.CheckAdmin(currentUser) {
			client.SendMessage("Unauthorized access\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		client.SendMessage("Room name: ", models.CommandStyle, currentUser.Terminal)
		roomName, err := currentUser.Terminal.ReadLine()
		if err != nil {
			client.SendMessage(err.Error()+"\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		client.SendMessage("Room type (0: public/1: private): ", models.CommandStyle, currentUser.Terminal)
		roomType, err := currentUser.Terminal.ReadLine()
		if err != nil {
			client.SendMessage(err.Error()+"\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		intRoomeType, errConv := strconv.Atoi(roomType)
		if errConv != nil {
			client.SendMessage(errConv.Error()+"\n", models.ErrorStyle, currentUser.Terminal)
			return
		}
		if roomName != "" && (intRoomeType == 0 || intRoomeType == 1) {
			services.CreateRoom(roomName, intRoomeType)
		} else {
			client.SendMessage("\nBad format (room name or type)\n\n", models.ErrorStyle, currentUser.Terminal)

		}
	case models.LeaveCmd:
		if currentRoom := models.Sessions[currentUser.Session]; currentRoom != nil {
			services.LeaveRoom(currentUser)
		}
	default:
		client.SendMessage("Unknown room command\n", models.ErrorStyle, currentUser.Terminal)
	}
}

func handleMainCommands(line string, currentUser *dbmodels.User) {
	switch {
	case models.ExitCmd.MatchString(line):
		client.GracefulExit(currentUser.Session, 0, true, currentUser)
	case models.ClearCmd.MatchString(line):
		clearScreen(currentUser.Terminal)
	case models.HelpCmd.MatchString(line):
		client.SendMessage("\n Usage:\n", models.CommandStyle, currentUser.Terminal)
		currentUser.Terminal.Write([]byte(models.HelpMsg()))
	default:
		client.SendMessage("Unknown command\n", models.ErrorStyle, currentUser.Terminal)
	}
}

// func enterRoom(line string, currentUser dbmodels.User) {
// 	parts := strings.Fields(line)
// 	if len(parts) < 3 {
// 		client.SendMessage("Invalid room\n", models.ErrorStyle, currentUser.Terminal)
// 		return
// 	}
// 	toEnter := parts[2]
// 	matching := models.Filter(models.AvailableRooms, func(r *models.Room) bool {
// 		return toEnter == r.Name
// 	})
// 	if len(matching) == 0 {
// 		client.SendMessage("Invalid room\n", models.ErrorStyle, currentUser.Terminal)
// 		return
// 	}
// 	if currentRoom := models.Sessions[currentUser.Session]; currentRoom != nil {
// 		currentRoom.Leave(currentUser.Session)
// 	}
// 	r := matching[0]
// 	if r.Type == 1 {
// 		line, err := currentUser.Terminal.ReadPassword("The room you're trying to join is private. Please enter the key:\n")
// 		if err != nil {
// 			return
// 		}
// 		hash := sha256.New()
// 		hash.Write([]byte(line))
// 		hashed := fmt.Sprintf("%x", hash.Sum(nil))
// 		if r.Password != hashed {
// 			currentUser.Terminal.Write([]byte("Invalid password\n"))
// 			return
// 		}
// 	}
// 	r.Enter(currentUser)
// 	models.Sessions[currentUser.Session] = r
// }

func clearScreen(t *xterm.Terminal) {
	_, err := t.Write([]byte("\033[H\033[2J"))
	if err != nil {
		log.Printf("Failed to clear screen: %v", err)
	}
}
