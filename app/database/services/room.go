package services

import (
	"errors"
	"fmt"
	"gochat/database"
	dbmodels "gochat/database/models"
	"gochat/models"
	"gochat/server"
	"gochat/terminals"
	"log"
	"strings"

	"github.com/pterm/pterm"
	"gorm.io/gorm"
)

func CreateRoom(name string, room_type int) (bool, string) {
	var existingRoom dbmodels.Room
	result := database.Db.Where("name = ?", name).First(&existingRoom)
	if result.Error == nil {
		return false, "Room already exists"
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error.Error()
	}
	createResult := database.Db.Create(&dbmodels.Room{Name: name, Type: room_type})
	if createResult.Error != nil {
		return false, createResult.Error.Error()
	}
	return true, ""
}

func ListRooms() ([]dbmodels.Room, []dbmodels.Room, error) {
	var rooms []dbmodels.Room
	result := database.Db.Preload("Users").Find(&rooms)

	if result.Error != nil {
		return nil, nil, fmt.Errorf("Error retrieving rooms: " + result.Error.Error())
	}

	var publicRooms []dbmodels.Room
	var privateRooms []dbmodels.Room

	for _, room := range rooms {
		if room.Type == 0 {
			publicRooms = append(publicRooms, room)
		} else {
			privateRooms = append(privateRooms, room)
		}
	}
	return publicRooms, privateRooms, nil
}

func JoinRoom(roomName string, currentUser *dbmodels.User) string {
	var room *dbmodels.Room
	result := database.Db.Preload("History", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("History.From").Where("name = ?", roomName).First(&room)

	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return "Error: room " + roomName + " not found"
	}

	if result.Error != nil {
		return "Error retrieving rooms: " + result.Error.Error()
	}

	room.Users = append(room.Users, *currentUser)
	roomUpdate := database.Db.Save(&room)
	if roomUpdate.Error != nil {
		return "Error joining room " + roomName + ": " + roomUpdate.Error.Error()
	}
	currentUser.Room = room
	userUpdate := database.Db.Save(&currentUser)
	if userUpdate.Error != nil {
		return "Error updating user " + currentUser.Username + ": " + userUpdate.Error.Error()
	}
	if room.Motd != "" {
		motd := models.EntryStyle.Sprintf(room.Motd)
		entryMsg := dbmodels.Message{From: dbmodels.User{Username: room.Name}, Message: []byte(motd)}
		send(*currentUser, ">", entryMsg)
	}
	var systemUser *dbmodels.User
	selectSystem := database.Db.Where("username = ?", "System").First(&systemUser)
	if selectSystem.Error != nil {
		return "Error selecting system user: " + selectSystem.Error.Error()
	}
	welcomeMessage := dbmodels.Message{
		From:    *systemUser,
		Message: []byte(models.SystemStyle.Sprintf(fmt.Sprintf(" %s has joined the room.", currentUser.Username))),
	}
	for _, m := range room.History {
		send(*currentUser, "> ", m)
	}
	SendRoomMessage(room, &welcomeMessage.From, ">", string(welcomeMessage.Message))
	return ""
}

func LeaveRoom(currentUser *dbmodels.User) error {
	if err := database.Db.Preload("Room").First(&currentUser, currentUser.ID).Error; err != nil {
		return errors.New("Error reloading user: " + err.Error())
	}
	if currentUser.Room == nil {
		return errors.New("user is not in any room")
	}

	room := currentUser.Room

	if err := database.Db.Preload("Users").Preload("History").First(&room, room.ID).Error; err != nil {
		return errors.New("error loading room: " + err.Error())
	}

	var updatedUsers []dbmodels.User
	for _, user := range room.Users {
		if user.ID != currentUser.ID {
			updatedUsers = append(updatedUsers, user)
		}
	}
	room.Users = updatedUsers

	if err := database.Db.Save(&room).Error; err != nil {
		return errors.New("error updating room users: " + err.Error())
	}

	currentUser.Room = nil
	if err := database.Db.Save(currentUser).Error; err != nil {
		return errors.New("error updating  user: " + err.Error())
	}

	var systemUser *dbmodels.User
	if err := database.Db.Where("username = ?", "System").First(&systemUser).Error; err != nil {
		return errors.New("error selecting system user: " + err.Error())
	}

	leaveMessage := dbmodels.Message{
		From:    *systemUser,
		Message: []byte(models.SystemStyle.Sprintf(fmt.Sprintf(" %s has left the room.", currentUser.Username))),
	}

	SendRoomMessage(room, &leaveMessage.From, ">", string(leaveMessage.Message))

	return nil
}

func SendRoomMessage(room *dbmodels.Room, from *dbmodels.User, separator, message string) {
	if err := database.Db.Select("id, name").Preload("Users").Preload("History").Find(&room).Error; err != nil {
		log.Panicln("Error reloading users of room "+room.Name+" :", err)
		return
	}
	term := terminals.GetTerminal(from)
	encryptedData, encryptErr := server.EncryptData(*from, message)
	if encryptErr != nil {
		term.Write([]byte(models.ErrorStyle.Sprintf(encryptErr.Error() + "\n")))
		log.Println("Error sending message to room "+": ", encryptErr.Error())
		return
	}
	messageObj := dbmodels.Message{From: *from, Message: []byte(encryptedData)}
	message = "REDACTED"
	if len(room.History) > 100 {
		room.History = []dbmodels.Message{}
	}
	room.History = append(room.History, messageObj)
	roomUpdate := database.Db.Save(&room)
	if roomUpdate.Error != nil {
		log.Println("Error updating room history: " + roomUpdate.Error.Error())
		return
	}
	for _, u := range room.Users {
		if u.Username != from.Username {
			send(u, separator, messageObj)
		}
	}
}

func send(u dbmodels.User, separator string, m dbmodels.Message) {
	term := terminals.GetTerminal(&u)
	plainMessage, err := server.DecryptData(m.From, m.Message)
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	raw := pterm.NewStyle(pterm.Color(m.From.Color)).Sprintf(m.From.Username+separator) + plainMessage + "\n"
	if m.From.Username == u.Username {
		raw = models.UserStyle.Sprintf(m.From.Username+separator) + plainMessage + "\n"
	}
	if strings.Contains(plainMessage, "@"+u.Username) || strings.Contains(plainMessage, "@everyone") {
		raw = m.From.Username + separator + models.MentionStyle.Sprintf(plainMessage+"\n")
	}
	term.Write([]byte(raw))
}
