package tbot

import (
	"fmt"
	"github.com/sampgo/sampgo"
	"strconv"
	"strings"
)

const (
	BotPrefix = "TBot"
	MaxBots   = 50
)

type BotInfo struct {
	Number     int
	BotGroupID int
	Skin       int
	Car        int
	SeatID     int
	IsSingle   bool
	State      int
	// set only in runtime
	id                int
	nickTextDraw      int
	recordingPlayerID int
	ready             bool
}

func (b *BotInfo) String() string {
	return fmt.Sprintf("ID: %d, GROUP: %d, SKIN_ID: %d, SINGLE: %t", b.id, b.BotGroupID, b.Skin, b.IsSingle)
}

func IsBotConnected(botNum int) bool {
	if _, ok := Bots[botNum]; !ok {
		return false
	}
	return Bots[botNum].id != BotNotConnected
}

func ConnectBot(id int) {
	sampgo.SpawnPlayer(id)

	var name string
	sampgo.GetPlayerName(id, &name, 24)
	if !strings.HasPrefix(name, BotPrefix) {
		return
	}

	botNum := BotNumberFromName(name)

	if IsNicksVisible {
		AttachBotNick(id, botNum)
	}

	sampgo.SetPlayerSkin(id, Bots[botNum].Skin)

	if Bots[botNum].Car != NoCar {
		sampgo.PutPlayerInVehicle(id, Bots[botNum].Car, Bots[botNum].SeatID)
	}

	Bots[botNum].id = id
	Players[id] = &PlayerInfo{
		BotInfo: Bots[botNum],
	}
}

func GetFreeBotNum() (int, bool) {
	for i := 0; i < MaxBots; i++ {
		if _, ok := Bots[i]; !ok {
			return i, true
		}
	}
	return 0, false
}

func AttachBotNick(id int, botNumber int) {
	nick := fmt.Sprintf("%s%d", BotPrefix, botNumber)
	label := sampgo.Create3DTextLabel(nick, 0x28BA9AFF, 0, 0, 0, 200, -1, false)
	sampgo.Attach3DTextLabelToPlayer(label, id, 0, 0, 0.3)
	Bots[botNumber].nickTextDraw = label
}

func BotNumberFromName(name string) int {
	idStr := strings.TrimPrefix(name, BotPrefix)
	id, _ := strconv.Atoi(idStr)
	return id
}

func IsBot(id int) bool {
	var name string
	sampgo.GetPlayerName(id, &name, 24)
	return strings.HasPrefix(name, BotPrefix)
}

func DisconnectBot(botNumber int) {
	bot, ok := Bots[botNumber]
	if !ok {
		return
	}
	if bot.nickTextDraw != 0 {
		bot.nickTextDraw = 0
		sampgo.Delete3DTextLabel(bot.nickTextDraw)
	}

	delete(Players, bot.id)
	Bots[botNumber].id = BotNotConnected
}
