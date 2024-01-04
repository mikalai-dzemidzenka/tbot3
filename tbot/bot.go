package tbot

import (
	"fmt"
	"github.com/sampgo/sampgo"
	"strconv"
	"strings"
)

var (
	Bots           = make(map[botNumber]*BotInfo)
	IsNicksVisible = false
)

const (
	BotPrefix = "TBot"
	MaxBots   = 50
)

const (
	NotReady = iota
	Ready
	Done
)

type BotInfo struct {
	BotNumber  int
	BotGroupID int
	Skin       int
	CarInfo
	SeatID   int
	IsSingle bool
	State    int

	BotRuntimeInfo `json:"-"`
}

type BotRuntimeInfo struct {
	id                int
	nickTextDraw      int
	recordingPlayerID int
	car               int
	tickToStart       uint
}

func NewBotRuntimeInfo(car int) BotRuntimeInfo {
	return BotRuntimeInfo{
		id:                BotNotConnected,
		nickTextDraw:      0,
		recordingPlayerID: NoRecordingBotNumber,
		car:               car,
		tickToStart:       0,
	}
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

	sampgo.SetPlayerSkin(id, Bots[botNum].Skin)

	if Bots[botNum].car != NoCar {
		sampgo.PutPlayerInVehicle(id, Bots[botNum].car, Bots[botNum].SeatID)
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

func AttachBotNick(botNumber int) {
	bot, ok := Bots[botNumber]
	if !ok {
		return
	}

	var nick string
	if bot.IsSingle {
		nick = fmt.Sprintf("%sS%d id%d g%d", BotPrefix, botNumber, bot.id, bot.BotGroupID)
	} else {
		nick = fmt.Sprintf("%s%d id%d g%d", BotPrefix, botNumber, bot.id, bot.BotGroupID)
	}
	label := sampgo.Create3DTextLabel(nick, 0x28BA9AFF, 0, 0, 0, 200, -1, false)
	sampgo.Attach3DTextLabelToPlayer(label, bot.id, 0, 0, 0.3)
	Bots[botNumber].nickTextDraw = label
}

func DetachBotNick(botNumber int) {
	_, ok := Bots[botNumber]
	if !ok {
		return
	}

	if Bots[botNumber].nickTextDraw != 0 {
		sampgo.Delete3DTextLabel(Bots[botNumber].nickTextDraw)
		Bots[botNumber].nickTextDraw = 0
	}
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

	if IsNicksVisible {
		DetachBotNick(botNumber)
	}

	delete(Players, bot.id)
	Bots[botNumber].BotRuntimeInfo = NewBotRuntimeInfo(Bots[botNumber].car)
}
