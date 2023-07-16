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
	ID                int
	Number            int
	BotGroupID        int
	NickTextDraw      int
	Skin              int
	Car               int
	SeatID            int
	IsSingle          bool
	RecordingPlayerID int
	State             int
	Ready             bool
}

func (b *BotInfo) String() string {
	return fmt.Sprintf("ID: %d, GROUP: %d, SKIN_ID: %d, SINGLE: %t", b.ID, b.BotGroupID, b.Skin, b.IsSingle)
}

func (t *T) IsBotConnected(botNum int) bool {
	if _, ok := t.bots[botNum]; !ok {
		return false
	}
	return t.bots[botNum].ID != BotNotConnected
}

func (t *T) ConnectBot(id int) {
	sampgo.SpawnPlayer(id)

	var name string
	sampgo.GetPlayerName(id, &name, 24)
	if !strings.HasPrefix(name, BotPrefix) {
		return
	}

	botNum := BotNumberFromName(name)

	if t.isNicksVisible {
		t.AttachBotNick(id, botNum)
	}

	sampgo.SetPlayerSkin(id, t.bots[botNum].Skin)

	if t.bots[botNum].Car != NoCar {
		sampgo.PutPlayerInVehicle(id, t.bots[botNum].Car, t.bots[botNum].SeatID)
	}

	t.bots[botNum].ID = id
	t.players[id] = &PlayerInfo{
		BotInfo: t.bots[botNum],
	}
}

func (t *T) GetFreeBotNum() (int, bool) {
	for i := 0; i < MaxBots; i++ {
		if _, ok := t.bots[i]; !ok {
			return i, true
		}
	}
	return 0, false
}

func (t *T) AttachBotNick(id int, botNumber int) {
	nick := fmt.Sprintf("%s%d", BotPrefix, botNumber)
	label := sampgo.Create3DTextLabel(nick, 0x28BA9AFF, 0, 0, 0, 200, -1, false)
	sampgo.Attach3DTextLabelToPlayer(label, id, 0, 0, 0.3)
	t.bots[botNumber].NickTextDraw = label
}

func BotNumberFromName(name string) int {
	idStr := strings.TrimPrefix(name, BotPrefix)
	id, _ := strconv.Atoi(idStr)
	return id
}

func (t *T) IsBot(id int) bool {
	var name string
	sampgo.GetPlayerName(id, &name, 24)
	return strings.HasPrefix(name, BotPrefix)
}

func (t *T) DisconnectBot(botNumber int) {
	bot, ok := t.bots[botNumber]
	if !ok {
		return
	}
	if bot.NickTextDraw != 0 {
		bot.NickTextDraw = 0
		sampgo.Delete3DTextLabel(bot.NickTextDraw)
	}

	delete(t.players, t.bots[botNumber].ID)
	t.bots[botNumber].ID = BotNotConnected
}
