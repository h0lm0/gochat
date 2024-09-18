package terminals

import (
	dbmodels "gochat/database/models"
	"gochat/keydb"
	"strconv"

	xterm "golang.org/x/term"
)

var (
	TermIndex = -1
	Terminals = []*xterm.Terminal{}
)

func GetTerminal(currentUser *dbmodels.User) *xterm.Terminal {
	strIndex, _ := keydb.GetKey(currentUser.Username)
	index, _ := strconv.Atoi(strIndex)
	return Terminals[index]
}
