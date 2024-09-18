package models

import (
	dbmodels "gochat/database/models"

	"github.com/gliderlabs/ssh"
)

type Room struct {
	Name string
	Type int
	// History  []Message
	Users    []dbmodels.User
	Password string
	Motd     string
}

var (
	Sessions map[ssh.Session]*Room
	// AvailableRooms []*Room
)

// func (r *Room) Enter(u dbmodels.User) {
// 	r.Users = append(r.Users, u)
// 	if r.Motd != "" {
// 		motd := EntryStyle.Sprintf(r.Motd)
// 		name := EntryStyle.Sprintf(r.Name + "> ")
// 		entryMsg := Message{From: name, Message: motd}
// 		send(u, "", entryMsg)
// 	}
// 	welcomeMessage := Message{
// 		From:    SystemStyle.Sprintf("System> "),
// 		Message: SystemStyle.Sprintf(fmt.Sprintf("%s has joined the room.", u.Username)),
// 	}
// 	for _, m := range r.History {
// 		if strings.Contains(m.From, u.Session.User()) {
// 			send(u, UserStyle.Sprintf("> "), Message{From: UserStyle.Sprintf(m.From), Message: m.Message})
// 		} else if strings.Contains(m.From, "System") {
// 			send(u, "", m)
// 		} else {
// 			send(u, "> ", m)
// 		}
// 	}
// 	r.SendMessage(welcomeMessage.From, "", welcomeMessage.Message)
// }

// func (r *Room) Leave(sess ssh.Session) {
// 	r.Users = RemoveByUsername(r.Users, sess.User())
// 	byeStyle := pterm.NewStyle(pterm.FgLightRed)
// 	// Envoyer un message pour notifier les autres utilisateurs que cet utilisateur a quitté la salle
// 	leaveMessage := Message{
// 		From:    byeStyle.Sprintf("System> "),
// 		Message: byeStyle.Sprintf(fmt.Sprintf("%s has left the room.", sess.User())),
// 	}
// 	r.SendMessage(leaveMessage.From, "", leaveMessage.Message)
// 	Sessions[sess] = nil
// }

// func (r *Room) SendMessage(from, separator string, message string) {
// 	messageObj := Message{From: from, Message: message}
// 	r.History = append(r.History, messageObj)

// 	if len(r.History) > 100 {
// 		r.History = r.History[1:]
// 	}

// 	for _, u := range r.Users {
// 		if u.Session.User() != from {
// 			send(u, separator, messageObj)
// 		}
// 	}
// }

// func RemoveByUsername(s []dbmodels.User, n string) []dbmodels.User {
// 	for i, u := range s {
// 		if u.Session.User() == n {
// 			return append(s[:i], s[i+1:]...)
// 		}
// 	}
// 	return s
// }

// func send(u dbmodels.User, separator string, m Message) {
// 	raw := m.From + separator + m.Message + "\n"
// 	if strings.Contains(m.Message, "@"+u.Session.User()) || strings.Contains(m.Message, "@everyone") {
// 		mentionStyle := pterm.NewStyle(pterm.FgLightBlue, pterm.Bold)
// 		raw = m.From + separator + mentionStyle.Sprintf(m.Message+"\n")
// 	}
// 	u.Terminal.Write([]byte(raw))
// }

// func LoadRoomsFromFile(filePath string) ([]*Room, error) {
// 	file, err := os.Open(filePath)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	bytes, err := io.ReadAll(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var rooms []*Room
// 	if err := json.Unmarshal(bytes, &rooms); err != nil {
// 		return nil, err
// 	}

// 	return rooms, nil
// }

// func ListRooms() string {
// 	var sb strings.Builder
// 	for _, r := range AvailableRooms {
// 		sb.WriteString(SuccessStyle.Sprintf(" #" + r.Name + "\n"))
// 		for _, u := range r.Users {
// 			sb.WriteString(" " + u.Session.User() + "\n")
// 		}
// 	}
// 	return sb.String()
// }

// func Filter[T any](s []T, cond func(t T) bool) []T {
// 	res := []T{}
// 	for _, v := range s {
// 		if cond(v) {
// 			res = append(res, v)
// 		}
// 	}
// 	return res
// }
