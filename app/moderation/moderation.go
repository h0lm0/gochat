package moderation

import (
	"gochat/client"
	dbmodels "gochat/database/models"
	"gochat/models"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
	xterm "golang.org/x/term"
)

func CheckAdmin(user *dbmodels.User) bool {
	return user.Role == models.RoleAdmin
}

func KickUser(currentSession ssh.Session, username string, reason string) string {
	for _, user := range models.ConnectedUsers {
		if user.Session.User() == username {
			kickReason := reason
			if kickReason == "" {
				kickReason = "no reason"
			}
			kickMessage := models.ErrorStyle.Sprintf("System> kicked for: " + kickReason + "\n")
			userIp := strings.Split(user.Session.RemoteAddr().String(), ":")[0]
			user.Terminal.Write([]byte(kickMessage))
			err := user.Session.Exit(403)
			if err != nil {
				return models.ErrorStyle.Sprintf("User: " + username + " not kicked: " + err.Error() + "\n")
			}
			if currentSession != nil {
				log.Println("User " + username + " | IP: " + userIp + " | kicked by " + currentSession.User() + " for: " + kickReason)
			} else {
				log.Println("User " + username + " | IP: " + userIp + " | kicked by system for: " + kickReason)
			}
			return models.SuccessStyle.Sprintf("User " + username + " kicked ! \n")
		}
	}
	return models.ErrorStyle.Sprintf("User " + username + " not found\n")
}

func BanIP(currentSession ssh.Session, username string, reason string) string {
	for _, user := range models.ConnectedUsers {
		if user.Session.User() == username {
			banReason := reason
			if banReason == "" {
				banReason = "no reason"
			}
			ipToBan := strings.Split(user.Session.RemoteAddr().String(), ":")[0]
			models.BannedIPS = append(models.BannedIPS, ipToBan)
			banMessage := models.ErrorStyle.Sprintf("System> banned for: " + banReason + "\n")
			user.Terminal.Write([]byte(banMessage))
			err := user.Session.Exit(403)
			if err != nil {
				return models.ErrorStyle.Sprintf("User " + username + " not banned: " + err.Error() + "\n")
			}
			log.Println("User: " + username + " | IP: " + ipToBan + " | Banned by " + currentSession.User() + " for: " + banReason)
			return models.SuccessStyle.Sprintf("User " + username + " banned ! \n")
		}
	}
	return models.ErrorStyle.Sprintf("User " + username + " not found\n")
}

func UnbanIP(currentSession ssh.Session, ip string) string {
	for i, bannedIp := range models.BannedIPS {
		if bannedIp == ip {
			models.BannedIPS = append(models.BannedIPS[:i], models.BannedIPS[i+1:]...)
			log.Println("IP: " + ip + " unbanned by " + currentSession.User())
			return models.SuccessStyle.Sprintf(ip + " unbanned\n")
		}
	}
	return models.ErrorStyle.Sprintf("IP " + ip + " not banned\n")
}

func CheckBan(s ssh.Session, term *xterm.Terminal) {
	for _, ip := range models.BannedIPS {
		if ip == strings.Split(s.RemoteAddr().String(), ":")[0] {
			term.Write([]byte(models.ErrorStyle.Sprintln("System> IP blacklist")))
			client.GracefulExit(s, 403, false, nil)
		}
	}
}

func ListBanned() string {
	list := ""
	for _, bannedIp := range models.BannedIPS {
		list += models.SystemStyle.Sprintf(" " + bannedIp + "\n")
	}
	return list
}
