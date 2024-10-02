package services

import (
	"gochat/database"
	dbmodels "gochat/database/models"
	"gochat/models"
	"log"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/bcrypt"
	xterm "golang.org/x/term"
	"gorm.io/gorm"
)

func CreateUser(username string, password string, role string) (bool, string) {
	// Vérifier si un utilisateur avec ce nom d'utilisateur existe déjà
	var existingUser dbmodels.User
	result := database.Db.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		return false, "Username already exists"
	} else if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return false, result.Error.Error()
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	createResult := database.Db.Create(&dbmodels.User{Username: username, Password: string(hashedPassword), Role: role})
	if createResult.Error != nil {
		return false, createResult.Error.Error()
	}
	return true, ""
}

func DeleteUser(username string) (bool, string) {
	_, connectedUsers, _ := ListUsers()
	for _, user := range connectedUsers {
		if user.Username == username {
			DisconnectUser(user)
		}
	}

	deleteResult := database.Db.Where("username = ?", username).Delete(&dbmodels.User{})
	if deleteResult.Error != nil {
		return false, deleteResult.Error.Error()
	}

	if deleteResult.RowsAffected == 0 {
		return false, "No user found with the given username"
	}

	return true, ""
}

func BanUser(currentSession ssh.Session, username string, reason string) string {
	var user dbmodels.User
	result := database.Db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return models.ErrorStyle.Sprintf("User " + username + " not found\n")
	}
	user.Banned = true
	updateResult := database.Db.Save(&user)
	if updateResult.Error != nil {
		log.Println("Error during ban of " + username + ": " + updateResult.Error.Error())
	}
	banReason := reason
	if banReason == "" {
		banReason = "no reason"
	}
	banMessage := models.ErrorStyle.Sprintf("System> banned for: " + banReason + "\n")
	_, connectedUsers, _ := ListUsers()

	for _, user := range connectedUsers {
		if user.Session.User() == username {
			user.Terminal.Write([]byte(banMessage))
			err := user.Session.Exit(403)
			if err != nil {
				return models.ErrorStyle.Sprintf("User " + username + " not kicked: " + err.Error() + "\n")
			}
		}
	}
	log.Println("User: " + username + " banned by " + currentSession.User() + " for: " + banReason)
	return models.SuccessStyle.Sprintf("User " + username + " banned ! \n")
}

func UnbanUser(username string) bool {
	var user dbmodels.User
	result := database.Db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return false
	}
	user.Banned = false
	updateResult := database.Db.Save(&user)
	return updateResult.Error == nil
}

func CheckUser(username string, plainpassword string) (bool, *dbmodels.User) {
	var user dbmodels.User
	result := database.Db.Where("username = ? and banned = ? and connected = ?", username, false, false).First(&user)
	if result.Error != nil {
		return false, nil
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainpassword))
	if err != nil {
		return false, nil
	}
	return true, &user
}

func ConnectUser(s ssh.Session, term *xterm.Terminal, user *dbmodels.User) bool {
	user.Session = s
	user.Color = models.RandomColor()
	user.Terminal = term
	models.ConnectedUsers = append(models.ConnectedUsers, *user)
	user.Connected = true
	updateResult := database.Db.Save(&user)

	refreshResult := database.Db.First(user, user.ID)
	if refreshResult.Error != nil {
		log.Println("Error refreshing data for user "+user.Username+" on ConnectUser(): ", refreshResult.Error)
		return false
	}

	return updateResult.Error == nil
}

func DisconnectUser(user dbmodels.User) bool {
	user.Connected = false
	updateResult := database.Db.Save(&user)
	user.Terminal.Write([]byte("Disconnected\n"))
	user.Session.Close()
	return updateResult.Error == nil
}

func ListUsers() (string, []dbmodels.User, []dbmodels.User) {
	var users []dbmodels.User
	result := database.Db.Find(&users)
	if result.Error != nil {
		return "Error retrieving users: " + result.Error.Error(), nil, nil
	}

	var connectedUsers []dbmodels.User
	var disconnectedUsers []dbmodels.User

	for _, user := range users {
		if user.Connected {
			connectedUsers = append(connectedUsers, user)
		} else {
			disconnectedUsers = append(disconnectedUsers, user)
		}
	}
	return "", connectedUsers, disconnectedUsers

}
