package models

import (
	dbmodels "gochat/database/models"
	"regexp"

	"math/rand"

	"github.com/pterm/pterm"
)

const (
	KickCmd       = "kick"
	BanCmd        = "ban"
	BanIPCmd      = "banip"
	UnbanCmd      = "unban"
	UnbanIPCmd    = "unbanip"
	ListBannedCmd = "listbanned"
	ListCmd       = "list"
	JoinCmd       = "join"
	CreateCmd     = "create"
	DeleteCmd     = "delete"
	LeaveCmd      = "leave"
)

var (
	SshKeys []string
	// AdminPassword  string
	BannedIPS      []string
	ConnectedUsers []dbmodels.User

	SuccessStyle = pterm.NewStyle(pterm.FgGreen)
	ErrorStyle   = pterm.NewStyle(pterm.FgRed)
	AlertStyle   = pterm.NewStyle(pterm.FgLightRed, pterm.BgBlack, pterm.Fuzzy)
	EntryStyle   = pterm.NewStyle(pterm.FgMagenta, pterm.Bold)
	RoomStyle    = pterm.NewStyle(pterm.FgMagenta)
	CommandStyle = pterm.NewStyle(pterm.Italic, pterm.FgCyan)
	SystemStyle  = pterm.NewStyle(pterm.FgLightYellow)
	UserStyle    = pterm.NewStyle(pterm.FgGreen)
	MentionStyle = pterm.NewStyle(pterm.FgLightBlue, pterm.Bold)
	BasicStyle   = pterm.NewStyle(pterm.FgWhite)

	// Main commands
	// EnterCmd = regexp.MustCompile(`^/enter .*`)
	HelpCmd  = regexp.MustCompile(`^/help$`)
	ExitCmd  = regexp.MustCompile(`^/exit$`)
	ClearCmd = regexp.MustCompile(`^/clear$`)
	// leaveCmd       = regexp.MustCompile(`^/leave.*`)

	AdminCmd = regexp.MustCompile(`^/admin\s+(kick|ban|banip|unban|unbanip|listbanned)\s*(.*)$`)
	UserCmd  = regexp.MustCompile(`^/user\s+(list|create|delete)\s*(.*)$`)
	RoomCmd  = regexp.MustCompile(`^/room\s+(list|join|leave|create|delete)\s*(.*)$`)
)

func HelpMsg() string {
	return `
 /help: To display this message
 /clear: Clear terminal screen
 /user: Display help related to users
 /admin: Display help related to administration
 /room: Display help related to rooms
 /exit: To leave the server

`
}

func AdminHelpMsg() string {
	return `
 /admin: To display this message
 /admin | kick <username> [reason]: To kick a conncted user
	| ban <username> [reason]: To ban a user (db)
	| banip <username> [reason]: To ban a connected user (ip)
	| unban <username>: To unban a user
	| unbanip <ip>: To unban an ip
	| listbanned: To list IPs banned

`
}

func UserHelpMsg() string {
	return `
 /user: To display this message
 /user  | list: To list all connected users
	| create <username>: To create a new user
	| delete <username>: To delete an existing user

`
}

func RoomHelpMsg() string {
	return `
 /room: To display this message
 /room  | list: To list available rooms
 	| join <room_name>: To enter a room
 	| leave: To leave current room
 	| create <room_name>: To create a new room
	| delete <room_name>: To delete an existing room

`
}

func RandomColor() int {
	firstColor := rand.Intn(5) + 33
	secondColor := rand.Intn(7) + 91
	choice := rand.Intn(2)
	if choice == 0 {
		return firstColor
	}
	return secondColor
}
